package config

import (
	"slices"
	"strings"
	"testing"
)

func testConfig() *Config {
	return &Config{
		Available:    []string{"bash", "blender", "go", "nvim", "quickshell", "system"},
		DefaultRoles: []string{"system", "bash", "go", "nvim"},
	}
}

func TestResolveWithoutTagsSelectsDefaults(t *testing.T) {
	got, err := testConfig().Resolve(nil)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	want := []string{"system", "bash", "go", "nvim"}
	if !slices.Equal(got.Roles, want) {
		t.Errorf("roles = %v, want %v", got.Roles, want)
	}
	// No -t must be passed, so main.yml sees ansible_run_tags == ['all'] and
	// leaves tag filtering off, exactly like a bare `dotfiles`.
	if got.Explicit {
		t.Error("Explicit = true, want false for the default selection")
	}
}

func TestResolveFullDefaultSetStaysImplicit(t *testing.T) {
	// Same set, different order: still the default run.
	got, err := testConfig().Resolve([]string{"nvim,go", "bash", "system"})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got.Explicit {
		t.Error("Explicit = true, want false when the selection matches the defaults")
	}
}

func TestResolveSubsetIsExplicit(t *testing.T) {
	got, err := testConfig().Resolve([]string{"nvim", "go"})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	want := []string{"nvim", "go"}
	if !slices.Equal(got.Roles, want) {
		t.Errorf("roles = %v, want %v", got.Roles, want)
	}
	if !got.Explicit {
		t.Error("Explicit = false, want true for a subset")
	}
}

func TestResolveAcceptsRolesOutsideDefaults(t *testing.T) {
	// blender and quickshell are not in default_roles but must stay reachable.
	if _, err := testConfig().Resolve([]string{"blender", "quickshell"}); err != nil {
		t.Errorf("Resolve() error = %v, want nil", err)
	}
}

func TestResolveNormalizesInput(t *testing.T) {
	got, err := testConfig().Resolve([]string{" nvim , go ", "nvim", ""})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	want := []string{"nvim", "go"}
	if !slices.Equal(got.Roles, want) {
		t.Errorf("roles = %v, want %v (trimmed, split and deduplicated)", got.Roles, want)
	}
}

func TestResolveRejectsUnknownRoleWithSuggestion(t *testing.T) {
	_, err := testConfig().Resolve([]string{"nvm"})
	if err == nil {
		t.Fatal("Resolve() error = nil, want an error for an unknown role")
	}
	if !strings.Contains(err.Error(), `"nvim"`) {
		t.Errorf("error = %q, want it to suggest nvim", err)
	}
}

func TestResolveRejectsUnknownRoleWithoutSuggestion(t *testing.T) {
	_, err := testConfig().Resolve([]string{"kubernetes"})
	if err == nil {
		t.Fatal("Resolve() error = nil, want an error for an unknown role")
	}
	if !strings.Contains(err.Error(), "available:") {
		t.Errorf("error = %q, want it to list the available roles", err)
	}
}
