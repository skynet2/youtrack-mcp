package cli

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/skynet2/youtrack-mcp/cmd/youtrack-mcp/mcp"
	"github.com/skynet2/youtrack-mcp/internal/config"
	"github.com/skynet2/youtrack-mcp/internal/transport"
	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "youtrack-mcp",
	Short: "MCP server for YouTrack (stdio or streamable HTTP)",
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

		if cfg.Transport == config.TransportHTTP {
			return serveHTTP(cmd.Context(), srv, cfg, logger)
		}

		logger.Info().Str("url", cfg.URL).Msg("starting MCP stdio server")
		return server.ServeStdio(srv)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yaml", "path to config file")
}

func serveHTTP(parent context.Context, srv *server.MCPServer, cfg config.Config, logger zerolog.Logger) error {
	streamable := server.NewStreamableHTTPServer(srv)
	httpSrv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: transport.BearerAuth(streamable, cfg.APIKey),
	}

	ctx, stop := signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("http shutdown failed")
		}
	}()

	logger.Info().Str("url", cfg.URL).Str("addr", cfg.ListenAddr).Msg("starting MCP streamable HTTP server")
	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
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
