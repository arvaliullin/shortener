package mcp

import (
	"context"
	"encoding/json"

	"github.com/arvaliullin/shortener/internal/agent/chat"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// List читает инструменты, которые можно предложить модели.
func (s *Session) List(ctx context.Context) ([]chat.Tool, error) {
	var tools []chat.Tool
	var cursor string
	for {
		result, err := s.session.ListTools(ctx, &sdkmcp.ListToolsParams{Cursor: cursor})
		if err != nil {
			return nil, err
		}
		for _, tool := range result.Tools {
			tools = append(tools, chat.Tool{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  parametersOf(tool.InputSchema),
			})
		}
		if result.NextCursor == "" {
			return tools, nil
		}
		cursor = result.NextCursor
	}
}

// parametersOf копирует JSON Schema инструмента.
func parametersOf(schema any) json.RawMessage {
	if schema == nil {
		return json.RawMessage(`{"type":"object","properties":{}}`)
	}
	body, err := json.Marshal(schema)
	if err != nil || len(body) == 0 || string(body) == "null" {
		return json.RawMessage(`{"type":"object","properties":{}}`)
	}
	return body
}
