package handlers_test

import (
	"encoding/json"
	"io"
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

func TestURLHandler_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		Update(gomock.Any(), "abc12345", "https://example.com").
		Return(nil)

	h := handlers.NewURLHandler(svc, "http://localhost:8080")
	rec := httptest.NewRecorder()

	h.Update(rec, requestWithBody("abc12345", strings.NewReader(`{"url":"https://example.com"}`)))

	require.Equal(t, http.StatusOK, rec.Code)
	var got dto.URLItem
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, dto.URLItem{
		ID:     "abc12345",
		URL:    "https://example.com",
		Result: "http://localhost:8080/abc12345",
	}, got)
}

func requestWithBody(id string, body io.Reader) *http.Request {
	req := requestWithChiID(id)
	req.Body = io.NopCloser(body)
	return req
}
