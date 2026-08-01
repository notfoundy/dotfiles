// Package system prepares the host: distribution detection, base packages,
// SSH keys and the dotfiles checkout.
package system

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Distro is the detected distribution identity, as reported by os-release.
type Distro struct {
	ID   string // e.g. "fedora"
	Like []string
	Name string // pretty name, for messages
}

// Supported reports whether the CLI knows how to prepare this host.
//
// Only Fedora is handled, mirroring the roles: ten of them dispatch on
// {{ ansible_distribution | lower }}.yml and only ship a fedora.yml.
func (d Distro) Supported() bool {
	if d.ID == "fedora" {
		return true
	}
	for _, like := range d.Like {
		if like == "fedora" {
			return true
		}
	}
	return false
}

// DetectDistro reads /etc/os-release.
func DetectDistro() (Distro, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return Distro{}, fmt.Errorf("cannot read /etc/os-release: %w", err)
	}
	defer func() { _ = f.Close() }()

	fields := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if !found {
			continue
		}
		fields[key] = strings.Trim(value, `"'`)
	}
	if err := scanner.Err(); err != nil {
		return Distro{}, fmt.Errorf("cannot read /etc/os-release: %w", err)
	}

	d := Distro{ID: fields["ID"], Name: fields["PRETTY_NAME"]}
	if like := fields["ID_LIKE"]; like != "" {
		d.Like = strings.Fields(like)
	}
	if d.Name == "" {
		d.Name = d.ID
	}
	return d, nil
}
