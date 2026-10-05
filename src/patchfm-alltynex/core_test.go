package main

import "testing"

func TestEqualHashes(t *testing.T) {
	a := Hashes{MD5: "aa", SHA256: "bb", Size: 1}
	if !equalHashes(a, a) {
		t.Fatal("same hash rejected")
	}
	b := a
	b.Size = 2
	if equalHashes(a, b) {
		t.Fatal("size mismatch accepted")
	}
}
