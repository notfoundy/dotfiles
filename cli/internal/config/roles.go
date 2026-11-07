package config

import (
	"fmt"
	"strings"
)

type Selection struct {
	Roles []string
	// False for the default set: without -t, Ansible does not filter the
	// meta dependencies (hyprdot -> starship) out.
	Explicit bool
}

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
