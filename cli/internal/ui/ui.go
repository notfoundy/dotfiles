package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const clearLine = "\r\033[K"

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

type Task struct {
	p    *Printer
	name string
	// The pending line is still the current one: settle rewrites it.
	inline bool
}

func (p *Printer) Start(format string, args ...any) *Task {
	t := &Task{p: p, name: fmt.Sprintf(format, args...), inline: true}
	if p.tty {
		_, _ = fmt.Fprintf(p.out, "%s %s  %s", clearLine, p.muted.Render("[ ]"), p.muted.Render(t.name))
	}
	return t
}

// StartInteractive ends its line so that a prompting child writes below it.
func (p *Printer) StartInteractive(format string, args ...any) *Task {
	t := &Task{p: p, name: fmt.Sprintf(format, args...)}
	_, _ = fmt.Fprintf(p.out, "%s %s  %s\n", p.linePrefix(), p.muted.Render("[ ]"), p.muted.Render(t.name))
	return t
}

func (t *Task) Done() {
	t.settle(t.p.ok.Render("[✓]"), t.p.ok.Render(t.name))
}

func (t *Task) Skip(reason string) {
	label := t.name
	if reason != "" {
		label += " — " + reason
	}
	t.settle(t.p.muted.Render("[-]"), t.p.muted.Render(label))
}

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

func (p *Printer) linePrefix() string {
	if p.tty {
		return clearLine
	}
	return ""
}

func (p *Printer) Info(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "      %s\n", p.muted.Render(fmt.Sprintf(format, args...)))
}

func (p *Printer) Warn(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "      %s\n", p.warn.Render(fmt.Sprintf(format, args...)))
}

func (p *Printer) Detail(output string) {
	if output == "" {
		return
	}
	_, _ = fmt.Fprintf(p.out, "%s\n", p.muted.Render(indent(output)))
}

func (p *Printer) Notice(format string, args ...any) {
	_, _ = fmt.Fprintf(p.out, "%s %s\n", p.accent.Render("▶"), p.notice.Render(fmt.Sprintf(format, args...)))
}

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
