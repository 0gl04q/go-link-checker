package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/0gl04q/go-link-checker/internal/cli"
	"github.com/0gl04q/go-link-checker/internal/config"
	"github.com/lmittmann/tint"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: getLogLevel(cfg.CLI.LogLevel),
	}))
	slog.SetDefault(logger)

	if err := cli.New(cfg).Run(); err != nil {
		log.Fatalf("cli: %v", err)
	}
}

func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
