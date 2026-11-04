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
	verifier := middleware.NewVerifier(cfg.KeycloakIssuer, cfg.KeycloakJWKSURL, cfg.KeycloakClientID)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(verifier))
		r.Post("/api/shorten", urlHandler.Shorten)
		r.Get("/api/urls", urlHandler.List)
		r.Put("/api/urls/{id}", urlHandler.Update)
		r.Delete("/api/urls/{id}", urlHandler.Delete)
	})
	router.Handle("/mcp", mcphandlers.NewURLHandler(svc, cfg.BaseURL))
	router.Get("/{id}", urlHandler.Redirect)

	return router
}
