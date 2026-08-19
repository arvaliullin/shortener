package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arvaliullin/shortener/internal/api/http/router"
	"github.com/arvaliullin/shortener/internal/config"
	"github.com/rs/zerolog"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router.New(cfg, logger),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()

	server.Shutdown(context.Background())
}
