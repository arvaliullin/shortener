package handlers_test

import (
	"encoding/json"
	"errors"
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

	tests := []struct {
		name       string
		setupMock  func(m *mocks.MockURLService)
		wantStatus int
		wantItems  []dto.URLItem
		wantBody   string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					List(gomock.Any()).
					Return([]domain.ShortURL{{ID: "abc12345", Original: "https://example.com"}}, nil)
			},
			wantStatus: http.StatusOK,
			wantItems: []dto.URLItem{{
				ID:     "abc12345",
				URL:    "https://example.com",
				Result: "http://localhost:8080/abc12345",
			}},
		},
		{
			name: "empty",
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().List(gomock.Any()).Return([]domain.ShortURL{}, nil)
			},
			wantStatus: http.StatusOK,
			wantItems:  []dto.URLItem{},
		},
		{
			name: "service error",
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().List(gomock.Any()).Return(nil, errors.New("list failed"))
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
			tt.setupMock(svc)

			h := handlers.NewURLHandler(svc, "http://localhost:8080")
			req := httptest.NewRequest(http.MethodGet, "/api/urls", nil)
			rec := httptest.NewRecorder()

			h.List(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				require.Equal(t, tt.wantBody, rec.Body.String())
				return
			}

			var got []dto.URLItem
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
			require.Equal(t, tt.wantItems, got)
		})
	}
}
