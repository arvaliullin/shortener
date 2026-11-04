package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().Delete(gomock.Any(), "abc12345").Return(nil)

	h := handlers.NewURLHandler(svc, "http://localhost:8080")
	rec := httptest.NewRecorder()

	h.Delete(rec, requestWithChiID("abc12345"))

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Empty(t, rec.Body.String())
}
