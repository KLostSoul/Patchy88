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

func TestSessionLoggerCreatesUniqueFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "logs")
	when := time.Date(2026, 10, 26, 20, 30, 0, 0, time.Local)
	a, err := newSessionLoggerAt([]string{dir}, when)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, err := newSessionLoggerAt([]string{dir}, when)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	if a.Path() == b.Path() || !strings.HasSuffix(a.Path(), "20261026_203000.log") || !strings.HasSuffix(b.Path(), "20261026_203000_001.log") {
		t.Fatalf("unexpected paths %q %q", a.Path(), b.Path())
	}
	if err := a.WriteLine("한글 패치 시작\n[xdelta3 오류/진단] stdout/stderr 테스트"); err != nil {
		t.Fatal(err)
	}
	if err := a.Sync(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(a.Path())
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"2026-", "한글 패치 시작", "[xdelta3 오류/진단] stdout/stderr 테스트"} {
		if !strings.Contains(string(data), s) {
			t.Fatalf("log is missing %q: %s", s, data)
		}
	}
	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	if err := a.WriteLine("should fail"); err == nil {
		t.Fatal("writing after Close succeeded")
	}
}

func TestSessionLoggerFallback(t *testing.T) {
	base := t.TempDir()
	blocked := filepath.Join(base, "not-a-directory")
	if err := os.WriteFile(blocked, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	fallback := filepath.Join(base, "fallback")
	logger, err := newSessionLoggerAt([]string{blocked, fallback}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	if filepath.Dir(logger.Path()) != fallback {
		t.Fatalf("fallback not used: %q", logger.Path())
	}
}

func TestSessionLoggerConcurrentWrites(t *testing.T) {
	logger, err := newSessionLoggerAt([]string{t.TempDir()}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	var wg sync.WaitGroup
	for i := 0; i < 80; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := logger.WriteLine("[진행] 로그 저장 테스트"); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if err := logger.Sync(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(logger.Path())
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(data), "[진행] 로그 저장 테스트"); n != 80 {
		t.Fatalf("lost log entries: %d of 80", n)
	}
}

func TestSessionLoggerNoStorage(t *testing.T) {
	base := t.TempDir()
	blocked := filepath.Join(base, "blocked")
	if err := os.WriteFile(blocked, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := newSessionLoggerAt([]string{blocked}, time.Now()); err == nil {
		t.Fatal("expected failure when no log storage available")
	}
}

func TestSessionLoggerRecordsVerifiedPatch(t *testing.T) {
	root := t.TempDir()
	original := []byte("original payload for log fixture")
	target := []byte("verified replacement")
	cue := []byte("FILE \"translated.bin\" BINARY\n")
	if err := os.WriteFile(filepath.Join(root, "original.img"), original, 0600); err != nil {
		t.Fatal(err)
	}
	logger, err := newSessionLoggerAt([]string{filepath.Join(root, "logs")}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	// In the fixture only: emit the exact expected target bytes, then let the
	// real engine verify hashes and publish BIN/CUE.
	decoder := scriptDecoder(t, "printf 'verified replacement' > \"$5\"")
	e := Engine{
		Root:    root,
		Decoder: decoder,
		Manifest: Manifest{
			Source: hashBytes(original),
			Target: OutputHashes{Hashes: hashBytes(target), Filename: "translated.bin"},
			Cue:    OutputHashes{Hashes: hashBytes(cue), Filename: "translated.cue"},
		},
		CueTemplate: cue,
	}
	logf := func(s string) {
		if err := logger.WriteLine(s); err != nil {
			t.Error(err)
		}
	}
	_, _, err = e.ApplyContext(context.Background(), root, logf)
	if err != nil {
		t.Fatal(err)
	}
	if err := logger.Sync(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(logger.Path())
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"[검증 성공] BIN", "[검증 성공] CUE", "[xdelta3] 정상 종료", "MD5", "SHA-256"} {
		if !strings.Contains(string(data), s) {
			t.Fatalf("expected %q in session log: %s", s, data)
		}
	}
}
