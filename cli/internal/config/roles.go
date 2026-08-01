package config

import (
	"fmt"
	"strings"
)

// Selection is the outcome of resolving what the user asked for against what
// the repo actually offers.
type Selection struct {
	Roles []string
	// Explicit reports whether -t must be passed to ansible-playbook.
	//
	// When the selection is exactly the default set we deliberately pass no
	// tags at all, so main.yml sees ansible_run_tags == ['all'] and skips tag
	// filtering entirely — the behaviour of a bare `dotfiles`. Passing the full
	// list via -t would instead switch Ansible into filtering mode, where role
	// dependencies pulled in through meta/main.yml (hyprdot -> starship, kitty,
	// awww, geoclue; desktop-shell -> quickshell) are subject to the filter.
	Explicit bool
}

// Resolve validates the requested tags against the roles present in the repo.
// An empty request selects the default roles.
func (c *Config) Resolve(tags []string) (Selection, error) {
	tags = normalize(tags)
	if len(tags) == 0 {
		return Selection{Roles: c.DefaultRoles, Explicit: false}, nil
	}

	for _, tag := range tags {
		if c.hasRole(tag) {
			continue
		}
		return Selection{}, unknownRoleError(tag, c.Available)
	}

	return Selection{Roles: tags, Explicit: !sameSet(tags, c.DefaultRoles)}, nil
}

// normalize splits comma-separated values, trims them and drops duplicates
// while preserving the order the user typed.
func normalize(tags []string) []string {
	var out []string
	seen := map[string]bool{}

	for _, raw := range tags {
		for _, tag := range strings.Split(raw, ",") {
			tag = strings.TrimSpace(tag)
			if tag == "" || seen[tag] {
				continue
			}
			seen[tag] = true
			out = append(out, tag)
		}
	}
	return out
}

func (c *Config) hasRole(name string) bool {
	for _, r := range c.Available {
		if r == name {
			return true
		}
	}
	return false
}

// IsDefault reports whether a role is part of the default set. Roles outside it
// stay selectable through -t but are not pre-checked in the picker.
func (c *Config) IsDefault(name string) bool {
	for _, r := range c.DefaultRoles {
		if r == name {
			return true
		}
	}
	return false
}

func unknownRoleError(tag string, available []string) error {
	if suggestion, ok := closest(tag, available); ok {
		return fmt.Errorf("unknown role %q — did you mean %q?", tag, suggestion)
	}
	return fmt.Errorf("unknown role %q (available: %s)", tag, strings.Join(available, ", "))
}

// closest returns the nearest available role name, provided it is close enough
// to plausibly be a typo.
func closest(tag string, available []string) (string, bool) {
	best, bestDistance := "", -1

	for _, candidate := range available {
		d := levenshtein(tag, candidate)
		if bestDistance == -1 || d < bestDistance {
			best, bestDistance = candidate, d
		}
	}

	// A third of the length, so "nvm" -> "nvim" matches but "foo" -> "go" does not.
	maxDistance := len(tag)/3 + 1
	if bestDistance > maxDistance {
		return "", false
	}
	return best, true
}

func levenshtein(a, b string) int {
	ar, br := []rune(a), []rune(b)

	prev := make([]int, len(br)+1)
	curr := make([]int, len(br)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(ar); i++ {
		curr[0] = i
		for j := 1; j <= len(br); j++ {
			cost := 1
			if ar[i-1] == br[j-1] {
				cost = 0
			}
			curr[j] = min(prev[j]+1, min(curr[j-1]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[len(br)]
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	in := make(map[string]bool, len(a))
	for _, v := range a {
		in[v] = true
	}
	for _, v := range b {
		if !in[v] {
			return false
		}
	}
	return true
}
