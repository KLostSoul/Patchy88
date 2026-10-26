package main

import (
	"context"
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
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	programName    = "PatchyDOS Dune II v1.01"
	manifestFile   = "Dune2_Kor_v1.01.json"
	expectedSchema = "patchydos.dune2.xdelta.v1"
	expectedName   = "Dune II Battle for Arrakis Kor v1.01"
	expectedBIN    = "Dune2-Battle-for-Arrakis-Korean-v1.01.bin"
	expectedCUE    = "Dune2-Battle-for-Arrakis-Korean-v1.01.cue"
)

type Hashes struct {
	MD5    string `json:"md5"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}
type OutputHashes struct {
	Hashes
	Filename string `json:"filename"`
}
type PatchHashes struct {
	Filename string `json:"filename"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
}
type Manifest struct {
	Name             string       `json:"name"`
	Schema           string       `json:"schema"`
	Source           Hashes       `json:"source"`
	Target           OutputHashes `json:"target"`
	Cue              OutputHashes `json:"cue"`
	Patch            PatchHashes  `json:"patch"`
	Xdelta3SHA256    string       `json:"xdelta3_sha256"`
	Xdelta3X64SHA256 string       `json:"xdelta3_x64_sha256"`
}
type Engine struct {
	Root        string
	Manifest    Manifest
	Decoder     string
	Patch       string
	CueTemplate []byte
}

func hashFile(path string) (Hashes, error) {
	f, err := os.Open(path)
	if err != nil {
		return Hashes{}, err
	}
	defer f.Close()
	h1, h2 := md5.New(), sha256.New()
	sz, err := io.Copy(io.MultiWriter(h1, h2), f)
	if err != nil {
		return Hashes{}, err
	}
	return Hashes{MD5: hex.EncodeToString(h1.Sum(nil)), SHA256: hex.EncodeToString(h2.Sum(nil)), Size: sz}, nil
}
func hashBytes(data []byte) Hashes {
	a, b := md5.Sum(data), sha256.Sum256(data)
	return Hashes{MD5: hex.EncodeToString(a[:]), SHA256: hex.EncodeToString(b[:]), Size: int64(len(data))}
}
func equalHashes(a, b Hashes) bool {
	return a.Size == b.Size && strings.EqualFold(a.MD5, b.MD5) && strings.EqualFold(a.SHA256, b.SHA256)
}
func NewEngine(root string) (*Engine, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(abs, manifestFile))
	if err != nil {
		return nil, fmt.Errorf("매니페스트 읽기 실패: %w", err)
	}
	var m Manifest
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("매니페스트 JSON 오류: %w", err)
	}
	if m.Name != expectedName || m.Schema != expectedSchema || m.Target.Filename != expectedBIN || m.Cue.Filename != expectedCUE || m.Patch.Filename != "Dune2-Battle-for-Arrakis-Korean-v1.01.bin.xdelta" {
		return nil, errors.New("매니페스트 이름 또는 형식이 올바르지 않습니다")
	}
	if m.Source.Size <= 0 || m.Target.Size <= 0 || m.Cue.Size <= 0 || m.Patch.Size <= 0 {
		return nil, errors.New("매니페스트의 파일 크기가 잘못됐습니다")
	}
	patch := filepath.Join(abs, "patches", m.Patch.Filename)
	ph, err := hashFile(patch)
	if err != nil {
		return nil, fmt.Errorf("xdelta 읽기 실패: %w", err)
	}
	if ph.Size != m.Patch.Size || !strings.EqualFold(ph.SHA256, m.Patch.SHA256) {
		return nil, errors.New("xdelta 크기 또는 SHA-256 불일치")
	}
	cuePath := filepath.Join(abs, m.Cue.Filename)
	cue, err := os.ReadFile(cuePath)
	if err != nil {
		return nil, fmt.Errorf("결과 CUE 템플릿 읽기 실패: %w", err)
	}
	if !equalHashes(hashBytes(cue), m.Cue.Hashes) {
		return nil, errors.New("결과 CUE 템플릿 해시 불일치")
	}
	helper := "xdelta3.exe"
	expectedHelper := m.Xdelta3SHA256
	if runtime.GOARCH == "amd64" {
		helper = "xdelta3-x64.exe"
		expectedHelper = m.Xdelta3X64SHA256
	}
	decoder := filepath.Join(abs, helper)
	dh, err := hashFile(decoder)
	if err != nil {
		return nil, fmt.Errorf("xdelta3 실행파일 읽기 실패: %w", err)
	}
	if !strings.EqualFold(dh.SHA256, expectedHelper) {
		return nil, fmt.Errorf("%s 검증 실패", helper)
	}
	return &Engine{Root: abs, Manifest: m, Decoder: decoder, Patch: patch, CueTemplate: cue}, nil
}

// FindSource walks the selected folder (including its subfolders). It recognizes
// the original raw MODE1/2352 data track by the *full* size, MD5 and SHA-256,
// regardless of file name. No cue/ccd parsing is needed: these files only point
// at the raw track; the xdelta source must be the raw IMG/BIN itself.
// Symlinked files or directories are never followed.
func (e *Engine) FindSource(folder string) (string, error) {
	abs, err := filepath.Abs(folder)
	if err != nil {
		return "", err
	}
	st, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("원본 폴더를 열 수 없습니다: %w", err)
	}
	if !st.IsDir() {
		return "", errors.New("파일이 아니라 원본 폴더를 지정해야 합니다")
	}
	var matches []string
	var candidates int
	err = filepath.WalkDir(abs, func(path string, entry os.DirEntry, readErr error) error {
		if readErr != nil {
			return fmt.Errorf("폴더 검색 실패 (%s): %w", path, readErr)
		}
		if !entry.Type().IsRegular() {
			return nil
		} // ignore symlinks, directories and special files
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".img" && ext != ".bin" {
			return nil
		}
		st, err := entry.Info()
		if err != nil {
			return err
		}
		if st.Size() != e.Manifest.Source.Size {
			return nil
		}
		candidates++
		h, err := hashFile(path)
		if err != nil {
			return fmt.Errorf("후보 파일 읽기 실패 (%s): %w", path, err)
		}
		if equalHashes(h, e.Manifest.Source) {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("선택한 폴더에서 지원 원본 IMG/BIN을 찾지 못했습니다 (크기 일치 후보: %d개).\n원본 이미지의 내용 전체가 요구된 해시와 일치해야 합니다", candidates)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("동일한 원본 이미지가 %d개 발견됐습니다. 패치할 원본이 하나만 포함된 하위 폴더를 선택하세요:\n%s", len(matches), strings.Join(matches, "\n"))
	}
}

// ValidateSource also re-verifies the source just before patching. A changed
// original must never be accepted merely because it passed the folder scan.
func (e *Engine) ValidateSource(path string) error {
	h, err := hashFile(path)
	if err != nil {
		return fmt.Errorf("원본 IMG/BIN 읽기 실패: %w", err)
	}
	if !equalHashes(h, e.Manifest.Source) {
		return fmt.Errorf("원본이 검사 후 변경됐거나 지원 버전이 아닙니다\n크기: %d\nMD5: %s\nSHA-256: %s", h.Size, h.MD5, h.SHA256)
	}
	return nil
}
func (e *Engine) ValidateBIN(path string) error {
	h, err := hashFile(path)
	if err != nil {
		return err
	}
	if !equalHashes(h, e.Manifest.Target.Hashes) {
		return fmt.Errorf("한글 BIN 검증 실패: 크기 %d / MD5 %s / SHA-256 %s", h.Size, h.MD5, h.SHA256)
	}
	return nil
}
func (e *Engine) ValidateCUE(path string) error {
	h, err := hashFile(path)
	if err != nil {
		return err
	}
	if !equalHashes(h, e.Manifest.Cue.Hashes) {
		return errors.New("같은 이름의 CUE가 있으나 검증값이 다릅니다")
	}
	return nil
}
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func stageCUE(dirname string, contents []byte) (string, error) {
	f, err := os.CreateTemp(dirname, "PatchyDOS-Dune2-v1.01-*.cue.tmp")
	if err != nil {
		return "", err
	}
	tmp := f.Name()
	if _, err = f.Write(contents); err != nil {
		f.Close()
		os.Remove(tmp)
		return "", err
	}
	if err = f.Close(); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return tmp, nil
}

// Apply only ever publishes verified results. Existing mismatched target files are not overwritten.
func (e *Engine) Apply(folder string, logf func(string)) (string, string, error) {
	return e.ApplyContext(context.Background(), folder, logf)
}

// ApplyContext supports user cancellation, a decoder timeout and progress monitoring.
func (e *Engine) ApplyContext(ctx context.Context, folder string, logf func(string)) (string, string, error) {
	src, err := e.FindSource(folder)
	if err != nil {
		return "", "", err
	}
	if logf != nil {
		logf("[자동 인식] " + src)
	}
	if err := ctx.Err(); err != nil {
		return "", "", fmt.Errorf("작업이 중단되었습니다: %w", err)
	}
	if err := e.ValidateSource(src); err != nil {
		return "", "", err
	}
	outDir := filepath.Dir(src)
	outBIN := filepath.Join(outDir, e.Manifest.Target.Filename)
	outCUE := filepath.Join(outDir, e.Manifest.Cue.Filename)
	binExists, err := fileExists(outBIN)
	if err != nil {
		return "", "", err
	}
	cueExists, err := fileExists(outCUE)
	if err != nil {
		return "", "", err
	}
	if binExists {
		if err = e.ValidateBIN(outBIN); err != nil {
			return "", "", fmt.Errorf("기존 BIN을 덮어쓰지 않습니다: %w", err)
		}
	}
	if cueExists {
		if err = e.ValidateCUE(outCUE); err != nil {
			return "", "", fmt.Errorf("기존 CUE를 덮어쓰지 않습니다: %w", err)
		}
	}
	if binExists && cueExists {
		if logf != nil {
			logf("[확인] 기존 BIN/CUE가 검증값과 일치합니다.")
		}
		return outBIN, outCUE, nil
	}
	tmpBIN := ""
	if !binExists {
		f, er := os.CreateTemp(outDir, "PatchyDOS-Dune2-v1.01-*.bin.tmp")
		if er != nil {
			return "", "", er
		}
		tmpBIN = f.Name()
		f.Close()
		os.Remove(tmpBIN)
		defer os.Remove(tmpBIN)
		if logf != nil {
			logf("원본 CD IMG/BIN에서 xdelta 적용 중...")
		}
		if er := e.runDecoder(ctx, src, tmpBIN, logf); er != nil {
			return "", "", er
		}
		if logf != nil {
			logf("생성 BIN 전체 크기/MD5/SHA-256 검증 중...")
		}
		if er = e.ValidateBIN(tmpBIN); er != nil {
			return "", "", er
		}
	}
	tmpCUE := ""
	if !cueExists {
		tmpCUE, err = stageCUE(outDir, e.CueTemplate)
		if err != nil {
			return "", "", err
		}
		defer os.Remove(tmpCUE)
		if err = e.ValidateCUE(tmpCUE); err != nil {
			return "", "", err
		}
	}
	// Avoid overwrites if another process created an output during processing.
	if !binExists {
		exists, er := fileExists(outBIN)
		if er != nil {
			return "", "", er
		}
		if exists {
			return "", "", errors.New("패치 진행 중 결과 BIN이 다른 프로세스에서 생성되었습니다")
		}
	}
	if !cueExists {
		exists, er := fileExists(outCUE)
		if er != nil {
			return "", "", er
		}
		if exists {
			return "", "", errors.New("패치 진행 중 결과 CUE가 다른 프로세스에서 생성되었습니다")
		}
	}
	// A late GUI cancellation must not publish a result after the decoder exits.
	if err := ctx.Err(); err != nil {
		return "", "", fmt.Errorf("사용자 요청으로 결과 확정 전에 중단했습니다: %w", err)
	}
	if !binExists {
		if err = os.Rename(tmpBIN, outBIN); err != nil {
			return "", "", fmt.Errorf("검증된 BIN 확정 실패: %w", err)
		}
	}
	if !cueExists {
		if err = os.Rename(tmpCUE, outCUE); err != nil {
			return "", "", fmt.Errorf("검증된 CUE 확정 실패: %w", err)
		}
	}
	return outBIN, outCUE, nil
}

// The 5 minute deadline is a safety bound, not an expected runtime. If no
// output growth and no diagnostic output occur for 90 seconds, the decoder
// is considered unresponsive. Keep values overridable in isolated tests.
var (
	decoderTotalTimeout     = 5 * time.Minute
	decoderStallTimeout     = 90 * time.Second
	decoderPollInterval     = 2 * time.Second
	decoderProgressInterval = 5 * time.Second
)

// diagnosticWriter records bounded stdout/stderr and forwards complete lines.
// xdelta3 usually writes diagnostics to stderr, even in successful runs.
type diagnosticWriter struct {
	mu       sync.Mutex
	prefix   string
	buffer   string
	tail     string
	logf     func(string)
	activity *atomic.Int64
}

func (w *diagnosticWriter) Write(data []byte) (int, error) {
	w.activity.Store(time.Now().UnixNano())
	w.mu.Lock()
	w.tail += strings.ToValidUTF8(string(data), "?")
	if len(w.tail) > 16384 {
		w.tail = w.tail[len(w.tail)-16384:]
	}
	w.buffer += strings.ToValidUTF8(string(data), "?")
	var lines []string
	for {
		idx := strings.IndexByte(w.buffer, '\n')
		if idx < 0 {
			break
		}
		line := strings.TrimSpace(w.buffer[:idx])
		w.buffer = w.buffer[idx+1:]
		if len(line) > 1024 {
			line = line[:1024] + "…"
		}
		if line != "" {
			lines = append(lines, w.prefix+line)
		}
	}
	if len(w.buffer) > 4096 {
		lines = append(lines, w.prefix+w.buffer[:1024]+"…")
		w.buffer = ""
	}
	w.mu.Unlock()
	if w.logf != nil {
		for _, line := range lines {
			w.logf(line)
		}
	}
	return len(data), nil
}
func (w *diagnosticWriter) summary() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	if strings.TrimSpace(w.tail) == "" {
		return "(디코더 진단 출력 없음)"
	}
	return strings.TrimSpace(w.tail)
}
func (w *diagnosticWriter) flush() {
	w.mu.Lock()
	s := strings.TrimSpace(w.buffer)
	w.buffer = ""
	w.mu.Unlock()
	if s != "" && w.logf != nil {
		if len(s) > 1024 {
			s = s[:1024] + "…"
		}
		w.logf(w.prefix + s)
	}
}
func fmtBytes(n int64) string {
	if n < 0 {
		return "아직 생성되지 않음"
	}
	return fmt.Sprintf("%d bytes (%.2f MiB)", n, float64(n)/1048576.0)
}
func partialSize(path string) int64 {
	st, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return st.Size()
}

func (e *Engine) runDecoder(ctx context.Context, src, tmpBIN string, logf func(string)) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("사용자 요청으로 패치 중단: %w", err)
	}
	totalCtx, cancel := context.WithTimeout(ctx, decoderTotalTimeout)
	defer cancel()
	var activity atomic.Int64
	start := time.Now()
	activity.Store(start.UnixNano())
	errWriter := &diagnosticWriter{prefix: "[xdelta3 오류/진단] ", logf: logf, activity: &activity}
	outWriter := &diagnosticWriter{prefix: "[xdelta3 출력] ", logf: logf, activity: &activity}
	cmd := exec.CommandContext(totalCtx, e.Decoder, "-d", "-s", src, e.Patch, tmpBIN)
	cmd.Dir = e.Root
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	// If the child inherits a pipe into another process, do not hang forever.
	cmd.WaitDelay = 3 * time.Second
	if logf != nil {
		logf(fmt.Sprintf("[xdelta3] 실행 시작 / 전체 제한 %s / 무응답 판정 %s", decoderTotalTimeout, decoderStallTimeout))
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("xdelta3 실행 실패: %w (보안 프로그램 차단 또는 실행 권한 확인)", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	tick := time.NewTicker(decoderPollInterval)
	defer tick.Stop()
	previous := partialSize(tmpBIN)
	lastProgress := start
	nextReport := start
	stalled := false
	for {
		select {
		case err := <-done:
			errWriter.flush()
			outWriter.flush()
			elapsed := time.Since(start).Round(time.Second)
			size := partialSize(tmpBIN)
			if logf != nil {
				logf(fmt.Sprintf("[xdelta3] 종료 / 경과 %s / 임시 BIN %s", elapsed, fmtBytes(size)))
			}
			if stalled {
				return fmt.Errorf("xdelta3 응답 없음: %s 동안 출력 파일 증가나 진단 출력이 없어 자동 중단했습니다. 원본 폴더의 쓰기 권한, 보안 프로그램의 차단 및 저장 공간을 확인하세요. 임시 파일은 정리합니다.\n진단: %s", decoderStallTimeout, errWriter.summary())
			}
			if ctx.Err() != nil {
				return fmt.Errorf("사용자 요청으로 작업을 중단했습니다. 임시 파일은 정리합니다: %w", ctx.Err())
			}
			if totalCtx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("xdelta3 실행 시간 제한(%s)을 초과해 중단했습니다. 임시 파일은 정리합니다.\n진단: %s", decoderTotalTimeout, errWriter.summary())
			}
			if err != nil {
				return fmt.Errorf("xdelta3 적용 실패 (경과 %s, 임시 BIN %s): %w\n[xdelta3 stderr] %s\n[xdelta3 stdout] %s", elapsed, fmtBytes(size), err, errWriter.summary(), outWriter.summary())
			}
			if logf != nil {
				logf("[xdelta3] 정상 종료; 결과 이미지 무결성을 검사합니다.")
			}
			return nil
		case now := <-tick.C:
			sz := partialSize(tmpBIN)
			if sz != previous {
				lastProgress = now
				previous = sz
			}
			lastOutput := time.Unix(0, activity.Load())
			idle := now.Sub(lastProgress)
			if oidle := now.Sub(lastOutput); oidle < idle {
				idle = oidle
			}
			if now.Sub(nextReport) >= decoderProgressInterval {
				if logf != nil {
					logf(fmt.Sprintf("[진행] 경과 %s / 임시 BIN %s / 마지막 활동 %s 전", now.Sub(start).Round(time.Second), fmtBytes(sz), idle.Round(time.Second)))
				}
				nextReport = now
			}
			if idle >= decoderStallTimeout && totalCtx.Err() == nil {
				stalled = true
				if logf != nil {
					logf("[중단] xdelta3 무응답 감지. 프로세스를 종료하고 임시 파일을 정리합니다.")
				}
				cancel()
			}
		}
	}
}
