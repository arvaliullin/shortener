package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/core/domain"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		req        *http.Request
		setupMock  func(m *mocks.MockURLService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			req:  requestWithChiID("abc12345"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().Delete(gomock.Any(), "abc12345").Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "empty id",
			req:        httptest.NewRequest(http.MethodDelete, "/api/urls/", nil),
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
		},
		{
			name: "not found",
			req:  requestWithChiID("missing1"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().Delete(gomock.Any(), "missing1").Return(domain.ErrURLNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "not found\n",
		},
		{
			name: "service error",
			req:  requestWithChiID("abc12345"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().Delete(gomock.Any(), "abc12345").Return(errors.New("delete failed"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			svc := mocks.NewMockURLService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			h := handlers.NewURLHandler(svc, "http://localhost:8080")
			rec := httptest.NewRecorder()
			h.Delete(rec, tt.req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, tt.wantBody, rec.Body.String())
		})
	}
}
