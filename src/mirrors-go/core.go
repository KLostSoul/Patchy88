package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	programName    = "Mirrors_Kor1.00"
	manifestFile   = "Mirrors_Kor1.00.json"
	manifestSchema = "mirrors.kor1.00.xdelta"
)

var exts = []string{"ccd", "img", "sub"}
var editions = []string{"Japanese", "English"}

type Hashes struct {
	MD5    string `json:"md5"`
	SHA256 string `json:"sha256"`
}
type SourceDef struct {
	Hashes
	Patch       string `json:"patch"`
	PatchSHA256 string `json:"patch_sha256"`
}
type TargetDef struct {
	Hashes
	Filename string `json:"filename"`
}
type ExtraDef struct {
	Filename string `json:"filename"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
}
type Manifest struct {
	Name         string                          `json:"name"`
	Schema       string                          `json:"schema"`
	Source       map[string]map[string]SourceDef `json:"source"`
	Target       map[string]TargetDef            `json:"target"`
	Extras       []ExtraDef                      `json:"extras"`
	XdeltaSHA256 string                          `json:"xdelta3_sha256"`
}
type Engine struct {
	Root     string
	Manifest Manifest
	Decoder  string
}
type ScanResult struct {
	Folder        string
	Edition       string // Japanese, English, or AlreadyPatched
	Inputs        map[string]string
	Already       map[string]bool
	ExtrasPresent map[string]bool
	Notes         []string
}

type fileHashes struct {
	Hashes
	Size int64
}

func hashFile(path string) (fileHashes, error) {
	f, err := os.Open(path)
	if err != nil {
		return fileHashes{}, err
	}
	defer f.Close()
	m, s := md5.New(), sha256.New()
	n, err := io.Copy(io.MultiWriter(m, s), f)
	if err != nil {
		return fileHashes{}, err
	}
	return fileHashes{Hashes: Hashes{MD5: hex.EncodeToString(m.Sum(nil)), SHA256: hex.EncodeToString(s.Sum(nil))}, Size: n}, nil
}
func hashEqual(h fileHashes, ref Hashes) bool {
	return strings.EqualFold(h.MD5, ref.MD5) && strings.EqualFold(h.SHA256, ref.SHA256)
}
func shaFile(path string) (string, error) {
	h, err := hashFile(path)
	if err != nil {
		return "", err
	}
	return h.SHA256, nil
}
func assetPath(root, prefix, name string) (string, error) {
	if name == "" || filepath.IsAbs(name) || strings.Contains(name, "..") || strings.ContainsAny(name, ":\\") {
		return "", fmt.Errorf("잘못된 자산 경로: %q", name)
	}
	p := filepath.Clean(filepath.Join(root, prefix, filepath.FromSlash(name)))
	base := filepath.Clean(filepath.Join(root, prefix))
	rel, err := filepath.Rel(base, p)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("자산 외부 경로: %q", name)
	}
	return p, nil
}
func cleanFilename(s string) bool {
	return s != "" && s == filepath.Base(s) && !strings.ContainsAny(s, "/\\:") && s != "." && s != ".."
}
func NewEngine(root string) (*Engine, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(filepath.Join(root, manifestFile))
	if err != nil {
		return nil, fmt.Errorf("매니페스트 읽기 실패: %w", err)
	}
	var m Manifest
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("매니페스트 형식 오류: %w", err)
	}
	if m.Name != programName || m.Schema != manifestSchema {
		return nil, errors.New("매니페스트 이름/형식이 일치하지 않습니다")
	}
	if len(m.Source) != 2 || len(m.Target) != 3 || len(m.Extras) != 3 {
		return nil, errors.New("매니페스트 대상 개수가 올바르지 않습니다")
	}
	for _, edition := range editions {
		src, ok := m.Source[edition]
		if !ok || len(src) != 3 {
			return nil, fmt.Errorf("%s 패치 목록이 불완전합니다", edition)
		}
		for _, ext := range exts {
			x, ok := src[ext]
			if !ok {
				return nil, fmt.Errorf("%s %s 패치가 없습니다", edition, ext)
			}
			if len(x.MD5) != 32 || len(x.SHA256) != 64 || len(x.PatchSHA256) != 64 {
				return nil, fmt.Errorf("%s %s 해시가 올바르지 않습니다", edition, ext)
			}
			path, err := assetPath(root, "patches", filepath.Base(x.Patch))
			if err != nil {
				return nil, err
			}
			if x.Patch != "patches/"+filepath.Base(x.Patch) {
				return nil, fmt.Errorf("패치 경로가 잘못됐습니다: %s", x.Patch)
			}
			blob, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("%s %s 패치 읽기 실패: %w", edition, ext, err)
			}
			if len(blob) < 4 || blob[0] != 0xd6 || blob[1] != 0xc3 || blob[2] != 0xc4 {
				return nil, fmt.Errorf("%s %s: VCDIFF(xdelta) 형식이 아닙니다", edition, ext)
			}
			got := sha256.Sum256(blob)
			if hex.EncodeToString(got[:]) != strings.ToLower(x.PatchSHA256) {
				return nil, fmt.Errorf("%s %s 패치 SHA-256이 다릅니다", edition, ext)
			}
		}
	}
	for _, ext := range exts {
		t, ok := m.Target[ext]
		if !ok || !cleanFilename(t.Filename) || len(t.MD5) != 32 || len(t.SHA256) != 64 {
			return nil, fmt.Errorf("%s 결과 정보가 잘못됐습니다", ext)
		}
	}
	for _, x := range m.Extras {
		if !cleanFilename(x.Filename) || len(x.SHA256) != 64 || x.Size <= 0 {
			return nil, fmt.Errorf("추가 파일 정보가 잘못됐습니다: %s", x.Filename)
		}
		p, err := assetPath(root, "extras", x.Filename)
		if err != nil {
			return nil, err
		}
		h, err := hashFile(p)
		if err != nil {
			return nil, fmt.Errorf("추가 파일 읽기 실패 %s: %w", x.Filename, err)
		}
		if h.Size != x.Size || !strings.EqualFold(h.SHA256, x.SHA256) {
			return nil, fmt.Errorf("추가 파일 손상: %s", x.Filename)
		}
	}
	helper := filepath.Join(root, "xdelta3.exe")
	h, err := shaFile(helper)
	if err != nil {
		return nil, fmt.Errorf("xdelta3.exe 읽기 실패: %w", err)
	}
	if !strings.EqualFold(h, m.XdeltaSHA256) {
		return nil, errors.New("xdelta3.exe SHA-256이 다릅니다")
	}
	// Do not require whole-source SHA of any other file: source checks below use both published MD5 and SHA-256.
	return &Engine{Root: root, Manifest: m, Decoder: helper}, nil
}
func (e *Engine) Scan(folder string) (*ScanResult, error) {
	folder, err := filepath.Abs(folder)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(folder)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("폴더를 열 수 없습니다: %s", folder)
	}
	items, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return strings.ToLower(items[i].Name()) < strings.ToLower(items[j].Name()) })
	candidates := map[string]map[string][]string{}
	for _, edition := range editions {
		candidates[edition] = map[string][]string{}
	}
	already := map[string]bool{}
	canonical := map[string]string{}
	targetCandidates := map[string][]string{}
	for _, ext := range exts {
		canonical[ext] = e.Manifest.Target[ext].Filename
	}
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(item.Name()), "."))
		if ext != "ccd" && ext != "img" && ext != "sub" {
			continue
		}
		path := filepath.Join(folder, item.Name())
		h, err := hashFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s 검사 오류: %w", item.Name(), err)
		}
		if strings.EqualFold(item.Name(), canonical[ext]) {
			if !hashEqual(h, e.Manifest.Target[ext].Hashes) {
				return nil, fmt.Errorf("결과 파일명 %s에 예상과 다른 데이터가 이미 있습니다. 덮어쓰지 않습니다", item.Name())
			}
			already[ext] = true
			targetCandidates[ext] = append(targetCandidates[ext], path)
			// A Japanese CCD/SUB is byte-identical to the corresponding final Korean file;
			// allow canonical names as fallbacks only if no other Japanese source exists.
			continue
		}
		for _, edition := range editions {
			def := e.Manifest.Source[edition][ext]
			if hashEqual(h, def.Hashes) {
				candidates[edition][ext] = append(candidates[edition][ext], path)
			}
		}
	}
	extrasPresent := map[string]bool{}
	for _, x := range e.Manifest.Extras {
		dest := filepath.Join(folder, x.Filename)
		s, err := os.Stat(dest)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if s.IsDir() {
			return nil, fmt.Errorf("추가 파일명에 폴더가 존재합니다: %s", x.Filename)
		}
		h, err := shaFile(dest)
		if err != nil {
			return nil, err
		}
		if s.Size() != x.Size || !strings.EqualFold(h, x.SHA256) {
			return nil, fmt.Errorf("다른 내용의 %s 파일이 이미 있습니다. 덮어쓰지 않습니다", x.Filename)
		}
		extrasPresent[x.Filename] = true
	}
	allAlready := true
	for _, ext := range exts {
		if !already[ext] {
			allAlready = false
		}
	}
	if allAlready {
		return &ScanResult{Folder: folder, Edition: "AlreadyPatched", Inputs: map[string]string{}, Already: already, ExtrasPresent: extrasPresent, Notes: []string{"한글판 CCD/IMG/SUB가 이미 모두 있고 해시가 일치합니다."}}, nil
	}
	// The Japanese CCD/SUB hashes equal the Korean ones, so a single canonical file can act as both.
	for _, ext := range exts {
		if len(candidates["Japanese"][ext]) == 0 && already[ext] && len(targetCandidates[ext]) == 1 {
			if e.Manifest.Source["Japanese"][ext].MD5 == e.Manifest.Target[ext].MD5 && e.Manifest.Source["Japanese"][ext].SHA256 == e.Manifest.Target[ext].SHA256 {
				candidates["Japanese"][ext] = append(candidates["Japanese"][ext], targetCandidates[ext][0])
			}
		}
	}
	valid := []string{}
	for _, edition := range editions {
		complete := true
		for _, ext := range exts {
			if len(candidates[edition][ext]) != 1 {
				complete = false
			}
		}
		if complete {
			valid = append(valid, edition)
		}
	}
	if len(valid) != 1 {
		desc := []string{}
		for _, edition := range editions {
			counts := []string{}
			for _, ext := range exts {
				counts = append(counts, fmt.Sprintf("%s:%d", strings.ToUpper(ext), len(candidates[edition][ext])))
			}
			desc = append(desc, edition+" ["+strings.Join(counts, ", ")+"]")
		}
		return nil, fmt.Errorf("일본판 또는 영문판 원본 3개를 유일하게 식별할 수 없습니다. 파일 누락/중복/판본 혼합 여부를 확인하세요. %s", strings.Join(desc, "; "))
	}
	edition := valid[0]
	inputs := map[string]string{}
	for _, ext := range exts {
		inputs[ext] = candidates[edition][ext][0]
	}
	notes := []string{fmt.Sprintf("%s 원본 CCD/IMG/SUB 3개 식별 완료 (MD5 + SHA-256)", edition)}
	for _, ext := range exts {
		if already[ext] {
			notes = append(notes, strings.ToUpper(ext)+" 결과 파일은 이미 검증됨")
		}
	}
	return &ScanResult{Folder: folder, Edition: edition, Inputs: inputs, Already: already, ExtrasPresent: extrasPresent, Notes: notes}, nil
}

type stagedFile struct{ Temp, Dest string }

func tempName(folder string) (string, error) {
	f, err := os.CreateTemp(folder, ".Mirrors_Kor1.00-*.tmp")
	if err != nil {
		return "", err
	}
	p := f.Name()
	if err = f.Close(); err != nil {
		os.Remove(p)
		return "", err
	}
	if err = os.Remove(p); err != nil {
		return "", err
	}
	return p, nil
}
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, cpy := io.Copy(out, in)
	if cpy == nil {
		cpy = out.Sync()
	}
	closeErr := out.Close()
	if cpy != nil {
		os.Remove(dst)
		return cpy
	}
	if closeErr != nil {
		os.Remove(dst)
		return closeErr
	}
	return nil
}
func (e *Engine) Apply(s *ScanResult, log func(string)) error {
	if log == nil {
		log = func(string) {}
	}
	if s == nil {
		return errors.New("먼저 폴더를 검사해야 합니다")
	}
	// Independently rescan immediately before applying; no earlier GUI scan result is trusted.
	fresh, err := e.Scan(s.Folder)
	if err != nil {
		return err
	}
	if s.Edition != "AlreadyPatched" && fresh.Edition != s.Edition {
		return errors.New("원본 판본이 검사 이후 변경됐습니다")
	}
	s = fresh
	staged := []stagedFile{}
	committed := []string{}
	success := false
	defer func() {
		for _, f := range staged {
			os.Remove(f.Temp)
		}
		if !success {
			for i := len(committed) - 1; i >= 0; i-- {
				os.Remove(committed[i])
			}
		}
	}()
	// All output names are reserved/checked before decoding anything.
	for _, ext := range exts {
		dest := filepath.Join(s.Folder, e.Manifest.Target[ext].Filename)
		if s.Already[ext] {
			continue
		}
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("결과 파일이 이미 있습니다: %s", dest)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	for _, x := range e.Manifest.Extras {
		dest := filepath.Join(s.Folder, x.Filename)
		if s.ExtrasPresent[x.Filename] {
			continue
		}
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("추가 파일이 이미 있습니다: %s", dest)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	if s.Edition != "AlreadyPatched" {
		for _, ext := range exts {
			target := e.Manifest.Target[ext]
			if s.Already[ext] {
				log(strings.ToUpper(ext) + ": 이미 올바른 결과 파일이 있습니다")
				continue
			}
			def := e.Manifest.Source[s.Edition][ext]
			src := s.Inputs[ext]
			current, err := hashFile(src)
			if err != nil {
				return err
			}
			if !hashEqual(current, def.Hashes) {
				return fmt.Errorf("%s 원본 데이터가 검사 후 변경됐습니다", strings.ToUpper(ext))
			}
			patch := filepath.Join(e.Root, filepath.FromSlash(def.Patch))
			tmp, err := tempName(s.Folder)
			if err != nil {
				return err
			}
			staged = append(staged, stagedFile{Temp: tmp, Dest: filepath.Join(s.Folder, target.Filename)})
			log(fmt.Sprintf("%s: %s 전용 xdelta 적용", strings.ToUpper(ext), s.Edition))
			cmd := exec.Command(e.Decoder, "-q", "-d", "-s", src, patch, tmp)
			configureCommand(cmd)
			output, err := cmd.CombinedOutput()
			if err != nil {
				message := strings.TrimSpace(string(output))
				if len(message) > 600 {
					message = message[:600]
				}
				return fmt.Errorf("%s xdelta 복호화 실패: %w; %s", strings.ToUpper(ext), err, message)
			}
			got, err := hashFile(tmp)
			if err != nil {
				return err
			}
			if !hashEqual(got, target.Hashes) {
				return fmt.Errorf("%s 패치 결과 MD5/SHA-256 검증 실패: MD5=%s, SHA-256=%s", strings.ToUpper(ext), got.MD5, got.SHA256)
			}
			log(strings.ToUpper(ext) + ": 결과 MD5/SHA-256 검증 통과")
		}
	} else {
		log("이미 한글판 파일 3개가 검증됐습니다. 누락된 추가 파일만 복사합니다.")
	}
	for _, x := range e.Manifest.Extras {
		if s.ExtrasPresent[x.Filename] {
			log(x.Filename + ": 이미 존재하며 해시가 일치합니다")
			continue
		}
		tmp, err := tempName(s.Folder)
		if err != nil {
			return err
		}
		staged = append(staged, stagedFile{Temp: tmp, Dest: filepath.Join(s.Folder, x.Filename)})
		asset := filepath.Join(e.Root, "extras", x.Filename)
		if err = copyFile(asset, tmp); err != nil {
			return fmt.Errorf("%s 추가 파일 복사 실패: %w", x.Filename, err)
		}
		h, err := hashFile(tmp)
		if err != nil {
			return err
		}
		if h.Size != x.Size || !strings.EqualFold(h.SHA256, x.SHA256) {
			return fmt.Errorf("%s 임시 복사 검증 실패", x.Filename)
		}
		log(x.Filename + ": 추가 파일 복사 검증 통과")
	}
	// No original is overwritten: conflicts are rechecked immediately before commit.
	for _, f := range staged {
		if _, err := os.Stat(f.Dest); err == nil {
			return fmt.Errorf("확정 직전 출력 파일 충돌: %s", f.Dest)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	for _, f := range staged {
		if err = os.Rename(f.Temp, f.Dest); err != nil {
			return fmt.Errorf("결과 파일 확정 실패 (%s): %w", filepath.Base(f.Dest), err)
		}
		committed = append(committed, f.Dest)
		log("완료: " + filepath.Base(f.Dest))
	}
	success = true
	return nil
}
