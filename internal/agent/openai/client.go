// Package openai вызывает локальную модель по протоколу Chat Completions.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arvaliullin/shortener/internal/agent/chat"
)

const modelTimeout = 120 * time.Second

// Client отправляет диалог в Ollama по адресу совместимому с OpenAI.
type Client struct {
	baseURL string
	model   string
	http    *http.Client
}

// New создает клиента модели.
func New(baseURL, model string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		http:    directClient(modelTimeout),
	}
}

// directClient ходит по заданному адресу и не берет HTTP-прокси из окружения.
func directClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return nil, nil
			},
		},
	}
}

// Complete запрашивает следующий шаг и разбирает текст или вызовы инструментов.
func (c *Client) Complete(ctx context.Context, messages []chat.Message, tools []chat.Tool) (chat.Completion, error) {
	body, err := json.Marshal(completionRequest(c.model, messages, tools))
	if err != nil {
		return chat.Completion{}, fmt.Errorf("encode model request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return chat.Completion{}, fmt.Errorf("create model request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer ollama")

	response, err := c.http.Do(request)
	if err != nil {
		return chat.Completion{}, fmt.Errorf("call model: %w", err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return chat.Completion{}, fmt.Errorf("read model response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return chat.Completion{}, fmt.Errorf("model status %d", response.StatusCode)
	}

	return decodeCompletion(payload)
}
