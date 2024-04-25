package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	appLog "sso/internal/lib/log"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := appLog.Setup(cfg.Env)

	application := app.New(log, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application...", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("aplication stopped")
}
