package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nevermore/internal/transport/handler"
	"nevermore/pkg/logger"
	"time"

	"github.com/gammazero/workerpool"

	"nevermore/config"
	"nevermore/internal/service"
	"nevermore/internal/storage"
	"nevermore/pkg/auth"
	"nevermore/pkg/hash"
)

type App struct {
	server *http.Server
	wp     *workerpool.WorkerPool
}

func New() (*App, error) {
	db, err := storage.New(config.Psql(), config.Rds(), config.Photoes()) // исправлена опечатка: Photoes -> Photos
	if err != nil {
		return nil, err
	}

	hasher := hash.NewSHA1Hasher("aboba")

	jwtSecret := config.JwtSecret()

	manager, err := auth.NewManager(jwtSecret)
	if err != nil {
		return nil, err
	}

	wp := workerpool.New(100)

	srv := service.New(db, hasher, manager, wp)

	result := &App{
		server: &http.Server{
			Addr:              config.Srv(),
			Handler:           handler.New(srv, manager),
			ReadHeaderTimeout: 5 * time.Second,
		},
		wp: wp,
	}

	return result, nil
}

func (a *App) Run(ctx context.Context) error {
	log := logger.Get()

	go func() {
		// Упрощенная версия: используем простое получение из канала
		<-ctx.Done()

		fmt.Println("Shutting down the server...")

		err := a.server.Shutdown(context.Background())
		if err != nil {
			fmt.Println("Error during server shutdown:", err)
		}

		a.wp.StopWait()

		fmt.Println("Server shut down successfully")
	}()

	log.Info().Msg("Server started")

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	log.Info().Msg("Server stopped")

	return nil
}
