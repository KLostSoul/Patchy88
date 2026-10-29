//go:build !windows

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBoundedDecoderDiagnostics(t *testing.T) {
	var b boundedDiagnostics
	big := bytes.Repeat([]byte("X"), maxDecoderDiagnostics*3)
	n, err := b.Write(big)
	if err != nil || n != len(big) || len(b.data) != maxDecoderDiagnostics || !b.truncated {
		t.Fatalf("diagnostics were not bounded, n=%d len=%d err=%v", n, len(b.data), err)
	}
	if _, err := b.Write([]byte("tail")); err != nil || len(b.data) > maxDecoderDiagnostics || !strings.HasSuffix(b.String(), "tail") {
		t.Fatalf("diagnostics tail missing: %v", err)
	}
}

func fixtureEngine(t *testing.T, script string) (*Engine, string, string, string) {
	t.Helper()
	root := t.TempDir()
	original := filepath.Join(root, "original.IMG")
	targetBytes := []byte("korean translation build v1.02 image")
	cue := []byte("FILE \"Temptation-Korean-v1.02.bin\" BINARY\r\n  TRACK 01 MODE1/2352\r\n    INDEX 01 00:00:00\r\n")
	if err := os.WriteFile(original, []byte("original FM TOWNS disc"), 0600); err != nil {
		t.Fatal(err)
	}
	helper := filepath.Join(root, "test-xdelta3")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\n"+script+"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	e := &Engine{Root: root, Decoder: helper, Patch: "fake.xdelta", CueTemplate: cue,
		Manifest: Manifest{Source: hashBytes([]byte("original FM TOWNS disc")), Target: OutputHashes{Hashes: hashBytes(targetBytes), Filename: expectedBIN}, Cue: OutputHashes{Hashes: hashBytes(cue), Filename: expectedCUE}}}
	return e, root, original, helper
}

func TestApplyValidatesAndKeepsOriginal(t *testing.T) {
	e, root, original, _ := fixtureEngine(t, "printf 'korean translation build v1.02 image' > \"$5\"")
	bin, cue, err := e.Apply(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := e.ValidateBIN(bin); err != nil {
		t.Fatal(err)
	}
	if err := e.ValidateCUE(cue); err != nil {
		t.Fatal(err)
	}
	if err := e.ValidateSource(original); err != nil {
		t.Fatal("original changed:", err)
	}
	// A completed patch is safe to reuse without running the decoder again.
	e.Decoder = filepath.Join(root, "does_not_exist")
	if _, _, err := e.Apply(root, nil); err != nil {
		t.Fatal("idempotence:", err)
	}
	if len(mustTmpGlob(t, root)) != 0 {
		t.Fatal("temporary files remain")
	}
}

func mustTmpGlob(t *testing.T, root string) []string {
	t.Helper()
	b, e := filepath.Glob(filepath.Join(root, "PatchyFM-Yuuwaku-v1.02-*"))
	if e != nil {
		t.Fatal(e)
	}
	return b
}

func TestApplyRejectsBadDecoderAndCleansTemps(t *testing.T) {
	e, root, _, _ := fixtureEngine(t, "printf 'bad' > \"$5\"")
	if _, _, err := e.Apply(root, nil); err == nil || !strings.Contains(err.Error(), "검증 실패") {
		t.Fatalf("bad target accepted: %v", err)
	}
	if len(mustTmpGlob(t, root)) != 0 {
		t.Fatal("temporary file left")
	}
	if _, err := os.Stat(filepath.Join(root, expectedBIN)); !os.IsNotExist(err) {
		t.Fatalf("published a bad BIN: %v", err)
	}
}

func TestApplyPreservesExistingMismatchedOutput(t *testing.T) {
	e, root, _, _ := fixtureEngine(t, "exit 77")
	dest := filepath.Join(root, expectedBIN)
	if err := os.WriteFile(dest, []byte("user data"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := e.Apply(root, nil); err == nil || !strings.Contains(err.Error(), "덮어쓰지") {
		t.Fatalf("not rejected: %v", err)
	}
	data, err := os.ReadFile(dest)
	if err != nil || string(data) != "user data" {
		t.Fatal("modified existing output")
	}
}

func TestApplyTimeoutAndCleansTemps(t *testing.T) {
	old := decoderRunTimeout
	decoderRunTimeout = 150 * time.Millisecond
	defer func() { decoderRunTimeout = old }()
	e, root, _, _ := fixtureEngine(t, "printf 'unfinished' > \"$5\"\nexec sleep 4")
	start := time.Now()
	_, _, err := e.Apply(root, nil)
	if err == nil || !strings.Contains(err.Error(), "시간 제한") {
		t.Fatalf("expected timeout, got %v", err)
	}
	if time.Since(start) > 3*time.Second {
		t.Fatal("timeout took too long")
	}
	if len(mustTmpGlob(t, root)) != 0 {
		t.Fatal("timeout left .tmp file")
	}
}

func TestApplyRejectsFailedHelper(t *testing.T) {
	e, root, _, _ := fixtureEngine(t, "echo 'failure code' >&2\nexit 8")
	_, _, err := e.Apply(root, nil)
	if err == nil || !strings.Contains(err.Error(), "failure code") {
		t.Fatalf("stderr missing: %v", err)
	}
	if len(mustTmpGlob(t, root)) != 0 {
		t.Fatal("failure left .tmp file")
	}
}
