package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func requestWithChiID(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestURLHandler_Redirect(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		Resolve(gomock.Any(), "abc12345").
		Return("https://practicum.yandex.ru/", nil)

	h := handlers.NewURLHandler(svc, "http://localhost:8080")
	rec := httptest.NewRecorder()

	h.Redirect(rec, requestWithChiID("abc12345"))

	require.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	require.Equal(t, "https://practicum.yandex.ru/", rec.Header().Get("Location"))
}
