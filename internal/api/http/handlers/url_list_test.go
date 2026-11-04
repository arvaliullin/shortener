package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/dto"
	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/core/domain"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_List(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		List(gomock.Any()).
		Return([]domain.ShortURL{{ID: "abc12345", Original: "https://example.com"}}, nil)

	h := handlers.NewURLHandler(svc, "http://localhost:8080")
	req := httptest.NewRequest(http.MethodGet, "/api/urls", nil)
	rec := httptest.NewRecorder()

	h.List(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var got []dto.URLItem
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, []dto.URLItem{{
		ID:     "abc12345",
		URL:    "https://example.com",
		Result: "http://localhost:8080/abc12345",
	}}, got)
}
