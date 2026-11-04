package handlers_test

import (
	"testing"

	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		Shorten(gomock.Any(), "https://practicum.yandex.ru/").
		Return("abc12345", nil)

	session := connect(t, svc)
	output := decodeOutput[struct {
		ID     string `json:"id"`
		Result string `json:"result"`
	}](t, callTool(t, session, "shorten_url", map[string]any{
		"url": "https://practicum.yandex.ru/",
	}))

	require.Equal(t, "abc12345", output.ID)
	require.Equal(t, "http://localhost:8080/abc12345", output.Result)
}
