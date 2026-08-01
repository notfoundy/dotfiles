// Package config reads the Ansible-side configuration of the dotfiles repo —
// group_vars/all.yml and the set of roles available under roles/ — and owns the
// run log every other package writes to.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// Secret is one entry of op.ssh.<provider>.<user> in group_vars/all.yml.
type Secret struct {
	Name      string `yaml:"name"`
	VaultPath string `yaml:"vault_path"`
}

// groupVars mirrors the subset of group_vars/all.yml the CLI cares about.
// Everything else in that file is consumed by Ansible only.
type groupVars struct {
	DefaultRoles []string `yaml:"default_roles"`
	ExcludeRoles []string `yaml:"exclude_roles"`
	OP           struct {
		// provider -> account -> secrets, e.g. ssh.github.notfoundy[]
		SSH map[string]map[string][]Secret `yaml:"ssh"`
	} `yaml:"op"`
}

// Config is the resolved view of a dotfiles repo checkout.
type Config struct {
	Root         string   // repo root holding main.yml
	DefaultRoles []string // default_roles minus exclude_roles
	ExcludeRoles []string
	Available    []string // roles/ entries that actually hold tasks/main.yml
	Secrets      []Secret // flattened op.ssh.*
}

// Load reads group_vars/all.yml and scans roles/ under root.
func Load(root string) (*Config, error) {
	raw, err := os.ReadFile(filepath.Join(root, "group_vars", "all.yml"))
	if err != nil {
		return nil, fmt.Errorf("reading group_vars/all.yml: %w", err)
	}

	var gv groupVars
	if err := yaml.Unmarshal(raw, &gv); err != nil {
		return nil, fmt.Errorf("parsing group_vars/all.yml: %w", err)
	}

	available, err := scanRoles(filepath.Join(root, "roles"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Root:         root,
		DefaultRoles: subtract(gv.DefaultRoles, gv.ExcludeRoles),
		ExcludeRoles: gv.ExcludeRoles,
		Available:    available,
		Secrets:      flattenSecrets(gv.OP.SSH),
	}

	if len(cfg.DefaultRoles) == 0 {
		return nil, fmt.Errorf("no default_roles found in %s", filepath.Join(root, "group_vars", "all.yml"))
	}
	return cfg, nil
}

// scanRoles lists directories under rolesDir that carry a tasks/main.yml.
// A bare directory is not a usable role and must not be offered as a tag.
func scanRoles(rolesDir string) ([]string, error) {
	entries, err := os.ReadDir(rolesDir)
	if err != nil {
		return nil, fmt.Errorf("reading roles/: %w", err)
	}

	var roles []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(rolesDir, e.Name(), "tasks", "main.yml")); err != nil {
			continue
		}
		roles = append(roles, e.Name())
	}
	sort.Strings(roles)

	if len(roles) == 0 {
		return nil, fmt.Errorf("no roles found in %s", rolesDir)
	}
	return roles, nil
}

func flattenSecrets(ssh map[string]map[string][]Secret) []Secret {
	var out []Secret
	providers := sortedKeys(ssh)
	for _, p := range providers {
		for _, account := range sortedKeys(ssh[p]) {
			out = append(out, ssh[p][account]...)
		}
	}
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func subtract(from, remove []string) []string {
	excluded := make(map[string]bool, len(remove))
	for _, r := range remove {
		excluded[r] = true
	}

	out := make([]string, 0, len(from))
	for _, item := range from {
		if !excluded[item] {
			out = append(out, item)
		}
	}
	return out
}
