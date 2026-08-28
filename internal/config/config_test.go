package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{name: "valid", config: &Config{Port: "8080", InstagramRedirectURI: "https://app.example/api/auth/callback", FrontendURL: "https://app.example"}},
		{name: "invalid port", config: &Config{Port: "70000", InstagramRedirectURI: "https://app.example/callback", FrontendURL: "https://app.example"}, wantErr: true},
		{name: "invalid redirect", config: &Config{Port: "8080", InstagramRedirectURI: "javascript:alert(1)", FrontendURL: "https://app.example"}, wantErr: true},
		{name: "invalid frontend", config: &Config{Port: "8080", InstagramRedirectURI: "https://app.example/callback", FrontendURL: "//app.example"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.Validate()
			if (err != nil) != test.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestGetEnvReadsSecretFile(t *testing.T) {
	directory := t.TempDir()
	filename := filepath.Join(directory, "app-secret")
	if err := os.WriteFile(filename, []byte("  secret-from-file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TEST_SECRET", "")
	t.Setenv("TEST_SECRET_FILE", filename)
	if got := getEnv("TEST_SECRET", "fallback"); got != "secret-from-file" {
		t.Fatalf("getEnv() = %q, want secret from file", got)
	}
}
