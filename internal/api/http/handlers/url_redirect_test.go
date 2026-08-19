package handlers_test

import (
	"context"
	"errors"
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

	tests := []struct {
		name         string
		req          *http.Request
		setupMock    func(m *mocks.MockURLService)
		wantStatus   int
		wantLocation string
		wantBody     string
	}{
		{
			name: "success",
			req:  requestWithChiID("abc12345"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Resolve(gomock.Any(), "abc12345").
					Return("https://practicum.yandex.ru/", nil)
			},
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://practicum.yandex.ru/",
		},
		{
			name:       "empty id",
			req:        httptest.NewRequest(http.MethodGet, "/", nil),
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
		},
		{
			name: "resolve error",
			req:  requestWithChiID("missing"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Resolve(gomock.Any(), "missing").
					Return("", errors.New("not found"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
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

			h.Redirect(rec, tt.req)

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantLocation != "" {
				require.Equal(t, tt.wantLocation, rec.Header().Get("Location"))
			}
			if tt.wantBody != "" {
				require.Equal(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
