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
		name           string
		yaml           string
		env            map[string]string
		wantURL        string
		wantToken      string
		wantTimeout    time.Duration
		wantLogLevel   string
		wantTransport  string
		wantListenAddr string
		wantAPIKey     string
	}{
		{
			name:           "yaml only",
			yaml:           "url: https://yt.example.com\ntoken: abc123\n",
			wantURL:        "https://yt.example.com",
			wantToken:      "abc123",
			wantTimeout:    30 * time.Second,
			wantLogLevel:   "info",
			wantTransport:  "stdio",
			wantListenAddr: ":8080",
		},
		{
			name:           "http transport with api key",
			yaml:           "url: https://yt.example.com\ntoken: abc123\ntransport: http\nlisten_addr: \":9000\"\napi_key: topsecret\n",
			wantURL:        "https://yt.example.com",
			wantToken:      "abc123",
			wantTimeout:    30 * time.Second,
			wantLogLevel:   "info",
			wantTransport:  "http",
			wantListenAddr: ":9000",
			wantAPIKey:     "topsecret",
		},
		{
			name: "env overrides yaml url",
			yaml: "url: https://yt.example.com\ntoken: abc123\n",
			env: map[string]string{
				"YOUTRACK_URL": "https://override.example.com",
			},
			wantURL:        "https://override.example.com",
			wantToken:      "abc123",
			wantTimeout:    30 * time.Second,
			wantLogLevel:   "info",
			wantTransport:  "stdio",
			wantListenAddr: ":8080",
		},
		{
			name: "env only without yaml file",
			env: map[string]string{
				"YOUTRACK_URL":   "https://env.example.com",
				"YOUTRACK_TOKEN": "envtoken",
			},
			wantURL:        "https://env.example.com",
			wantToken:      "envtoken",
			wantTimeout:    30 * time.Second,
			wantLogLevel:   "info",
			wantTransport:  "stdio",
			wantListenAddr: ":8080",
		},
		{
			name:           "defaults for timeout and log_level",
			yaml:           "url: https://yt.example.com\ntoken: abc123\ntimeout: 10s\nlog_level: debug\n",
			wantURL:        "https://yt.example.com",
			wantToken:      "abc123",
			wantTimeout:    10 * time.Second,
			wantLogLevel:   "debug",
			wantTransport:  "stdio",
			wantListenAddr: ":8080",
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
			assert.Equal(t, tt.wantTransport, cfg.Transport)
			assert.Equal(t, tt.wantListenAddr, cfg.ListenAddr)
			assert.Equal(t, tt.wantAPIKey, cfg.APIKey)
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
		{
			name:    "http transport without api key",
			yaml:    "url: https://yt.example.com\ntoken: abc123\ntransport: http\n",
			wantErr: "config: api_key is required for http transport",
		},
		{
			name:    "unknown transport",
			yaml:    "url: https://yt.example.com\ntoken: abc123\ntransport: grpc\n",
			wantErr: `config: unknown transport "grpc"`,
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
