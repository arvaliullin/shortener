package handlers_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/shortener/internal/api/mcp/handlers"
	"github.com/arvaliullin/shortener/internal/core/ports"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

func connect(t *testing.T, svc ports.URLService) *sdkmcp.ClientSession {
	t.Helper()

	httpServer := httptest.NewServer(handlers.NewURLHandler(svc, "http://localhost:8080"))
	t.Cleanup(httpServer.Close)

	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	session, err := client.Connect(context.Background(), &sdkmcp.StreamableClientTransport{
		Endpoint:             httpServer.URL,
		DisableStandaloneSSE: true,
	}, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, session.Close())
	})
	return session
}

func callTool(t *testing.T, session *sdkmcp.ClientSession, name string, arguments map[string]any) *sdkmcp.CallToolResult {
	t.Helper()

	result, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name:      name,
		Arguments: arguments,
	})
	require.NoError(t, err)
	return result
}

func decodeOutput[T any](t *testing.T, result *sdkmcp.CallToolResult) T {
	t.Helper()

	require.False(t, result.IsError)
	body, err := json.Marshal(result.StructuredContent)
	require.NoError(t, err)

	var output T
	require.NoError(t, json.Unmarshal(body, &output))
	return output
}
