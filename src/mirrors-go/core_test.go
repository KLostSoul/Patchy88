package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testHashes(data []byte) Hashes {
	m := md5.Sum(data)
	s := sha256.Sum256(data)
	return Hashes{MD5: hex.EncodeToString(m[:]), SHA256: hex.EncodeToString(s[:])}
}
func testSHA(data []byte) string { return testHashes(data).SHA256 }
func writeTestFile(t *testing.T, p string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, data, 0755); err != nil {
		t.Fatal(err)
	}
}
func fixture(t *testing.T) (*Engine, string, map[string]map[string][]byte, map[string][]byte) {
	t.Helper()
	root := t.TempDir()
	assets := filepath.Join(root, "assets")
	folder := filepath.Join(root, "input")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}
	sources := map[string]map[string][]byte{}
	target := map[string][]byte{"ccd": []byte("Japanese original CCD; Korean result identical"), "img": []byte("Korean img test output with all expected content"), "sub": []byte("Japanese original SUB; Korean result identical")}
	m := Manifest{Name: programName, Schema: manifestSchema, Source: map[string]map[string]SourceDef{}, Target: map[string]TargetDef{}}
	for _, ext := range exts {
		m.Target[ext] = TargetDef{Hashes: testHashes(target[ext]), Filename: "Mirrors_Kor1.00." + ext}
	}
	for _, edition := range editions {
		sources[edition] = map[string][]byte{}
		m.Source[edition] = map[string]SourceDef{}
		for _, ext := range exts {
			var data []byte
			if edition == "Japanese" && (ext == "ccd" || ext == "sub") {
				data = target[ext]
			} else {
				data = []byte(edition + " original " + ext + " test file")
			}
			sources[edition][ext] = data
			patch := edition + "_" + strings.ToUpper(ext) + ".xdelta"
			fakeVCDIFF := append([]byte{0xd6, 0xc3, 0xc4, 0x00}, []byte(edition+ext)...)
			writeTestFile(t, filepath.Join(assets, "patches", patch), fakeVCDIFF)
			m.Source[edition][ext] = SourceDef{Hashes: testHashes(data), Patch: "patches/" + patch, PatchSHA256: testSHA(fakeVCDIFF)}
		}
	}
	for _, name := range []string{"Mirrors_Kor1.00.cue", "disk1main.d88", "disk2game.d88"} {
		var data []byte
		if name == "Mirrors_Kor1.00.cue" {
			data = []byte("FILE \"Mirrors_Kor1.00.img\" BINARY\n")
		} else {
			data = []byte("same fake D88 payload")
		}
		writeTestFile(t, filepath.Join(assets, "extras", name), data)
		m.Extras = append(m.Extras, ExtraDef{Filename: name, SHA256: testSHA(data), Size: int64(len(data))})
	}
	// Fake xdelta3 mimics xdelta3's -q -d -s SOURCE PATCH OUTPUT arguments.
	// It emits deterministic test outputs and fails on a configurable patch.
	helper := []byte(`#!/bin/sh
[ "$1" = "-q" ] || exit 2
[ "$2" = "-d" ] || exit 2
[ "$3" = "-s" ] || exit 2
src="$4"; patch="$5"; out="$6"
case "$patch" in
    *_CCD.xdelta) printf '%s' 'Japanese original CCD; Korean result identical' > "$out" ;;
    *_IMG.xdelta) printf '%s' 'Korean img test output with all expected content' > "$out" ;;
    *_SUB.xdelta) printf '%s' 'Japanese original SUB; Korean result identical' > "$out" ;;
    *) exit 1 ;;
esac
`)
	writeTestFile(t, filepath.Join(assets, "xdelta3.exe"), helper)
	m.XdeltaSHA256 = testSHA(helper)
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(assets, manifestFile), b)
	e, err := NewEngine(assets)
	if err != nil {
		t.Fatal(err)
	}
	return e, folder, sources, target
}
func putInputs(t *testing.T, folder, edition string, sources map[string]map[string][]byte) {
	t.Helper()
	for _, ext := range exts {
		writeTestFile(t, filepath.Join(folder, "original_"+ext+"."+ext), sources[edition][ext])
	}
}
func assertOutput(t *testing.T, e *Engine, folder string) {
	t.Helper()
	for _, ext := range exts {
		p := filepath.Join(folder, e.Manifest.Target[ext].Filename)
		got, err := hashFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if !hashEqual(got, e.Manifest.Target[ext].Hashes) {
			t.Errorf("wrong output: %s", p)
		}
	}
	for _, x := range e.Manifest.Extras {
		p := filepath.Join(folder, x.Filename)
		got, err := hashFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if got.Size != x.Size || !strings.EqualFold(got.SHA256, x.SHA256) {
			t.Errorf("wrong extra: %s", p)
		}
	}
}
func TestJapanesePatchAndExtras(t *testing.T) {
	e, folder, sources, _ := fixture(t)
	putInputs(t, folder, "Japanese", sources)
	s, err := e.Scan(folder)
	if err != nil {
		t.Fatal(err)
	}
	if s.Edition != "Japanese" {
		t.Fatal(s.Edition)
	}
	if err = e.Apply(s, nil); err != nil {
		t.Fatal(err)
	}
	assertOutput(t, e, folder)
	s, err = e.Scan(folder)
	if err != nil {
		t.Fatal(err)
	}
	if s.Edition != "AlreadyPatched" {
		t.Fatal("rerun:", s.Edition)
	}
	if err = e.Apply(s, nil); err != nil {
		t.Fatal(err)
	}
}
func TestEnglishPatchAndExtras(t *testing.T) {
	e, folder, sources, _ := fixture(t)
	putInputs(t, folder, "English", sources)
	s, err := e.Scan(folder)
	if err != nil {
		t.Fatal(err)
	}
	if s.Edition != "English" {
		t.Fatal(s.Edition)
	}
	if err = e.Apply(s, nil); err != nil {
		t.Fatal(err)
	}
	assertOutput(t, e, folder)
}
func TestMixedEditionsRejectedNoWrites(t *testing.T) {
	e, folder, sources, _ := fixture(t)
	writeTestFile(t, filepath.Join(folder, "source.ccd"), sources["Japanese"]["ccd"])
	writeTestFile(t, filepath.Join(folder, "source.img"), sources["English"]["img"])
	writeTestFile(t, filepath.Join(folder, "source.sub"), sources["Japanese"]["sub"])
	if _, err := e.Scan(folder); err == nil {
		t.Fatal("mixed edition accepted")
	}
	items, _ := os.ReadDir(folder)
	if len(items) != 3 {
		t.Fatal("files changed", len(items))
	}
}
func TestCollisionRejectedBeforePatching(t *testing.T) {
	e, folder, sources, _ := fixture(t)
	putInputs(t, folder, "English", sources)
	writeTestFile(t, filepath.Join(folder, "Mirrors_Kor1.00.cue"), []byte("unrelated file"))
	if _, err := e.Scan(folder); err == nil {
		t.Fatal("extra collision accepted")
	}
	if _, err := os.Stat(filepath.Join(folder, e.Manifest.Target["img"].Filename)); !os.IsNotExist(err) {
		t.Fatal("output incorrectly created")
	}
}
func TestTamperedAssetsRejected(t *testing.T) {
	e, _, _, _ := fixture(t)
	p := filepath.Join(e.Root, "patches", "English_IMG.xdelta")
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.Write([]byte("oops")); err != nil {
		t.Fatal(err)
	}
	f.Close()
	if _, err = NewEngine(e.Root); err == nil {
		t.Fatal("tampered patch accepted")
	}
}
func TestFailureRollsBackGeneratedFiles(t *testing.T) {
	e, folder, sources, _ := fixture(t)
	putInputs(t, folder, "English", sources)
	helper := []byte("#!/bin/sh\nexit 7\n")
	p := filepath.Join(e.Root, "xdelta3.exe")
	writeTestFile(t, p, helper)
	e.Decoder = p // Test the execution failure, not the asset-integrity check.
	s, err := e.Scan(folder)
	if err != nil {
		t.Fatal(err)
	}
	if err = e.Apply(s, nil); err == nil {
		t.Fatal("fake decoder failure incorrectly succeeded")
	}
	entries, err := os.ReadDir(folder)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("left partial outputs: %v", entries)
	}
}
func TestExistingKoreanAddsExtrasOnly(t *testing.T) {
	e, folder, _, target := fixture(t)
	for _, ext := range exts {
		writeTestFile(t, filepath.Join(folder, e.Manifest.Target[ext].Filename), target[ext])
	}
	s, err := e.Scan(folder)
	if err != nil {
		t.Fatal(err)
	}
	if s.Edition != "AlreadyPatched" {
		t.Fatal(s.Edition)
	}
	if err = e.Apply(s, nil); err != nil {
		t.Fatal(err)
	}
	assertOutput(t, e, folder)
}
func TestProvidedBundleIntegrity(t *testing.T) {
	if _, err := os.Stat(filepath.Join("assets", manifestFile)); os.IsNotExist(err) {
		t.Skip("배포 자산이 없는 소스 저장소에서는 통합 검증을 건너뜁니다")
	}
	e, err := NewEngine("assets")
	if err != nil {
		t.Fatal(err)
	}
	if e.Manifest.Name != programName {
		t.Fatal(e.Manifest.Name)
	}
	// The user's two provided D88 payloads are byte-identical; keep both requested names.
	h1 := e.Manifest.Extras[1].SHA256
	h2 := e.Manifest.Extras[2].SHA256
	if h1 != h2 {
		t.Fatal("D88 pair unexpectedly differs; review both sources")
	}
	cue, err := os.ReadFile(filepath.Join("assets", "extras", "Mirrors_Kor1.00.cue"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cue), e.Manifest.Target["img"].Filename) {
		t.Fatal("CUE target IMG name mismatch")
	}
}

func TestCanonicalNamesAndCue(t *testing.T) {
    e, _, _, _ := fixture(t)
    for _, ext := range exts {
        want := "Mirrors_Kor1.00." + ext
        if got := e.Manifest.Target[ext].Filename; got != want {
            t.Fatalf("target %s: got %s want %s", ext, got, want)
        }
    }
    found := false
    for _, x := range e.Manifest.Extras {
        if x.Filename == "Mirrors_Kor1.00.cue" { found = true }
        if x.Filename == "Kor.cue" { t.Fatal("legacy CUE remains") }
    }
    if !found { t.Fatal("canonical CUE missing") }
    cue, err := os.ReadFile(filepath.Join(e.Root, "extras", "Mirrors_Kor1.00.cue"))
    if err != nil { t.Fatal(err) }
    if !strings.Contains(string(cue), `FILE "Mirrors_Kor1.00.img" BINARY`) {
        t.Fatal("CUE reference mismatch")
    }
    if strings.Contains(string(cue), "Mirrors_Korean_Mirrors_Tools_Full_Build") {
        t.Fatal("legacy IMG filename in CUE")
    }
}

// Both complete editions in one directory must demand a user-selected edition.
func TestBothEditionsRequireSelection(t *testing.T) {
	for _, chosen := range editions {
		t.Run(chosen, func(t *testing.T) {
			e, folder, sources, _ := fixture(t)
			originalPaths := map[string][]byte{}
			for _, edition := range editions {
				for _, ext := range exts {
					name := edition + "_source." + ext
					data := sources[edition][ext]
					writeTestFile(t, filepath.Join(folder, name), data)
					originalPaths[name] = data
				}
			}
			s, err := e.Scan(folder)
			if err != nil {
				t.Fatal(err)
			}
			if s.Edition != "" || len(s.Options) != 2 {
				t.Fatalf("not offered both editions: %+v", s)
			}
			if err = e.Apply(s, nil); err == nil {
				t.Fatal("accepted unselected edition")
			}
			if _, err = os.Stat(filepath.Join(folder, e.Manifest.Target["img"].Filename)); !os.IsNotExist(err) {
				t.Fatal("unselected apply wrote a result")
			}
			s, err = e.ScanWithEdition(folder, chosen)
			if err != nil {
				t.Fatal(err)
			}
			if s.Edition != chosen {
				t.Fatalf("chose %q but got %q", chosen, s.Edition)
			}
			patchLogs := []string{}
			if err = e.Apply(s, func(line string) { patchLogs = append(patchLogs, line) }); err != nil {
				t.Fatal(err)
			}
			uses := 0
			for _, line := range patchLogs {
				if strings.Contains(line, "전용 xdelta 적용") {
					if !strings.Contains(line, chosen) {
						t.Fatalf("wrong edition patch: %s", line)
					}
					uses++
				}
			}
			if uses != 3 {
				t.Fatalf("expected three %s patches, got %d", chosen, uses)
			}
			assertOutput(t, e, folder)
			for name, orig := range originalPaths {
				got, err := os.ReadFile(filepath.Join(folder, name))
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != string(orig) {
					t.Fatalf("changed original: %s", name)
				}
			}
		})
	}
}

func TestSingleEditionIgnoresIncorrectPreference(t *testing.T) {
	e, folder, sources, _ := fixture(t)
	putInputs(t, folder, "Japanese", sources)
	if _, err := e.ScanWithEdition(folder, "English"); err == nil {
		t.Fatal("accepted unavailable edition")
	}
	s, err := e.Scan(folder)
	if err != nil || s.Edition != "Japanese" {
		t.Fatalf("single edition detection: %+v / %v", s, err)
	}
}
