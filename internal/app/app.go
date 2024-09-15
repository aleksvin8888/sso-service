package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {

	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	log.Info("Initialise storage successful")

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	log.Info("Initialise Auth Service successful")

	grpcApp := grpcapp.New(log, authService, grpcPort)

	log.Info("Initialise GRPC server successful")
	return &App{
		GRPCSrv: grpcApp,
	}

}
