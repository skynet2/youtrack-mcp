package cli

import (
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/skynet2/youtrack-mcp/cmd/youtrack-mcp/mcp"
	"github.com/skynet2/youtrack-mcp/internal/config"
	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "youtrack-mcp",
	Short: "MCP stdio server for YouTrack",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

		cfg, err := config.Load(configPath)
		if err != nil {
			logger.Error().Err(err).Msg("config load failed")
			return err
		}
		logger = logger.Level(parseLevel(cfg.LogLevel))

		client := youtrack.New(youtrack.Config{
			BaseURL:    cfg.URL,
			Token:      cfg.Token,
			HTTPClient: &http.Client{Timeout: cfg.Timeout},
		})

		srv := mcp.NewServer(client, logger, "0.1.0")
		logger.Info().Str("url", cfg.URL).Msg("starting MCP stdio server")
		return server.ServeStdio(srv)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yaml", "path to config file")
}

func Execute() error {
	return rootCmd.Execute()
}

func parseLevel(s string) zerolog.Level {
	lvl, err := zerolog.ParseLevel(s)
	if err != nil {
		return zerolog.InfoLevel
	}
	return lvl
}
