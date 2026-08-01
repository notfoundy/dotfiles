package onepass

import (
	"slices"
	"testing"

	"github.com/notfoundy/dotfiles/cli/internal/config"
)

// TestParseSession covers the shapes `op signin` prints on stdout, which the
// bash script used to hand straight to eval.
func TestParseSession(t *testing.T) {
	tests := []struct {
		name string
		out  string
		want string
	}{
		{
			name: "export statement",
			out:  `export OP_SESSION_my="abc123token"` + "\n",
			want: "OP_SESSION_my=abc123token",
		},
		{
			name: "trailing semicolon",
			out:  `export OP_SESSION_example='tok';`,
			want: "OP_SESSION_example=tok",
		},
		{
			name: "without the export keyword",
			out:  `OP_SESSION_abc=plain`,
			want: "OP_SESSION_abc=plain",
		},
		{
			name: "surrounded by noise",
			out:  "Enter the password:\nexport OP_SESSION_x=\"tok\"\n",
			want: "OP_SESSION_x=tok",
		},
		{
			// Desktop app integration authenticates without emitting a session.
			name: "no session emitted",
			out:  "",
			want: "",
		},
		{
			name: "unrelated output only",
			out:  "already signed in\n",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSession(tt.out); got != tt.want {
				t.Errorf("parseSession() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestReferences pins the field references down: `op read` rejects an
// item-level op://vault/item path, so the fields must be appended.
func TestReferences(t *testing.T) {
	tests := []struct {
		name   string
		secret config.Secret
		want   []string
	}{
		{
			name:   "one reference per field",
			secret: config.Secret{Name: "github_key", VaultPath: "op://Private/github_key"},
			want: []string{
				"op://Private/github_key/public_key",
			},
		},
		{
			name:   "no vault path",
			secret: config.Secret{Name: "github_key"},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := references(tt.secret); !slices.Equal(got, tt.want) {
				t.Errorf("references() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAccountRefPrefersShorthand(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		want    string
	}{
		{"shorthand", Account{Shorthand: "my", URL: "my.1password.eu"}, "my"},
		{"url fallback", Account{URL: "my.1password.eu"}, "my.1password.eu"},
		{"uuid fallback", Account{AccountUUID: "ABC"}, "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.account.Ref(); got != tt.want {
				t.Errorf("Ref() = %q, want %q", got, tt.want)
			}
		})
	}
}
