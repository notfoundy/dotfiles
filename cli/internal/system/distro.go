package system

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const ostreeBooted = "/run/ostree-booted"

type Distro struct {
	ID   string
	Like []string
	Name string
}

func (d Distro) Supported() bool {
	if _, err := os.Stat(ostreeBooted); err != nil {
		return false
	}
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
