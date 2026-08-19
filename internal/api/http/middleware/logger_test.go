package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestLogging(t *testing.T) {
	t.Parallel()

	const body = "hello"

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	handler := middleware.Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(body))
	}))

	req := httptest.NewRequest(http.MethodPost, "/shorten?q=1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, body, rec.Body.String())

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))

	require.Equal(t, "info", entry["level"])
	require.Equal(t, req.RequestURI, entry["uri"])
	require.Equal(t, http.MethodPost, entry["method"])
	require.Equal(t, float64(http.StatusCreated), entry["status"])
	require.Equal(t, float64(len(body)), entry["size"])
	require.Contains(t, entry, "duration")
	require.Equal(t, "request completed", entry["message"])
}

func TestLogging_SkipsHealth(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	handler := middleware.Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "OK", rec.Body.String())
	require.Empty(t, buf.Bytes())
}
