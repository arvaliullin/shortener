package openai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/arvaliullin/shortener/internal/agent/chat"
)

// decodeCompletion читает текст ответа или вызовы инструментов.
func decodeCompletion(payload []byte) (chat.Completion, error) {
	var body wireResponse
	if err := json.Unmarshal(payload, &body); err != nil {
		return chat.Completion{}, fmt.Errorf("decode model response: %w", err)
	}
	if len(body.Choices) == 0 {
		return chat.Completion{}, fmt.Errorf("model returned no choices")
	}

	message := body.Choices[0].Message
	completion := chat.Completion{Content: contentString(message.Content)}
	for i, call := range message.ToolCalls {
		id := call.ID
		if id == "" {
			id = fmt.Sprintf("call_%d", i+1)
		}
		completion.ToolCalls = append(completion.ToolCalls, chat.ToolCall{
			ID:        id,
			Name:      call.Function.Name,
			Arguments: argumentsString(call.Function.Arguments),
		})
	}
	return completion, nil
}

// contentString достает текст, если модель вернула строку.
func contentString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	return strings.TrimSpace(string(raw))
}

// argumentsString оставляет аргументы одной JSON-строкой.
func argumentsString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return "{}"
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	return string(raw)
}

type wireResponse struct {
	Choices []struct {
		Message struct {
			Content   json.RawMessage `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Function struct {
					Name      string          `json:"name"`
					Arguments json.RawMessage `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}
