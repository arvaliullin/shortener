package server

import (
	"context"
	"net/http"

	"github.com/arvaliullin/shortener/internal/agent/chat"
	"github.com/arvaliullin/shortener/internal/agent/config"
	agentmcp "github.com/arvaliullin/shortener/internal/agent/mcp"
	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/api/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// New собирает маршруты агента.
func New(cfg *config.Config, model chat.Model, tools *agentmcp.Client, logger zerolog.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logging(logger))
	router.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	verifier := middleware.NewVerifier(cfg.KeycloakIssuer, cfg.KeycloakJWKSURL, cfg.KeycloakClientID)
	handler := NewHandler(model, mcpOpener{client: tools}, logger)
	router.With(middleware.Authenticate(verifier)).Post("/api/chat", handler.ServeHTTP)
	return router
}

// mcpOpener приводит конкретную сессию MCP к интерфейсу обработчика.
type mcpOpener struct {
	client *agentmcp.Client
}

// Open открывает сессию инструментов.
func (o mcpOpener) Open(ctx context.Context) (Session, error) {
	return o.client.Open(ctx)
}
