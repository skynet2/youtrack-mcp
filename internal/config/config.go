package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	URL        string        `mapstructure:"url"`
	Token      string        `mapstructure:"token"`
	Timeout    time.Duration `mapstructure:"timeout"`
	LogLevel   string        `mapstructure:"log_level"`
	Transport  string        `mapstructure:"transport"`
	ListenAddr string        `mapstructure:"listen_addr"`
	APIKey     string        `mapstructure:"api_key"`
}

const (
	TransportStdio = "stdio"
	TransportHTTP  = "http"
)

func Load(path string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("YOUTRACK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("timeout", "30s")
	v.SetDefault("log_level", "info")
	v.SetDefault("transport", TransportStdio)
	v.SetDefault("listen_addr", ":8080")

	for _, k := range []string{"url", "token", "timeout", "log_level", "transport", "listen_addr", "api_key"} {
		_ = v.BindEnv(k)
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config: read: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("config: unmarshal: %w", err)
	}

	if cfg.URL == "" {
		return Config{}, errors.New("config: url is required")
	}
	if cfg.Token == "" {
		return Config{}, errors.New("config: token is required")
	}
	switch cfg.Transport {
	case TransportStdio, TransportHTTP:
	default:
		return Config{}, fmt.Errorf("config: unknown transport %q", cfg.Transport)
	}
	return cfg, nil
}
