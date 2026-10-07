package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs-karasal/dreams-come-true/internal/config"
)

func TestEnvOrDefault(t *testing.T) {
	tests := map[string]struct {
		key, envValue, defValue, want string
	}{
		"set": {
			key:      "test-key",
			envValue: "from-env",
			defValue: "def-value",
			want:     "from-env",
		},
		"empty": {
			key:      "test-key",
			envValue: "",
			defValue: "def-value",
			want:     "def-value",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv(tc.key, tc.envValue)

			got := config.EnvOrDefault(tc.key, tc.defValue)
			if tc.want != got {
				t.Errorf(
					"getSecretOrDefault(%q, %q) = %q, want %q",
					tc.key,
					tc.defValue,
					got,
					tc.want,
				)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	tests := map[string]struct {
		dotenv  string
		envAddr string
		want    config.Config
		wantErr bool
	}{
		"defaults when nothing is set": {
			want: config.Config{Addr: ":8080"},
		},
		"reads .env": {
			dotenv: "HTTP_ADDR=:7070\n",
			want:   config.Config{Addr: ":7070"},
		},
		"real env beats .env": {
			dotenv:  "HTTP_ADDR=:7070",
			envAddr: ":9090",
			want:    config.Config{Addr: ":9090"},
		},
		"broken .env is an error": {
			dotenv:  "HTTP_ADDR=\":7070\n",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.dotenv != "" {
				path := filepath.Join(dir, ".env")
				if err := os.WriteFile(path, []byte(tc.dotenv), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			t.Chdir(dir)

			t.Setenv("HTTP_ADDR", tc.envAddr)
			if tc.envAddr == "" {
				if err := os.Unsetenv("HTTP_ADDR"); err != nil {
					t.Fatal(err)
				}
			}

			got, err := config.Load()
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetConfig() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.want != got {
				t.Errorf("GetConfig() = %+v, want %+v", got, tc.want)
			}
		})
	}
}
