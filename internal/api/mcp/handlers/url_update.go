package handlers

import (
	"context"

	"github.com/arvaliullin/shortener/internal/api/mcp/dto"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// addUpdateTool регистрирует смену исходного адреса.
func (h *URLHandler) addUpdateTool(server *sdkmcp.Server) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "update_url",
		Description: "Меняет исходный адрес короткой ссылки.",
		Annotations: &sdkmcp.ToolAnnotations{IdempotentHint: true},
	}, h.update)
}

// update меняет исходный адрес существующей короткой ссылки.
func (h *URLHandler) update(ctx context.Context, _ *sdkmcp.CallToolRequest, input dto.UpdateInput) (*sdkmcp.CallToolResult, dto.URLItem, error) {
	id, err := input.URLID()
	if err != nil {
		return nil, dto.URLItem{}, err
	}
	originalURL, err := input.OriginalURL()
	if err != nil {
		return nil, dto.URLItem{}, err
	}

	if err := h.service.Update(ctx, id, originalURL); err != nil {
		return nil, dto.URLItem{}, err
	}

	return nil, dto.NewURLItem(h.baseURL, id, originalURL), nil
}
