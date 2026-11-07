package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const maxSize = 5 << 20

type Logger struct {
	file *os.File
	path string
}

// Open never fails: a log that cannot be opened discards writes.
func Open() *Logger {
	path, err := defaultPath()
	if err != nil {
		return &Logger{}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return &Logger{}
	}

	rotate(path)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return &Logger{}
	}
	return &Logger{file: f, path: path}
}

func defaultPath() (string, error) {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return filepath.Join(dir, "dotfiles", "dotfiles.log"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "state", "dotfiles", "dotfiles.log"), nil
}

func rotate(path string) {
	info, err := os.Stat(path)
	if err != nil || info.Size() < maxSize {
		return
	}
	_ = os.Rename(path, path+".1")
}

func (l *Logger) Path() string {
	if l == nil {
		return ""
	}
	return l.path
}

func (l *Logger) Printf(format string, args ...any) {
	if l == nil || l.file == nil {
		return
	}
	_, _ = fmt.Fprintf(l.file, "[%s] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
}

func (l *Logger) Close() {
	if l == nil || l.file == nil {
		return
	}
	_ = l.file.Close()
}
