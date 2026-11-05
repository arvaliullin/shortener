package openai

import (
	"encoding/json"

	"github.com/arvaliullin/shortener/internal/agent/chat"
)

// completionRequest собирает тело Chat Completions.
func completionRequest(model string, messages []chat.Message, tools []chat.Tool) wireRequest {
	request := wireRequest{
		Model:    model,
		Messages: make([]wireMessage, 0, len(messages)),
	}
	for _, message := range messages {
		request.Messages = append(request.Messages, wireMessageFrom(message))
	}
	for _, tool := range tools {
		request.Tools = append(request.Tools, wireToolFrom(tool))
	}
	return request
}

// wireMessageFrom переводит реплику в формат OpenAI.
func wireMessageFrom(message chat.Message) wireMessage {
	wire := wireMessage{
		Role:       message.Role,
		Content:    message.Content,
		Name:       message.Name,
		ToolCallID: message.ToolCallID,
	}
	for _, call := range message.ToolCalls {
		wire.ToolCalls = append(wire.ToolCalls, wireToolCall{
			ID:   call.ID,
			Type: "function",
			Function: wireFunctionCall{
				Name:      call.Name,
				Arguments: call.Arguments,
			},
		})
	}
	return wire
}

// wireToolFrom переводит описание инструмента в формат OpenAI.
func wireToolFrom(tool chat.Tool) wireTool {
	parameters := tool.Parameters
	if len(parameters) == 0 {
		parameters = json.RawMessage(`{"type":"object","properties":{}}`)
	}
	return wireTool{
		Type: "function",
		Function: wireFunction{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  parameters,
		},
	}
}

type wireRequest struct {
	Model    string        `json:"model"`
	Messages []wireMessage `json:"messages"`
	Tools    []wireTool    `json:"tools,omitempty"`
}

type wireMessage struct {
	Role       string         `json:"role"`
	Content    string         `json:"content,omitempty"`
	Name       string         `json:"name,omitempty"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
	ToolCalls  []wireToolCall `json:"tool_calls,omitempty"`
}

type wireToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function wireFunctionCall `json:"function"`
}

type wireFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type wireTool struct {
	Type     string       `json:"type"`
	Function wireFunction `json:"function"`
}

type wireFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters"`
}
