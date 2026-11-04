package handlers_test

import (
	"context"
	"testing"

	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReadURLResource(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		Resolve(gomock.Any(), "abc12345").
		Return("https://example.com", nil)

	session := connect(t, svc)
	result, err := session.ReadResource(context.Background(), &sdkmcp.ReadResourceParams{
		URI: "shortener://urls/abc12345",
	})

	require.NoError(t, err)
	require.Len(t, result.Contents, 1)
	require.JSONEq(t, `{
		"id": "abc12345",
		"url": "https://example.com",
		"result": "http://localhost:8080/abc12345"
	}`, result.Contents[0].Text)
}
