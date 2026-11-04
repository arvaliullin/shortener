package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/dto"
	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/core/domain"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestURLHandler_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		req        *http.Request
		setupMock  func(m *mocks.MockURLService)
		wantStatus int
		wantItem   *dto.URLItem
		wantBody   string
	}{
		{
			name: "success",
			req:  requestWithBody("abc12345", strings.NewReader(`{"url":"https://example.com"}`)),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Update(gomock.Any(), "abc12345", "https://example.com").
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantItem: &dto.URLItem{
				ID:     "abc12345",
				URL:    "https://example.com",
				Result: "http://localhost:8080/abc12345",
			},
		},
		{
			name:       "empty id",
			req:        httptest.NewRequest(http.MethodPut, "/api/urls/", strings.NewReader(`{"url":"https://example.com"}`)),
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
		},
		{
			name:       "invalid body",
			req:        requestWithBody("abc12345", strings.NewReader(`{"url":"  "}`)),
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
		},
		{
			name: "not found",
			req:  requestWithBody("missing1", strings.NewReader(`{"url":"https://example.com"}`)),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Update(gomock.Any(), "missing1", "https://example.com").
					Return(domain.ErrURLNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "not found\n",
		},
		{
			name: "service error",
			req:  requestWithBody("abc12345", strings.NewReader(`{"url":"https://example.com"}`)),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Update(gomock.Any(), "abc12345", "https://example.com").
					Return(errors.New("update failed"))
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
			h.Update(rec, tt.req)

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantItem != nil {
				var got dto.URLItem
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
				require.Equal(t, *tt.wantItem, got)
				return
			}
			require.Equal(t, tt.wantBody, rec.Body.String())
		})
	}
}

func requestWithBody(id string, body io.Reader) *http.Request {
	req := requestWithChiID(id)
	req.Body = io.NopCloser(body)
	return req
}
