package chat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToolCallsFromFunctionName(t *testing.T) {
	t.Parallel()

	calls, content := toolCallsFromContent("<function_name>list_urls</function_name>", []Tool{{Name: "list_urls"}})

	require.Equal(t, []ToolCall{{ID: "call_1", Name: "list_urls", Arguments: "{}"}}, calls)
	require.Empty(t, content)
}

func TestToolCallsFromToolCallJSON(t *testing.T) {
	t.Parallel()

	raw := `<tool_call>{"name":"shorten_url","arguments":{"url":"https://example.com"}}</tool_call>`
	calls, content := toolCallsFromContent(raw, []Tool{{Name: "shorten_url"}})

	require.Len(t, calls, 1)
	require.Equal(t, "call_1", calls[0].ID)
	require.Equal(t, "shorten_url", calls[0].Name)
	require.JSONEq(t, `{"url":"https://example.com"}`, calls[0].Arguments)
	require.Empty(t, content)
}

func TestToolCallsFromContentKeepsPlainText(t *testing.T) {
	t.Parallel()

	calls, content := toolCallsFromContent("готово", []Tool{{Name: "list_urls"}})

	require.Empty(t, calls)
	require.Equal(t, "готово", content)
}

func TestToolCallsFromContentIgnoresUnknownName(t *testing.T) {
	t.Parallel()

	raw := "<function_name>delete_url</function_name>"
	calls, content := toolCallsFromContent(raw, []Tool{{Name: "list_urls"}})

	require.Empty(t, calls)
	require.Equal(t, raw, content)
}
