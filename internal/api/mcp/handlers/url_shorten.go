package handlers

import (
	"context"

	"github.com/arvaliullin/shortener/internal/api/mcp/dto"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// addShortenTool регистрирует создание короткой ссылки.
func (h *URLHandler) addShortenTool(server *sdkmcp.Server) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "shorten_url",
		Description: "Создает короткую ссылку для исходного адреса.",
	}, h.shorten)
}

// shorten создает короткую ссылку.
func (h *URLHandler) shorten(ctx context.Context, _ *sdkmcp.CallToolRequest, input dto.ShortenInput) (*sdkmcp.CallToolResult, dto.ShortenOutput, error) {
	originalURL, err := input.OriginalURL()
	if err != nil {
		return nil, dto.ShortenOutput{}, err
	}

	id, err := h.service.Shorten(ctx, originalURL)
	if err != nil {
		return nil, dto.ShortenOutput{}, err
	}

	return nil, dto.NewShortenOutput(h.baseURL, id), nil
}
