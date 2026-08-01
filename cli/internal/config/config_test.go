package config

import (
	"path/filepath"
	"slices"
	"testing"
)

// repoRoot is the checkout this CLI lives in, three levels up from
// cli/internal/config.
const repoRoot = "../../.."

// TestLoadRepo reads the real group_vars/all.yml and roles/ directory, so the
// CLI's view of the repo stays in sync with the Ansible side.
func TestLoadRepo(t *testing.T) {
	cfg, err := Load(repoRoot)
	if err != nil {
		t.Fatalf("Load(%q) error = %v", repoRoot, err)
	}

	if len(cfg.DefaultRoles) == 0 {
		t.Fatal("DefaultRoles is empty")
	}

	// Every default role must exist on disk, otherwise a plain `dotfiles`
	// would fail halfway through the playbook.
	for _, role := range cfg.DefaultRoles {
		if !cfg.hasRole(role) {
			t.Errorf("default role %q has no roles/%s/tasks/main.yml", role, role)
		}
	}

	// Roles reachable only through -t must still be discovered.
	for _, role := range []string{"1password", "starship", "quickshell"} {
		if !cfg.hasRole(role) {
			t.Errorf("role %q was not discovered", role)
		}
		if cfg.IsDefault(role) {
			t.Errorf("role %q should not be part of the default set", role)
		}
	}

	if !slices.IsSorted(cfg.Available) {
		t.Error("Available is not sorted")
	}
}

// TestLoadSecrets checks the op.ssh.<provider>.<account> tree is flattened.
func TestLoadSecrets(t *testing.T) {
	cfg, err := Load(repoRoot)
	if err != nil {
		t.Fatalf("Load(%q) error = %v", repoRoot, err)
	}

	if len(cfg.Secrets) == 0 {
		t.Skip("no op secrets declared in group_vars/all.yml")
	}
	for _, s := range cfg.Secrets {
		if s.Name == "" || s.VaultPath == "" {
			t.Errorf("secret %+v has an empty name or vault_path", s)
		}
	}
}

func TestLoadMissingRepo(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("Load() error = nil, want an error for a missing checkout")
	}
}
