package router

import (
	"net/http"

	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/api/http/middleware"
	mcphandlers "github.com/arvaliullin/shortener/internal/api/mcp/handlers"
	"github.com/arvaliullin/shortener/internal/config"
	"github.com/arvaliullin/shortener/internal/core/services"
	"github.com/arvaliullin/shortener/internal/repository/redis"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func New(cfg *config.Config, logger zerolog.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logging(logger))

	healthHandler := handlers.NewHealthHandler()
	router.Get("/health", healthHandler.ServeHTTP)

	repo := redis.NewURLRepository(cfg.RedisAddr)
	svc := services.NewURLService(repo)
	urlHandler := handlers.NewURLHandler(svc, cfg.BaseURL)

	router.Post("/api/shorten", urlHandler.Shorten)
	router.Get("/api/urls", urlHandler.List)
	router.Put("/api/urls/{id}", urlHandler.Update)
	router.Delete("/api/urls/{id}", urlHandler.Delete)
	router.Handle("/mcp", mcphandlers.NewURLHandler(svc, cfg.BaseURL))
	router.Get("/{id}", urlHandler.Redirect)

	return router
}
