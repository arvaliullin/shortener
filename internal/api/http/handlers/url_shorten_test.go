package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/dto"
	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_Shorten(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		Shorten(gomock.Any(), "https://practicum.yandex.ru/").
		Return("abc12345", nil)

	h := handlers.NewURLHandler(svc, "http://localhost:8080")
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://practicum.yandex.ru/"}`))
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var resp dto.ShortenResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "http://localhost:8080/abc12345", resp.Result)
}
