package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	tests := []struct {
		name         string
		yaml         string
		env          map[string]string
		wantURL      string
		wantToken    string
		wantTimeout  time.Duration
		wantLogLevel string
	}{
		{
			name:         "yaml only",
			yaml:         "url: https://yt.example.com\ntoken: abc123\n",
			wantURL:      "https://yt.example.com",
			wantToken:    "abc123",
			wantTimeout:  30 * time.Second,
			wantLogLevel: "info",
		},
		{
			name: "env overrides yaml url",
			yaml: "url: https://yt.example.com\ntoken: abc123\n",
			env: map[string]string{
				"YOUTRACK_URL": "https://override.example.com",
			},
			wantURL:      "https://override.example.com",
			wantToken:    "abc123",
			wantTimeout:  30 * time.Second,
			wantLogLevel: "info",
		},
		{
			name: "env only without yaml file",
			env: map[string]string{
				"YOUTRACK_URL":   "https://env.example.com",
				"YOUTRACK_TOKEN": "envtoken",
			},
			wantURL:      "https://env.example.com",
			wantToken:    "envtoken",
			wantTimeout:  30 * time.Second,
			wantLogLevel: "info",
		},
		{
			name:         "defaults for timeout and log_level",
			yaml:         "url: https://yt.example.com\ntoken: abc123\ntimeout: 10s\nlog_level: debug\n",
			wantURL:      "https://yt.example.com",
			wantToken:    "abc123",
			wantTimeout:  10 * time.Second,
			wantLogLevel: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfgPath := filepath.Join(t.TempDir(), "config.yaml")

			if tt.yaml != "" {
				require.NoError(t, os.WriteFile(cfgPath, []byte(tt.yaml), 0o644))
			}

			cfg, err := Load(cfgPath)

			require.NoError(t, err)
			assert.Equal(t, tt.wantURL, cfg.URL)
			assert.Equal(t, tt.wantToken, cfg.Token)
			assert.Equal(t, tt.wantTimeout, cfg.Timeout)
			assert.Equal(t, tt.wantLogLevel, cfg.LogLevel)
		})
	}
}

func TestLoad_Failure(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		env     map[string]string
		wantErr string
	}{
		{
			name:    "missing url",
			yaml:    "token: abc123\n",
			wantErr: "config: url is required",
		},
		{
			name:    "missing token",
			yaml:    "url: https://yt.example.com\n",
			wantErr: "config: token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfgPath := filepath.Join(t.TempDir(), "config.yaml")

			if tt.yaml != "" {
				require.NoError(t, os.WriteFile(cfgPath, []byte(tt.yaml), 0o644))
			}

			_, err := Load(cfgPath)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
