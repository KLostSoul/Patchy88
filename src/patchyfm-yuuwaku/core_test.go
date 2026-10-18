package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEqualHashes(t *testing.T) {
	a := Hashes{MD5: "AA", SHA256: "Bb", Size: 2}
	b := Hashes{MD5: "aa", SHA256: "bb", Size: 2}
	if !equalHashes(a, b) {
		t.Fatal("hash case comparison failed")
	}
	b.Size = 3
	if equalHashes(a, b) {
		t.Fatal("wrong size accepted")
	}
}

func TestFindSourceRecursivelyAndIgnoreOtherImages(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "disc", "archive")
	if err := os.MkdirAll(sub, 0700); err != nil {
		t.Fatal(err)
	}
	good := []byte("a validated CD data track for unit tests")
	goodPath := filepath.Join(sub, "something.BIN")
	if err := os.WriteFile(goodPath, good, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "bad.img"), []byte("different"), 0600); err != nil {
		t.Fatal(err)
	}
	e := Engine{Manifest: Manifest{Source: hashBytes(good)}}
	got, err := e.FindSource(root)
	if err != nil || got != goodPath {
		t.Fatalf("got %q, %v; wanted %q", got, err, goodPath)
	}
	if err := e.ValidateSource(got); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(goodPath, []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := e.ValidateSource(got); err == nil {
		t.Fatal("changed source accepted")
	}
	if _, err := e.FindSource(root); err == nil {
		t.Fatal("bad source located")
	}
}

func TestFindSourceRejectsDuplicates(t *testing.T) {
	root := t.TempDir()
	body := []byte("identical data")
	for _, name := range []string{"A.img", "B.bin"} {
		if err := os.WriteFile(filepath.Join(root, name), body, 0600); err != nil {
			t.Fatal(err)
		}
	}
	e := Engine{Manifest: Manifest{Source: hashBytes(body)}}
	if _, err := e.FindSource(root); err == nil || !strings.Contains(err.Error(), "2개") {
		t.Fatalf("ambiguous duplicate not rejected: %v", err)
	}
}

func TestFindSourceRejectsFileInsteadOfFolder(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "sample.img")
	if err := os.WriteFile(path, []byte("a"), 0600); err != nil {
		t.Fatal(err)
	}
	e := Engine{Manifest: Manifest{Source: hashBytes([]byte("a"))}}
	if _, err := e.FindSource(path); err == nil {
		t.Fatal("file path accepted as folder")
	}
}
