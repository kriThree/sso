package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCSrv grpcapp.App
}

func New(log *slog.Logger, port int, StoragePath string, tokenTTL time.Duration) *App {

	storage, err := sqlite.New(StoragePath)

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, port, authService)

	return &App{
		GRPCSrv: *grpcApp,
	}
}
