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
	"strings"
)

const (
	programName  = "PatchFM Alltynex v1.0"
	manifestFile = "Alltynex_Kor_v1.0.json"
)

type Hashes struct {
	MD5    string `json:"md5"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

type Manifest struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
	Source Hashes `json:"source"`
	Target struct {
		Hashes
		Filename string `json:"filename"`
	} `json:"target"`
	Patch struct {
		Filename string `json:"filename"`
		SHA256   string `json:"sha256"`
	} `json:"patch"`
	Xdelta3SHA256    string `json:"xdelta3_sha256"`
	Xdelta3X64SHA256 string `json:"xdelta3_x64_sha256"`
}

type Engine struct {
	Root     string
	Manifest Manifest
	Decoder  string
	Patch    string
}

type FileHashes struct{ Hashes }

func hashFile(path string) (FileHashes, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileHashes{}, err
	}
	defer f.Close()
	m, s := md5.New(), sha256.New()
	n, err := io.Copy(io.MultiWriter(m, s), f)
	if err != nil {
		return FileHashes{}, err
	}
	return FileHashes{Hashes{MD5: hex.EncodeToString(m.Sum(nil)), SHA256: hex.EncodeToString(s.Sum(nil)), Size: n}}, nil
}

func equalHashes(a, b Hashes) bool {
	return a.Size == b.Size && strings.EqualFold(a.MD5, b.MD5) && strings.EqualFold(a.SHA256, b.SHA256)
}

func NewEngine(root string) (*Engine, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(filepath.Join(abs, manifestFile))
	if err != nil {
		return nil, fmt.Errorf("매니페스트 읽기 실패: %w", err)
	}
	var m Manifest
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("매니페스트 형식 오류: %w", err)
	}
	if m.Name != "Alltynex Kor v1.0" || m.Schema != "patchfm.alltynex.xdelta.v1" {
		return nil, errors.New("매니페스트 이름/형식이 올바르지 않습니다")
	}
	if m.Target.Filename != "Alltynex (Kor v1.0).iso" {
		return nil, errors.New("결과 파일명이 올바르지 않습니다")
	}
	patch := filepath.Join(abs, "patches", m.Patch.Filename)
	ph, err := hashFile(patch)
	if err != nil {
		return nil, fmt.Errorf("xdelta 읽기 실패: %w", err)
	}
	if !strings.EqualFold(ph.SHA256, m.Patch.SHA256) {
		return nil, errors.New("xdelta SHA-256이 다릅니다")
	}
	helper := "xdelta3.exe"
	want := m.Xdelta3SHA256
	if runtime.GOARCH == "amd64" {
		helper = "xdelta3-x64.exe"
		want = m.Xdelta3X64SHA256
	}
	decoder := filepath.Join(abs, helper)
	dh, err := hashFile(decoder)
	if err != nil {
		return nil, fmt.Errorf("xdelta3 읽기 실패: %w", err)
	}
	if !strings.EqualFold(dh.SHA256, want) {
		return nil, fmt.Errorf("%s SHA-256이 다릅니다", helper)
	}
	return &Engine{Root: abs, Manifest: m, Decoder: decoder, Patch: patch}, nil
}

func (e *Engine) ValidateSource(path string) error {
	h, err := hashFile(path)
	if err != nil {
		return err
	}
	if !equalHashes(h.Hashes, e.Manifest.Source) {
		return fmt.Errorf("지원 원본이 아닙니다\n크기: %d\nMD5: %s\nSHA-256: %s", h.Size, h.MD5, h.SHA256)
	}
	return nil
}

func (e *Engine) ValidateTarget(path string) error {
	h, err := hashFile(path)
	if err != nil {
		return err
	}
	if !equalHashes(h.Hashes, e.Manifest.Target.Hashes) {
		return fmt.Errorf("결과 ISO 검증 실패\n크기: %d\nMD5: %s\nSHA-256: %s", h.Size, h.MD5, h.SHA256)
	}
	return nil
}

func (e *Engine) Apply(source string, logf func(string)) (string, error) {
	if err := e.ValidateSource(source); err != nil {
		return "", err
	}
	out := filepath.Join(filepath.Dir(source), e.Manifest.Target.Filename)
	if _, err := os.Stat(out); err == nil {
		if e.ValidateTarget(out) == nil {
			return out, nil
		}
		return "", fmt.Errorf("기존 결과 파일이 있으나 검증값이 다릅니다: %s", out)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	tmp, err := os.CreateTemp(filepath.Dir(source), "PatchFM-Alltynex-*.iso.tmp")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)
	defer os.Remove(tmpPath)
	if logf != nil {
		logf("xdelta 적용 중...")
	}
	cmd := exec.Command(e.Decoder, "-d", "-s", source, e.Patch, tmpPath)
	cmd.Dir = e.Root
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("xdelta 적용 실패: %v\n%s", err, strings.TrimSpace(string(data)))
	}
	if logf != nil {
		logf("생성 ISO 전체 MD5/SHA-256 검증 중...")
	}
	if err = e.ValidateTarget(tmpPath); err != nil {
		return "", err
	}
	if err = os.Rename(tmpPath, out); err != nil {
		return "", fmt.Errorf("결과 확정 실패: %w", err)
	}
	return out, nil
}
