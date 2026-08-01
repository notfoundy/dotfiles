// Package selfupdate replaces the running binary with the latest GitHub release.
//
// It relies on the release assets being named without a version, so the
// /releases/latest/download/ URLs are stable and no GitHub API call — nor jq —
// is needed. The same property is what lets bin/bootstrap.sh stay tiny.
package selfupdate

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	repo         = "notfoundy/dotfiles"
	latestURL    = "https://github.com/" + repo + "/releases/latest"
	downloadBase = "https://github.com/" + repo + "/releases/latest/download/"

	// maxDownload caps what we are willing to read from the network.
	maxDownload = 64 << 20
)

var client = &http.Client{Timeout: 60 * time.Second}

// AssetName is the archive published for the current platform.
func AssetName() string {
	return fmt.Sprintf("dotfiles_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
}

// LatestVersion resolves the tag the "latest" release points at.
func LatestVersion() (string, error) {
	resp, err := client.Get(latestURL)
	if err != nil {
		return "", fmt.Errorf("checking for the latest release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checking for the latest release: unexpected status %s", resp.Status)
	}

	// The redirect lands on .../releases/tag/vX.Y.Z
	tag := path.Base(resp.Request.URL.Path)
	if tag == "" || tag == "latest" {
		return "", fmt.Errorf("could not determine the latest release tag")
	}
	return tag, nil
}

// IsNewer reports whether latest differs from current, ignoring a leading "v".
// A development build is never considered up to date.
func IsNewer(current, latest string) bool {
	if current == "dev" || current == "" {
		return true
	}
	return normalize(current) != normalize(latest)
}

func normalize(v string) string { return strings.TrimPrefix(strings.TrimSpace(v), "v") }

// Apply downloads the latest archive, verifies its checksum and swaps the
// running binary for it.
func Apply() error {
	asset := AssetName()

	sum, err := expectedChecksum(asset)
	if err != nil {
		return err
	}

	archive, err := download(downloadBase + asset)
	if err != nil {
		return err
	}

	if got := sha256.Sum256(archive); hex.EncodeToString(got[:]) != sum {
		return fmt.Errorf("checksum mismatch for %s — refusing to install", asset)
	}

	binary, err := extract(archive, "dotfiles")
	if err != nil {
		return err
	}
	return replaceRunning(binary)
}

// expectedChecksum pulls the entry for asset out of the published checksums.txt.
func expectedChecksum(asset string) (string, error) {
	raw, err := download(downloadBase + "checksums.txt")
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(raw), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == asset {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("no checksum published for %s", asset)
}

func download(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("downloading %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading %s: unexpected status %s", url, resp.Status)
	}
	return io.ReadAll(io.LimitReader(resp.Body, maxDownload))
}

// extract pulls a single file out of a gzipped tarball.
func extract(archive []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("reading archive: %w", err)
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("%s not found in the release archive", name)
		}
		if err != nil {
			return nil, fmt.Errorf("reading archive: %w", err)
		}
		if header.Typeflag != tar.TypeReg || path.Base(header.Name) != name {
			continue
		}
		return io.ReadAll(io.LimitReader(tr, maxDownload))
	}
}

// replaceRunning swaps the executable atomically, by writing the new binary
// beside it and renaming over the old one.
func replaceRunning(binary []byte) error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	if resolved, err := filepath.EvalSymlinks(self); err == nil {
		self = resolved
	}

	tmp, err := os.CreateTemp(filepath.Dir(self), ".dotfiles-update-*")
	if err != nil {
		return fmt.Errorf("cannot write next to %s: %w", self, err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()

	if _, err := tmp.Write(binary); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o755); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), self)
}
