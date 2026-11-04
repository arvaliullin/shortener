package handlers_test

import (
	"context"
	"testing"

	"github.com/arvaliullin/shortener/internal/core/domain"
	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReadURLListResource(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		List(gomock.Any()).
		Return([]domain.ShortURL{{ID: "abc12345", Original: "https://example.com"}}, nil)

	session := connect(t, svc)
	result, err := session.ReadResource(context.Background(), &sdkmcp.ReadResourceParams{
		URI: "shortener://urls",
	})

	require.NoError(t, err)
	require.Len(t, result.Contents, 1)
	require.JSONEq(t, `[{
		"id": "abc12345",
		"url": "https://example.com",
		"result": "http://localhost:8080/abc12345"
	}]`, result.Contents[0].Text)
}
