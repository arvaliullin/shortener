package handlers_test

import (
	"context"
	"encoding/json"
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
	output := decodeOutput[struct {
		URLs []struct {
			ID     string `json:"id"`
			URL    string `json:"url"`
			Result string `json:"result"`
		} `json:"urls"`
	}](t, callTool(t, session, "list_urls", map[string]any{}))

	require.Equal(t, "abc12345", output.URLs[0].ID)
	require.Equal(t, "https://example.com", output.URLs[0].URL)
	require.Equal(t, "http://localhost:8080/abc12345", output.URLs[0].Result)
}

func TestToolOutputSchemasAreObjects(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	session := connect(t, mocks.NewMockURLService(ctrl))
	tools, err := session.ListTools(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, tools.Tools)

	for _, tool := range tools.Tools {
		body, err := json.Marshal(tool.OutputSchema)
		require.NoError(t, err)

		var schema struct {
			Type string `json:"type"`
		}
		require.NoError(t, json.Unmarshal(body, &schema))
		require.Equal(t, "object", schema.Type, tool.Name)
	}
}
