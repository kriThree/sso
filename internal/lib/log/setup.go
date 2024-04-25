package appLog

import (
	"log/slog"
	"os"
)

const (
	env_local = "local"
	env_prod  = "prod"
	env_dev   = "dev"
)

func Setup(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case env_local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env_dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env_prod:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}
