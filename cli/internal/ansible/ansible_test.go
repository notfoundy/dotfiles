package ansible

import (
	"slices"
	"testing"

	"github.com/notfoundy/dotfiles/cli/internal/config"
)

func TestPlaybookArgsDefaultRunPassesNoTags(t *testing.T) {
	got := playbookArgs(PlaybookOptions{
		Root:      "/home/user/.dotfiles",
		Selection: config.Selection{Roles: []string{"system", "bash"}, Explicit: false},
	})

	// This is the parity case with the bash script's bare `dotfiles`: no -t at
	// all, so main.yml resolves the roles itself from default_roles and Ansible
	// applies no tag filtering.
	want := []string{"/home/user/.dotfiles/main.yml", "-K"}
	if !slices.Equal(got, want) {
		t.Errorf("args = %v, want %v", got, want)
	}
}

func TestPlaybookArgsExplicitSelectionPassesTags(t *testing.T) {
	got := playbookArgs(PlaybookOptions{
		Root:      "/home/user/.dotfiles",
		Selection: config.Selection{Roles: []string{"nvim", "go"}, Explicit: true},
	})

	want := []string{"/home/user/.dotfiles/main.yml", "-t", "nvim,go", "-K"}
	if !slices.Equal(got, want) {
		t.Errorf("args = %v, want %v", got, want)
	}
}

func TestPlaybookArgsForwardsExtra(t *testing.T) {
	got := playbookArgs(PlaybookOptions{
		Root:  "/repo",
		Extra: []string{"--check", "-vvv"},
	})

	want := []string{"/repo/main.yml", "--check", "-vvv", "-K"}
	if !slices.Equal(got, want) {
		t.Errorf("args = %v, want %v", got, want)
	}
}

func TestPlaybookArgsDoesNotDuplicateBecomeFlag(t *testing.T) {
	for _, flag := range []string{"-K", "--ask-become-pass", "--become-password-file=/tmp/pw"} {
		got := playbookArgs(PlaybookOptions{Root: "/repo", Extra: []string{flag}})

		want := []string{"/repo/main.yml", flag}
		if !slices.Equal(got, want) {
			t.Errorf("args with %s = %v, want %v", flag, got, want)
		}
	}
}
