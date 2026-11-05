package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Call выполняет инструмент и возвращает текст, который можно отдать модели.
func (s *Session) Call(ctx context.Context, name, arguments string) (string, error) {
	args, err := argumentsOf(arguments)
	if err != nil {
		return "", err
	}

	result, err := s.session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return "", fmt.Errorf("call tool %s: %w", name, err)
	}

	text := toolText(result)
	if result.IsError && text == "" {
		return "", errors.New("tool failed")
	}
	return text, nil
}

// argumentsOf разбирает JSON-объект аргументов.
func argumentsOf(arguments string) (any, error) {
	if strings.TrimSpace(arguments) == "" {
		return map[string]any{}, nil
	}
	var args any
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return nil, errors.New("invalid tool arguments")
	}
	return args, nil
}

// toolText собирает ответ инструмента в одну строку.
func toolText(result *sdkmcp.CallToolResult) string {
	if result.StructuredContent != nil {
		body, err := json.Marshal(result.StructuredContent)
		if err == nil {
			return briefURLList(string(body))
		}
	}

	var parts []string
	for _, content := range result.Content {
		text, ok := content.(*sdkmcp.TextContent)
		if ok && text.Text != "" {
			parts = append(parts, text.Text)
		}
	}
	return briefURLList(strings.Join(parts, "\n"))
}
