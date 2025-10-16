package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nevermore/internal/transport/handler"
	"nevermore/pkg/logger"

	"github.com/gammazero/workerpool"

	"nevermore/config"
	"nevermore/internal/service"
	"nevermore/internal/storage"
	"nevermore/pkg/hash"
)

type App struct {
	server *http.Server
	wp     *workerpool.WorkerPool
}

func New() (*App, error) {
	cfg, err := config.Init()

	db, err := storage.New(cfg.Psql(), cfg.Photoes())
	if err != nil {
		return nil, err
	}

	hasher := hash.NewSHA1Hasher("aboba")

	if err != nil {
		return nil, err
	}

	wp := workerpool.New(100)

	srv := service.New(db, hasher, wp)

	result := &App{
		server: &http.Server{
			Addr:    cfg.Srv(),
			Handler: handler.New(srv),
		},
		wp: wp,
	}

	return result, nil
}

func (a *App) Run(ctx context.Context) error {
	log := logger.Get()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down the server...")

			err := a.server.Shutdown(context.Background())
			if err != nil {
				fmt.Println(err)
			}

			a.wp.StopWait()

			fmt.Println("Server shutting down successfully")

			return
		}
	}()

	log.Info().Msg("Server started")

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	log.Info().Msg("Server stopped")

	return nil
}
