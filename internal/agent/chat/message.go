package chat

import (
	"context"
	"encoding/json"
)

// Message хранит одну реплику диалога, который уходит в модель.
type Message struct {
	Role       string
	Content    string
	Name       string
	ToolCallID string
	ToolCalls  []ToolCall
}

// ToolCall описывает вызов инструмента, который попросила модель.
type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}

// Tool описывает инструмент, доступный модели.
type Tool struct {
	Name        string
	Description string
	Parameters  json.RawMessage
}

// Completion содержит ответ модели за один шаг.
type Completion struct {
	Content   string
	ToolCalls []ToolCall
}

// Model запрашивает у модели следующий шаг.
type Model interface {
	Complete(ctx context.Context, messages []Message, tools []Tool) (Completion, error)
}

// Tools читает и выполняет инструменты сервиса ссылок.
type Tools interface {
	List(ctx context.Context) ([]Tool, error)
	Call(ctx context.Context, name, arguments string) (string, error)
}
