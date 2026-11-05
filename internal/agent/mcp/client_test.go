package mcp_test

import (
	"context"
	"net/http/httptest"
	"testing"

	agentmcp "github.com/arvaliullin/shortener/internal/agent/mcp"
	"github.com/arvaliullin/shortener/internal/api/mcp/handlers"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCallShorten(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	service := mocks.NewMockURLService(ctrl)
	service.EXPECT().
		Shorten(gomock.Any(), "https://practicum.yandex.ru/").
		Return("abc12345", nil)

	httpServer := httptest.NewServer(handlers.NewURLHandler(service, "http://localhost:8080"))
	t.Cleanup(httpServer.Close)

	session, err := agentmcp.New(httpServer.URL).Open(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, session.Close())
	})

	tools, err := session.List(context.Background())
	require.NoError(t, err)

	var names []string
	for _, tool := range tools {
		names = append(names, tool.Name)
	}
	require.ElementsMatch(t, []string{"shorten_url", "list_urls", "resolve_url", "update_url", "delete_url"}, names)

	text, err := session.Call(context.Background(), "shorten_url", `{"url":"https://practicum.yandex.ru/"}`)
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"abc12345","result":"http://localhost:8080/abc12345"}`, text)
}
