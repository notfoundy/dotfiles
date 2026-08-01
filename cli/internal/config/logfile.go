// The persistent, append-only run log.
//
// The bash implementation truncated ~/.dotfiles.log on every single command,
// which left nothing to inspect after a failure. This log lives under the XDG
// state directory and is only rotated once it grows past maxSize.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const maxSize = 5 << 20 // 5 MiB

// Logger appends timestamped lines to the run log. A nil-safe zero value is
// never produced: use Open, and keep using the returned logger even on error.
type Logger struct {
	file *os.File
	path string
}

// Open returns a logger writing to $XDG_STATE_HOME/dotfiles/dotfiles.log.
// Logging is best effort: if the file cannot be opened the returned logger
// silently discards writes rather than failing the run.
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

// rotate keeps a single previous log around once the current one gets large.
func rotate(path string) {
	info, err := os.Stat(path)
	if err != nil || info.Size() < maxSize {
		return
	}
	_ = os.Rename(path, path+".1")
}

// Path returns the log location, or "" when logging is disabled.
func (l *Logger) Path() string {
	if l == nil {
		return ""
	}
	return l.path
}

// Printf appends one timestamped line.
func (l *Logger) Printf(format string, args ...any) {
	if l == nil || l.file == nil {
		return
	}
	_, _ = fmt.Fprintf(l.file, "[%s] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
}

// Close flushes and releases the log file.
func (l *Logger) Close() {
	if l == nil || l.file == nil {
		return
	}
	_ = l.file.Close()
}
