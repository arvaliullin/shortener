// Package mcp вызывает инструменты сервиса ссылок по Streamable HTTP.
package mcp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Client открывает сессию MCP на время одного ответа чата.
type Client struct {
	endpoint string
}

// New создает клиента MCP.
func New(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

// Session держит соединение со stateless-сервером инструментов.
type Session struct {
	session *sdkmcp.ClientSession
}

// Open подключается к /mcp без отдельного SSE-потока: сервер ссылок stateless.
func (c *Client) Open(ctx context.Context) (*Session, error) {
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "shortener-agent", Version: "v1.0.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{
		Endpoint:             c.endpoint,
		DisableStandaloneSSE: true,
		HTTPClient:           directClient(),
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("connect mcp: %w", err)
	}
	return &Session{session: session}, nil
}

// directClient ходит по заданному адресу и не берет HTTP-прокси из окружения.
func directClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return nil, nil
			},
		},
	}
}

// Close завершает сессию MCP.
func (s *Session) Close() error {
	return s.session.Close()
}
