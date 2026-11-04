package handlers_test

import (
	"testing"

	"github.com/arvaliullin/shortener/internal/core/domain"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListURLs(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		List(gomock.Any()).
		Return([]domain.ShortURL{{ID: "abc12345", Original: "https://example.com"}}, nil)

	session := connect(t, svc)
	output := decodeOutput[[]struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Result string `json:"result"`
	}](t, callTool(t, session, "list_urls", map[string]any{}))

	require.Equal(t, "abc12345", output[0].ID)
	require.Equal(t, "https://example.com", output[0].URL)
	require.Equal(t, "http://localhost:8080/abc12345", output[0].Result)
}
