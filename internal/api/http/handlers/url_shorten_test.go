package handlers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/http/handlers"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestURLHandler_Shorten(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		body            io.Reader
		setupMock       func(m *mocks.MockURLService)
		wantStatus      int
		wantBody        string
		wantContentType string
	}{
		{
			name: "success",
			body: strings.NewReader("https://practicum.yandex.ru/"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Shorten(gomock.Any(), "https://practicum.yandex.ru/").
					Return("abc12345", nil)
			},
			wantStatus:      http.StatusCreated,
			wantBody:        "http://localhost:8080/abc12345",
			wantContentType: "text/plain",
		},
		{
			name: "trims surrounding whitespace",
			body: strings.NewReader("  https://example.com  \n"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Shorten(gomock.Any(), "https://example.com").
					Return("abc12345", nil)
			},
			wantStatus:      http.StatusCreated,
			wantBody:        "http://localhost:8080/abc12345",
			wantContentType: "text/plain",
		},
		{
			name:       "empty body",
			body:       strings.NewReader(""),
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
		},
		{
			name:       "whitespace only",
			body:       strings.NewReader("   \n\t"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
		},
		{
			name:       "read body error",
			body:       errReader{},
			wantStatus: http.StatusBadRequest,
			wantBody:   "bad request\n",
		},
		{
			name: "service error",
			body: strings.NewReader("https://practicum.yandex.ru/"),
			setupMock: func(m *mocks.MockURLService) {
				m.EXPECT().
					Shorten(gomock.Any(), "https://practicum.yandex.ru/").
					Return("", errors.New("shorten failed"))
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

			h := handlers.NewURLHandler(svc)
			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			rec := httptest.NewRecorder()

			h.Shorten(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, tt.wantBody, rec.Body.String())
			if tt.wantContentType != "" {
				require.Equal(t, tt.wantContentType, rec.Header().Get("Content-Type"))
			}
		})
	}
}
