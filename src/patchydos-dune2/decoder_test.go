//go:build !windows

package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func scriptDecoder(t *testing.T, contents string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "fake-decoder.sh")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+contents+"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	return p
}
func withShortWatchdog(t *testing.T, total, stall time.Duration) {
	t.Helper()
	ot, os_, op, or := decoderTotalTimeout, decoderStallTimeout, decoderPollInterval, decoderProgressInterval
	decoderTotalTimeout = total
	decoderStallTimeout = stall
	decoderPollInterval = 20 * time.Millisecond
	decoderProgressInterval = 40 * time.Millisecond
	t.Cleanup(func() {
		decoderTotalTimeout = ot
		decoderStallTimeout = os_
		decoderPollInterval = op
		decoderProgressInterval = or
	})
}
func TestDecoderStderrAndNonzeroExit(t *testing.T) {
	e := Engine{Root: t.TempDir(), Decoder: scriptDecoder(t, "echo 'bad input' 1>&2; exit 9")}
	var mu sync.Mutex
	var lines []string
	logf := func(s string) { mu.Lock(); lines = append(lines, s); mu.Unlock() }
	err := e.runDecoder(context.Background(), "unused", filepath.Join(t.TempDir(), "x.tmp"), logf)
	if err == nil || !strings.Contains(err.Error(), "bad input") {
		t.Fatalf("stderr missing: %v", err)
	}
	mu.Lock()
	got := strings.Join(lines, " ")
	mu.Unlock()
	if !strings.Contains(got, "[xdelta3 오류/진단] bad input") {
		t.Fatalf("stderr wasn't streamed: %s", got)
	}
}
func TestDecoderStallKilled(t *testing.T) {
	withShortWatchdog(t, 2*time.Second, 120*time.Millisecond)
	e := Engine{Root: t.TempDir(), Decoder: scriptDecoder(t, "exec sleep 10")}
	start := time.Now()
	err := e.runDecoder(context.Background(), "unused", filepath.Join(t.TempDir(), "x.tmp"), nil)
	if err == nil || !strings.Contains(err.Error(), "응답 없음") {
		t.Fatalf("no stall warning: %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("not killed promptly")
	}
}
func TestDecoderTotalTimeout(t *testing.T) {
	withShortWatchdog(t, 140*time.Millisecond, 3*time.Second)
	e := Engine{Root: t.TempDir(), Decoder: scriptDecoder(t, "exec sleep 10")}
	err := e.runDecoder(context.Background(), "unused", filepath.Join(t.TempDir(), "x.tmp"), nil)
	if err == nil || !strings.Contains(err.Error(), "시간 제한") {
		t.Fatalf("expected time limit: %v", err)
	}
}
func TestDecoderCancellation(t *testing.T) {
	withShortWatchdog(t, 5*time.Second, 4*time.Second)
	e := Engine{Root: t.TempDir(), Decoder: scriptDecoder(t, "exec sleep 10")}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(80 * time.Millisecond); cancel() }()
	err := e.runDecoder(ctx, "unused", filepath.Join(t.TempDir(), "x.tmp"), nil)
	if err == nil || !strings.Contains(err.Error(), "사용자 요청") {
		t.Fatalf("expected cancellation: %v", err)
	}
}
func TestDecoderReportsProgress(t *testing.T) {
	withShortWatchdog(t, 2*time.Second, time.Second)
	e := Engine{Root: t.TempDir(), Decoder: scriptDecoder(t, "printf 'partial-data' > \"$5\"; sleep 0.20")}
	var mu sync.Mutex
	var lines []string
	logf := func(s string) { mu.Lock(); lines = append(lines, s); mu.Unlock() }
	tmp := filepath.Join(t.TempDir(), "part.tmp")
	err := e.runDecoder(context.Background(), "unused", tmp, logf)
	if err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	got := strings.Join(lines, " ")
	mu.Unlock()
	if !strings.Contains(got, "[진행]") || !strings.Contains(got, "12 bytes") {
		t.Fatalf("progress missing: %s", got)
	}
}

func TestApplyContextCancelCleansTempFile(t *testing.T) {
	withShortWatchdog(t, 5*time.Second, 4*time.Second)
	dir := t.TempDir()
	src := filepath.Join(dir, "original.img")
	source := []byte("original-data")
	if err := os.WriteFile(src, source, 0600); err != nil {
		t.Fatal(err)
	}
	cue := []byte("FILE \"translated.bin\" BINARY\n")
	e := Engine{
		Root:    dir,
		Decoder: scriptDecoder(t, "printf 'partial' > \"$5\"; exec sleep 10"),
		Manifest: Manifest{
			Source: hashBytes(source),
			Target: OutputHashes{Hashes: hashBytes([]byte("desired")), Filename: "translated.bin"},
			Cue:    OutputHashes{Hashes: hashBytes(cue), Filename: "translated.cue"},
		},
		CueTemplate: cue,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(130 * time.Millisecond); cancel() }()
	_, _, err := e.ApplyContext(ctx, dir, nil)
	if err == nil || !strings.Contains(err.Error(), "사용자 요청") {
		t.Fatalf("cancellation error missing: %v", err)
	}
	if b, er := os.ReadFile(src); er != nil || string(b) != string(source) {
		t.Fatal("original file changed")
	}
	tmps, er := filepath.Glob(filepath.Join(dir, "*.tmp"))
	if er != nil || len(tmps) != 0 {
		t.Fatalf("stale temp files: %v %v", tmps, er)
	}
	if _, er = os.Stat(filepath.Join(dir, "translated.bin")); !os.IsNotExist(er) {
		t.Fatal("partial final BIN was published")
	}
}
