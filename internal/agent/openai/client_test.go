package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvaliullin/shortener/internal/agent/chat"
	"github.com/stretchr/testify/require"
)

func TestCompleteReadsToolCall(t *testing.T) {
	t.Parallel()

	var seen wireRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/chat/completions", r.URL.Path)
		require.Equal(t, "Bearer ollama", r.Header.Get("Authorization"))
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, &seen))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":null,"tool_calls":[{"id":"call_1","type":"function","function":{"name":"list_urls","arguments":{}}}]}}]}`))
	}))
	t.Cleanup(server.Close)

	client := New(server.URL+"/v1", "qwen2.5:1.5b")
	client.http = server.Client()
	completion, err := client.Complete(context.Background(), []chat.Message{{Role: "user", Content: "покажи ссылки"}}, []chat.Tool{{
		Name:       "list_urls",
		Parameters: json.RawMessage(`{"type":"object"}`),
	}})

	require.NoError(t, err)
	require.Equal(t, "qwen2.5:1.5b", seen.Model)
	require.Equal(t, "list_urls", seen.Tools[0].Function.Name)
	require.Equal(t, []chat.ToolCall{{ID: "call_1", Name: "list_urls", Arguments: "{}"}}, completion.ToolCalls)
}

func TestCompleteReadsText(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"готово"}}]}`))
	}))
	t.Cleanup(server.Close)

	client := New(server.URL+"/v1", "qwen2.5:1.5b")
	client.http = server.Client()
	completion, err := client.Complete(context.Background(), []chat.Message{{Role: "user", Content: "привет"}}, nil)

	require.NoError(t, err)
	require.Equal(t, "готово", completion.Content)
	require.Empty(t, completion.ToolCalls)
}

func TestCompleteRejectsModelStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := New(server.URL+"/v1", "qwen2.5:1.5b")
	client.http = server.Client()
	_, err := client.Complete(context.Background(), []chat.Message{{Role: "user", Content: "привет"}}, nil)

	require.ErrorContains(t, err, "model status 404")
}
