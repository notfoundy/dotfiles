// Package ui renders the preparation phase.
//
// Output stays on a single line per step while it runs and is rewritten in
// place once it settles. Commands executed during a step have their output
// captured to the run log, so nothing interleaves; only a failure prints it.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	reset  = "\033[0m"
	red    = "\033[1;31m"
	green  = "\033[1;32m"
	yellow = "\033[0;33m"
	cyan   = "\033[0;36m"
	gray   = "\033[0;37m"
	sea    = "\033[38;5;49m"

	clearLine = "\r\033[K"
)

// Printer writes progress to a stream.
type Printer struct {
	out   io.Writer
	tty   bool
	color bool
}

// New returns a Printer for stdout, adapting to whether it is a terminal.
func New() *Printer {
	tty := term.IsTerminal(int(os.Stdout.Fd()))
	return &Printer{
		out:   os.Stdout,
		tty:   tty,
		color: tty && os.Getenv("NO_COLOR") == "",
	}
}

func (p *Printer) paint(code, s string) string {
	if !p.color {
		return s
	}
	return code + s + reset
}

// Task is a single step of the preparation phase.
type Task struct {
	p    *Printer
	name string
}

// Start announces a step. The line is rewritten by Done or Fail.
func (p *Printer) Start(format string, args ...any) *Task {
	t := &Task{p: p, name: fmt.Sprintf(format, args...)}
	if p.tty {
		_, _ = fmt.Fprintf(p.out, "%s %s  %s", clearLine, p.paint(gray, "[ ]"), p.paint(gray, t.name))
	}
	return t
}

// Done marks the step as successful.
func (t *Task) Done() {
	t.settle(t.p.paint(green, "[✓]"), t.p.paint(green, t.name))
}

// Skip marks the step as intentionally not performed.
func (t *Task) Skip(reason string) {
	label := t.name
	if reason != "" {
		label += " — " + reason
	}
	t.settle(t.p.paint(gray, "[-]"), t.p.paint(gray, label))
}

// Fail marks the step as failed. The caller is responsible for surfacing err.
func (t *Task) Fail() {
	t.settle(t.p.paint(red, "[✗]"), t.p.paint(red, t.name))
}

func (t *Task) settle(symbol, label string) {
	prefix := ""
	if t.p.tty {
		prefix = clearLine
	}
	_, _ = fmt.Fprintf(t.p.out, "%s %s  %s\n", prefix, symbol, label)
}

// Info prints a secondary, indented line.
func (p *Printer) Info(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "      %s\n", p.paint(gray, fmt.Sprintf(format, args...)))
}

// Warn prints an indented warning.
func (p *Printer) Warn(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "      %s\n", p.paint(yellow, fmt.Sprintf(format, args...)))
}

// Detail prints captured command output after a failure.
func (p *Printer) Detail(output string) {
	if output == "" {
		return
	}
	_, _ = fmt.Fprintf(p.out, "%s\n", p.paint(gray, indent(output)))
}

// Notice prints a standalone highlighted line.
func (p *Printer) Notice(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "%s %s\n", p.paint(sea, "▶"), p.paint(cyan, fmt.Sprintf(format, args...)))
}

// Success prints a standalone success line.
func (p *Printer) Success(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "%s %s\n", p.paint(green, "✔"), p.paint(green, fmt.Sprintf(format, args...)))
}

func indent(s string) string {
	var out strings.Builder
	out.WriteString("      ")
	for _, r := range s {
		out.WriteString(string(r))
		if r == '\n' {
			out.WriteString("      ")
		}
	}
	return out.String()
}
