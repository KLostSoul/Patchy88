package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// sessionLogger saves the same messages displayed by the Windows GUI.
// One log is created for every application launch, without overwriting older logs.
// A fallback inside the user's config directory supports read-only install folders.
type sessionLogger struct {
	mu     sync.Mutex
	file   *os.File
	path   string
	closed bool
}

func (l *sessionLogger) Path() string { return l.path }

func (l *sessionLogger) WriteLine(message string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return errors.New("로그 파일이 이미 닫혔습니다")
	}
	message = strings.ReplaceAll(message, "\r\n", "\n")
	message = strings.ReplaceAll(message, "\r", "\n")
	var b strings.Builder
	for _, part := range strings.Split(message, "\n") {
		b.WriteString(time.Now().Format("2006-01-02 15:04:05.000"))
		b.WriteString(" ")
		b.WriteString(part)
		b.WriteByte('\n')
	}
	_, err := l.file.WriteString(b.String())
	return err
}

func (l *sessionLogger) Sync() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return errors.New("로그 파일이 이미 닫혔습니다")
	}
	return l.file.Sync()
}

func (l *sessionLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return nil
	}
	l.closed = true
	return l.file.Close()
}

func newSessionLoggerAt(dirs []string, now time.Time) (*sessionLogger, error) {
	var errs []string
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %s", dir, err))
			continue
		}
		for n := 0; n < 1000; n++ {
			name := fmt.Sprintf("PatchyDOS_Dune2_%s", now.Format("20060102_150405"))
			if n > 0 {
				name += fmt.Sprintf("_%03d", n)
			}
			path := filepath.Join(dir, name+".log")
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			if errors.Is(err, os.ErrExist) {
				continue
			}
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %s", path, err))
				break
			}
			return &sessionLogger{file: f, path: path}, nil
		}
	}
	return nil, fmt.Errorf("로그 파일을 생성할 수 없습니다: %s", strings.Join(errs, "; "))
}

func newSessionLogger(executablePath string) (*sessionLogger, error) {
	dirs := []string{filepath.Join(filepath.Dir(executablePath), "logs")}
	if configDir, err := os.UserConfigDir(); err == nil {
		fallback := filepath.Join(configDir, "PatchyDOS", "Dune2", "logs")
		if fallback != dirs[0] {
			dirs = append(dirs, fallback)
		}
	}
	return newSessionLoggerAt(dirs, time.Now())
}
