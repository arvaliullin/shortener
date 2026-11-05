// Package chat проводит диалог с моделью и инструментами коротких ссылок.
package chat

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	maxSteps = 3

	systemPrompt = "Ты помощник сервиса коротких ссылок и отвечаешь по-русски. " +
		"Для ссылок вызывай инструменты и не выдумывай адреса. " +
		"delete_url и update_url вызывай только когда пользователь прямо просит удалить или изменить ссылку."
)

var (
	errEmptyResponse = errors.New("empty model response")
	errTooManySteps  = errors.New("too many tool steps")
)

// Loop повторяет запрос к модели, пока она не вернет текст.
type Loop struct {
	model Model
	tools Tools
}

// NewLoop собирает цикл из клиента модели и инструментов.
func NewLoop(model Model, tools Tools) *Loop {
	return &Loop{model: model, tools: tools}
}

// Reply возвращает текст ответа на уже накопленную историю.
func (l *Loop) Reply(ctx context.Context, history []Message) (string, error) {
	listed, err := l.tools.List(ctx)
	if err != nil {
		return "", fmt.Errorf("list tools: %w", err)
	}

	offered := visibleTools(history, listed)
	messages := make([]Message, 0, len(history)+1)
	messages = append(messages, Message{Role: "system", Content: systemPrompt})
	messages = append(messages, history...)

	for range maxSteps {
		completion, err := l.model.Complete(ctx, messages, offered)
		if err != nil {
			return "", fmt.Errorf("complete: %w", err)
		}
		if len(completion.ToolCalls) == 0 {
			if strings.TrimSpace(completion.Content) == "" {
				return "", errEmptyResponse
			}
			return completion.Content, nil
		}

		messages = append(messages, Message{
			Role:      "assistant",
			Content:   completion.Content,
			ToolCalls: completion.ToolCalls,
		})
		for _, call := range completion.ToolCalls {
			messages = append(messages, Message{
				Role:       "tool",
				Name:       call.Name,
				ToolCallID: call.ID,
				Content:    l.call(ctx, history, offered, call),
			})
		}
	}

	return "", errTooManySteps
}

// call выполняет один инструмент или возвращает модели причину отказа.
func (l *Loop) call(ctx context.Context, history []Message, offered []Tool, call ToolCall) string {
	if !mutationAllowed(history, call.Name) {
		return mutationRefused
	}
	if !knownTool(offered, call.Name) {
		return "Unknown tool."
	}

	result, err := l.tools.Call(ctx, call.Name, call.Arguments)
	if err != nil {
		return err.Error()
	}
	return result
}
