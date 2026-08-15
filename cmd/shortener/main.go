package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/arvaliullin/shortener/internal/api/http/router"
	"github.com/arvaliullin/shortener/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router.New(cfg),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()

	server.Shutdown(context.Background())
}
