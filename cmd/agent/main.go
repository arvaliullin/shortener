package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arvaliullin/shortener/internal/agent/config"
	agentmcp "github.com/arvaliullin/shortener/internal/agent/mcp"
	"github.com/arvaliullin/shortener/internal/agent/openai"
	"github.com/arvaliullin/shortener/internal/agent/server"
	"github.com/rs/zerolog"
)

// main запускает HTTP-сервер агента.
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	httpServer := &http.Server{
		Addr:              cfg.Address,
		Handler:           server.New(cfg, openai.New(cfg.OllamaBaseURL, cfg.OllamaModel), agentmcp.New(cfg.MCPURL), logger),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	httpServer.Shutdown(context.Background())
}
