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
				Value:   "info",
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
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
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
				Action: func(c *cli.Context) error {
					log.Info().Msg("Starting tunnel...")
					// TODO: implement tunnel run logic
					return nil
				},
			},
		},
	}
}

// versionCommand returns the CLI command for printing version info.
func versionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print version information",
		Action: func(c *cli.Context) error {
			fmt.Printf("cloudflared version %s\n", Version)
			fmt.Printf("Commit:     %s\n", BuildCommit)
			fmt.Printf("Build time: %s\n", BuildTime)
			return nil
		},
	}
}
