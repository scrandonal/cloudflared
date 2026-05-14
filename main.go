// cloudflared - A tunneling daemon that connects your infrastructure to Cloudflare.
// This is a fork of cloudflare/cloudflared with additional customizations.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// BuildTime is set at build time via ldflags
	BuildTime = "unknown"
	// BuildCommit is set at build time via ldflags
	BuildCommit = "unknown"
)

func main() {
	// Configure zerolog for human-friendly console output during development
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})

	app := &cli.App{
		Name:    "cloudflared",
		Usage:   "Cloudflare Tunnel daemon",
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, BuildCommit, BuildTime),
		Authors: []*cli.Author{
			{
				Name:  "Cloudflare",
				Email: "support@cloudflare.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to configuration file",
				EnvVars: []string{"CLOUDFLARED_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "debug", // changed from "info" to "debug" for easier local development
				Usage:   "Log level (debug, info, warn, error)",
				EnvVars: []string{"CLOUDFLARED_LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-format",
				Value:   "text",
				Usage:   "Log format (text, json)",
				EnvVars: []string{"CLOUDFLARED_LOG_FORMAT"},
			},
		},
		Before: func(c *cli.Context) error {
			return configureLogging(c.String("log-level"), c.String("log-format"))
		},
		Commands: []*cli.Command{
			tunnelCommand(),
			versionCommand(),
		},
		Action: func(c *cli.Context) error {
			// Default action: show help
			return cli.ShowAppHelp(c)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("cloudflared exited with error")
	}
}

// configureLogging sets up the global logger based on CLI flags.
func configureLogging(level, format string) error {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", level, err)
	}
	zerolog.SetGlobalLevel(parsedLevel)

	if format == "json" {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		// Use a slightly more readable time format for local console output
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05", // shortened from RFC3339 — easier to read during dev
			NoColor:    false,      // keep colors enabled; helps distinguish log levels at a glance
		})
	}

	return nil
}

// tunnelCommand returns the CLI command for managing tunnels.
func tunnelCommand() *cli.Command {
	return &cli.Command{
		Name:  "tunnel",
		Usage: "Manage and run Cloudflare Tunnels",
		Subcommands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run a tunnel",
				// TODO: add --dry-run flag
				// TODO: add --reconnect-interval flag to control backoff (default upstream is 5s, feels too aggressive)
