package handlers_test

import (
	"testing"

	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestResolveURL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		Resolve(gomock.Any(), "abc12345").
		Return("https://example.com", nil)

	session := connect(t, svc)
	output := decodeOutput[struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Result string `json:"result"`
	}](t, callTool(t, session, "resolve_url", map[string]any{"id": "abc12345"}))

	require.Equal(t, "abc12345", output.ID)
	require.Equal(t, "https://example.com", output.URL)
	require.Equal(t, "http://localhost:8080/abc12345", output.Result)
}
