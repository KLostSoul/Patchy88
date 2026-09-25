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
	"runtime"
	"sort"
	"strings"
)

const (
	programName    = "Mirrors_Kor1.01"
	manifestFile   = "Mirrors_Kor1.01.json"
	manifestSchema = "mirrors.kor1.01.xdelta"
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
	Size     int64  `json:"size"`
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
	Target       map[string]map[string]TargetDef `json:"target"`
	Extras       []ExtraDef                      `json:"extras"`
	XdeltaSHA256 string                          `json:"xdelta3_sha256"`
	XdeltaX64SHA256 string                       `json:"xdelta3_x64_sha256,omitempty"`
}
type Engine struct {
	Root     string
	Manifest Manifest
	Decoder  string
}
type ScanResult struct {
	Folder        string
	Edition       string   // Japanese, English, AlreadyPatched, or empty when both source editions are available
	OutputEdition string   // Source edition of existing canonical outputs
	Options       []string // Complete matching source editions; GUI must prompt when both exist
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
	if len(m.Source) != 2 || len(m.Target) != 2 || len(m.Extras) != 3 {
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
	for _, edition := range editions {
		targets, ok := m.Target[edition]
		if !ok || len(targets) != len(exts) {
			return nil, fmt.Errorf("%s 판본 결과 목록이 불완전합니다", edition)
		}
		for _, ext := range exts {
			t, ok := targets[ext]
			if !ok || !cleanFilename(t.Filename) || t.Filename != programName+"."+ext || t.Size <= 0 ||
				len(t.MD5) != 32 || len(t.SHA256) != 64 {
				return nil, fmt.Errorf("%s %s 결과 이름/크기/해시가 불완전합니다", edition, ext)
			}
			if _, err := hex.DecodeString(t.MD5); err != nil { return nil, fmt.Errorf("%s %s MD5 오류", edition, ext) }
			if _, err := hex.DecodeString(t.SHA256); err != nil { return nil, fmt.Errorf("%s %s SHA-256 오류", edition, ext) }
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
	// Official upstream v3.2.0 release on amd64, official-source Win32 build on 386.
	// A source-only test fixture may omit xdelta3_x64_sha256.
	helperName, expectedHash := "xdelta3.exe", m.XdeltaSHA256
	if runtime.GOARCH == "amd64" && m.XdeltaX64SHA256 != "" {
		helperName, expectedHash = "xdelta3-x64.exe", m.XdeltaX64SHA256
	}
	if len(expectedHash) != 64 {
		return nil, fmt.Errorf("%s SHA-256 값이 올바르지 않습니다", helperName)
	}
	helper := filepath.Join(root, helperName)
	h, err := shaFile(helper)
	if err != nil {
		return nil, fmt.Errorf("%s 읽기 실패: %w", helperName, err)
	}
	if !strings.EqualFold(h, expectedHash) {
		return nil, fmt.Errorf("%s SHA-256이 다릅니다", helperName)
	}
	// Do not require whole-source SHA of any other file: source checks below use both published MD5 and SHA-256.
	return &Engine{Root: root, Manifest: m, Decoder: helper}, nil
}
// Scan identifies available source editions without guessing when both are present.
func (e *Engine) Scan(folder string) (*ScanResult, error) {
	return e.ScanWithEdition(folder, "")
}

// ScanWithEdition allows an explicit choice only among completely verified source sets.
// An empty preference selects automatically only when exactly one edition matches.
func (e *Engine) ScanWithEdition(folder, preferred string) (*ScanResult, error) {
	if preferred != "" && preferred != "Japanese" && preferred != "English" && preferred != "AlreadyPatched" {
		return nil, fmt.Errorf("지원하지 않는 원본 판본: %q", preferred)
	}
	folder, err := filepath.Abs(folder)
	if err != nil { return nil, err }
	info, err := os.Stat(folder)
	if err != nil || !info.IsDir() { return nil, fmt.Errorf("폴더를 열 수 없습니다: %s", folder) }
	items, err := os.ReadDir(folder)
	if err != nil { return nil, err }
	sort.Slice(items, func(i, j int) bool { return strings.ToLower(items[i].Name()) < strings.ToLower(items[j].Name()) })
	candidates := map[string]map[string][]string{}
	outputMatches := map[string]map[string]bool{}
	for _, ed := range editions { candidates[ed]=map[string][]string{}; outputMatches[ed]=map[string]bool{} }
	outputPresent := map[string]bool{}
	canonicalPath := map[string]string{}
	for _, item := range items {
		if item.IsDir() { continue }
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(item.Name()), "."))
		if ext!="ccd" && ext!="img" && ext!="sub" { continue }
		path := filepath.Join(folder, item.Name())
		h, err := hashFile(path)
		if err!=nil { return nil, fmt.Errorf("%s 검사 오류: %w", item.Name(), err) }
		if strings.EqualFold(item.Name(), programName+"."+ext) {
			if outputPresent[ext] { return nil, fmt.Errorf("중복 결과 파일명: %s", item.Name()) }
			outputPresent[ext]=true
			canonicalPath[ext]=path
			matched:=false
			for _, ed:=range editions {
				if _,err:=verifyTarget(path,h,e.Manifest.Target[ed][ext]);err==nil {outputMatches[ed][ext]=true;matched=true}
			}
			if !matched { return nil, fmt.Errorf("다른 내용의 결과 파일 %s가 있습니다. 덮어쓰지 않습니다", item.Name()) }
			continue
		}
		for _, ed:=range editions {
			if hashEqual(h,e.Manifest.Source[ed][ext].Hashes) {
				candidates[ed][ext]=append(candidates[ed][ext],path)
			}
		}
	}
	extrasPresent:=map[string]bool{}
	for _, x:=range e.Manifest.Extras {
		dest:=filepath.Join(folder,x.Filename)
		st,err:=os.Stat(dest)
		if os.IsNotExist(err) { continue }
		if err!=nil {return nil,err}
		if st.IsDir() {return nil,fmt.Errorf("추가 파일명에 폴더가 존재합니다: %s",x.Filename)}
		h,err:=shaFile(dest)
		if err!=nil {return nil,err}
		if st.Size()!=x.Size || !strings.EqualFold(h,x.SHA256) {return nil,fmt.Errorf("다른 내용의 %s 파일이 있습니다",x.Filename)}
		extrasPresent[x.Filename]=true
	}
	// A verified canonical CCD/SUB may double as its source for the SAME edition.
	for _, ed:=range editions {
		for _, ext:=range exts {
			if len(candidates[ed][ext])==0 && outputMatches[ed][ext] {
				src,tgt:=e.Manifest.Source[ed][ext],e.Manifest.Target[ed][ext]
				if strings.EqualFold(src.MD5,tgt.MD5) && strings.EqualFold(src.SHA256,tgt.SHA256) {
					candidates[ed][ext]=append(candidates[ed][ext],canonicalPath[ext])
				}
			}
		}
	}
	completeOutputs:=[]string{}
	for _, ed:=range editions {
		ok:=true
		for _, ext:=range exts {if !outputMatches[ed][ext] {ok=false}}
		if ok {completeOutputs=append(completeOutputs,ed)}
	}
	if len(completeOutputs)>1 {return nil,errors.New("기존 결과가 두 판본에 모두 해당합니다. 원본을 별도 폴더로 분리하세요")}
	if len(completeOutputs)==1 {
		ed:=completeOutputs[0]
		if preferred!="" && preferred!="AlreadyPatched" && preferred!=ed {
			return nil,fmt.Errorf("%s 결과 파일이 이미 있습니다. %s 패치는 별도 폴더에 적용하세요",ed,preferred)
		}
		already:=map[string]bool{}
		for _,ext:=range exts {already[ext]=true}
		return &ScanResult{Folder:folder,Edition:"AlreadyPatched",OutputEdition:ed,
			Inputs:map[string]string{},Already:already,ExtrasPresent:extrasPresent,
			Notes:[]string{ed+" 한글판 CCD/IMG/SUB 해시·크기 검증 통과"}},nil
	}
	if preferred=="AlreadyPatched" {return nil,errors.New("검사 후 한글판 결과가 변경됐습니다")}
	valid:=[]string{}
	for _,ed:=range editions {
		ok:=true
		for _,ext:=range exts {
			if len(candidates[ed][ext])!=1 || (outputPresent[ext]&&!outputMatches[ed][ext]) {ok=false}
		}
		if ok {valid=append(valid,ed)}
	}
	if len(valid)==0 {
		if len(outputPresent)>0 {return nil,errors.New("다른 판본의 결과가 있거나 원본·결과가 혼합되어 있습니다. 별도 폴더를 사용하세요")}
		desc:=[]string{}
		for _,ed:=range editions {
			c:=[]string{}
			for _,ext:=range exts {c=append(c,fmt.Sprintf("%s:%d",strings.ToUpper(ext),len(candidates[ed][ext])))}
			desc=append(desc,ed+" ["+strings.Join(c,", ")+"]")
		}
		return nil,fmt.Errorf("완전한 원본 CCD/IMG/SUB를 찾지 못했습니다. %s",strings.Join(desc,"; "))
	}
	selected:=valid[0]
	if preferred!="" {
		ok:=false
		for _,ed:=range valid {if ed==preferred {ok=true}}
		if !ok {return nil,fmt.Errorf("선택한 %s 판본의 원본이 불완전합니다 (사용 가능: %s)",preferred,strings.Join(valid,", "))}
		selected=preferred
	} else if len(valid)>1 {
		return &ScanResult{Folder:folder,Edition:"",Options:valid,Inputs:map[string]string{},
			Already:map[string]bool{},ExtrasPresent:extrasPresent,
			Notes:[]string{"일본판과 영문판이 모두 확인됐습니다. 원본을 직접 선택하세요."}},nil
	}
	inputs:=map[string]string{}
	already:=map[string]bool{}
	notes:=[]string{fmt.Sprintf("%s 원본 CCD/IMG/SUB 3개 식별 완료",selected)}
	for _,ext:=range exts {
		inputs[ext]=candidates[selected][ext][0]
		already[ext]=outputMatches[selected][ext]
		if already[ext] {notes=append(notes,strings.ToUpper(ext)+": "+selected+" 패치 결과가 이미 검증됨")}
	}
	return &ScanResult{Folder:folder,Edition:selected,OutputEdition:selected,Options:valid,
		Inputs:inputs,Already:already,ExtrasPresent:extrasPresent,Notes:notes},nil
}

type stagedFile struct{ Temp, Dest string }

func tempName(folder string) (string, error) {
	f, err := os.CreateTemp(folder, ".Mirrors_Kor1.01-*.tmp")
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
	if s.Edition == "" {
		return errors.New("일본판과 영문판이 모두 있습니다. 패치할 판본을 먼저 선택하세요")
	}
	// Independently rescan the explicitly selected edition just before applying.
	rescanEdition := s.Edition
	if rescanEdition == "AlreadyPatched" { rescanEdition = s.OutputEdition }
	fresh, err := e.ScanWithEdition(s.Folder, rescanEdition)
	if err != nil {
		return err
	}
	if s.Edition != "AlreadyPatched" && fresh.Edition != s.Edition {
		return errors.New("원본 판본이 검사 이후 변경됐습니다")
	}
	if s.Edition == "AlreadyPatched" && (fresh.Edition != "AlreadyPatched" || fresh.OutputEdition != s.OutputEdition) {
		return errors.New("한글판 결과 판본이 검사 이후 변경됐습니다")
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
	targetEdition := s.Edition
	if targetEdition == "AlreadyPatched" { targetEdition = s.OutputEdition }
	for _, ext := range exts {
		dest := filepath.Join(s.Folder, e.Manifest.Target[targetEdition][ext].Filename)
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
			target := e.Manifest.Target[targetEdition][ext]
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
			method, err := verifyTarget(tmp,got,target)
			if err!=nil { return fmt.Errorf("%s 패치 결과 검증 실패: %w",strings.ToUpper(ext),err) }
			log(fmt.Sprintf("%s: %s 검증 통과 (계산된 MD5=%s, SHA-256=%s)",strings.ToUpper(ext),method,got.MD5,got.SHA256))
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
