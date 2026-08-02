package selfupdate

import (
	"net/url"
	"strings"
	"testing"
)

// TestTagFromURL pins down the pages the /releases/latest redirect can land on.
func TestTagFromURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{
			name: "release tag",
			raw:  "https://github.com/notfoundy/dotfiles/releases/tag/v1.2.0",
			want: "v1.2.0",
		},
		{
			name: "tag containing a slash",
			raw:  "https://github.com/notfoundy/dotfiles/releases/tag/cli/v1.2.0",
			want: "cli/v1.2.0",
		},
		{
			// No release published: GitHub serves the listing instead.
			name:    "releases listing",
			raw:     "https://github.com/notfoundy/dotfiles/releases",
			wantErr: true,
		},
		{
			name:    "no redirect followed",
			raw:     "https://github.com/notfoundy/dotfiles/releases/latest",
			wantErr: true,
		},
		{
			name:    "empty tag",
			raw:     "https://github.com/notfoundy/dotfiles/releases/tag/",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.raw)
			if err != nil {
				t.Fatalf("parsing %q: %v", tt.raw, err)
			}

			got, err := tagFromURL(u)
			if (err != nil) != tt.wantErr {
				t.Fatalf("tagFromURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("tagFromURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestAssetNameCarriesNoVersion guards the contract with .goreleaser.yaml: the
// archive name template omits the version, which is what keeps the
// releases/latest/download URLs stable.
func TestAssetNameCarriesNoVersion(t *testing.T) {
	got := AssetName()

	if !strings.HasPrefix(got, "dotfiles_") || !strings.HasSuffix(got, ".tar.gz") {
		t.Fatalf("AssetName() = %q, want dotfiles_<os>_<arch>.tar.gz", got)
	}
	if n := strings.Count(got, "_"); n != 2 {
		t.Errorf("AssetName() = %q, want exactly the os and arch fields", got)
	}
}

func TestIsNewer(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    bool
	}{
		{"development build", "dev", "v1.0.0", true},
		{"unset version", "", "v1.0.0", true},
		{"same version", "1.0.0", "v1.0.0", false},
		{"leading v on both sides", "v1.0.0", "v1.0.0", false},
		{"different version", "1.0.0", "v1.1.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNewer(tt.current, tt.latest); got != tt.want {
				t.Errorf("IsNewer(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.want)
			}
		})
	}
}
