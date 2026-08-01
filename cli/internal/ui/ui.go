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

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// clearLine returns to the start of the current line and erases it.
const clearLine = "\r\033[K"

// Printer writes progress to a stream.
type Printer struct {
	out io.Writer
	tty bool

	ok     lipgloss.Style
	fail   lipgloss.Style
	warn   lipgloss.Style
	notice lipgloss.Style
	muted  lipgloss.Style
	accent lipgloss.Style
}

// New returns a Printer for stdout, adapting to whether it is a terminal.
func New() *Printer {
	out := os.Stdout
	r := lipgloss.NewRenderer(out)

	return &Printer{
		out:    out,
		tty:    term.IsTerminal(int(out.Fd())),
		ok:     r.NewStyle().Bold(true).Foreground(lipgloss.Color("2")),
		fail:   r.NewStyle().Bold(true).Foreground(lipgloss.Color("1")),
		warn:   r.NewStyle().Foreground(lipgloss.Color("3")),
		notice: r.NewStyle().Foreground(lipgloss.Color("6")),
		muted:  r.NewStyle().Foreground(lipgloss.Color("7")),
		accent: r.NewStyle().Foreground(lipgloss.Color("49")),
	}
}

// Task is a single step of the preparation phase.
type Task struct {
	p    *Printer
	name string
	// inline records that the pending line is still the current one, so Done
	// and Fail rewrite it instead of printing below it.
	inline bool
}

// Start announces a step. The line is rewritten by Done or Fail.
func (p *Printer) Start(format string, args ...any) *Task {
	t := &Task{p: p, name: fmt.Sprintf(format, args...), inline: true}
	if p.tty {
		_, _ = fmt.Fprintf(p.out, "%s %s  %s", clearLine, p.muted.Render("[ ]"), p.muted.Render(t.name))
	}
	return t
}

// StartInteractive announces a step that hands the terminal over to a child
// process, such as the 1Password password prompt. The pending line is closed so
// the child prompts below it rather than on top of it, and Done or Fail then
// print their own line, since the child's output sits in between.
func (p *Printer) StartInteractive(format string, args ...any) *Task {
	t := &Task{p: p, name: fmt.Sprintf(format, args...)}
	_, _ = fmt.Fprintf(p.out, "%s %s  %s\n", p.linePrefix(), p.muted.Render("[ ]"), p.muted.Render(t.name))
	return t
}

// Done marks the step as successful.
func (t *Task) Done() {
	t.settle(t.p.ok.Render("[✓]"), t.p.ok.Render(t.name))
}

// Skip marks the step as intentionally not performed.
func (t *Task) Skip(reason string) {
	label := t.name
	if reason != "" {
		label += " — " + reason
	}
	t.settle(t.p.muted.Render("[-]"), t.p.muted.Render(label))
}

// Fail marks the step as failed. The caller is responsible for surfacing err.
func (t *Task) Fail() {
	t.settle(t.p.fail.Render("[✗]"), t.p.fail.Render(t.name))
}

func (t *Task) settle(symbol, label string) {
	prefix := ""
	if t.inline {
		prefix = t.p.linePrefix()
	}
	_, _ = fmt.Fprintf(t.p.out, "%s %s  %s\n", prefix, symbol, label)
}

// linePrefix erases the current line on a terminal, where something pending may
// already sit on it, and is a no-op when the output is redirected.
func (p *Printer) linePrefix() string {
	if p.tty {
		return clearLine
	}
	return ""
}

// Info prints a secondary, indented line.
func (p *Printer) Info(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "      %s\n", p.muted.Render(fmt.Sprintf(format, args...)))
}

// Warn prints an indented warning.
func (p *Printer) Warn(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "      %s\n", p.warn.Render(fmt.Sprintf(format, args...)))
}

// Detail prints captured command output after a failure.
func (p *Printer) Detail(output string) {
	if output == "" {
		return
	}
	_, _ = fmt.Fprintf(p.out, "%s\n", p.muted.Render(indent(output)))
}

// Notice prints a standalone highlighted line.
func (p *Printer) Notice(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "%s %s\n", p.accent.Render("▶"), p.notice.Render(fmt.Sprintf(format, args...)))
}

// Success prints a standalone success line.
func (p *Printer) Success(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "%s %s\n", p.ok.Render("✔"), p.ok.Render(fmt.Sprintf(format, args...)))
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
