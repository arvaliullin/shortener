package handlers

import (
	"context"

	"github.com/arvaliullin/shortener/internal/api/mcp/dto"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// addResolveTool регистрирует поиск исходного адреса по идентификатору.
func (h *URLHandler) addResolveTool(server *sdkmcp.Server) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "resolve_url",
		Description: "Возвращает исходный адрес короткой ссылки.",
		Annotations: &sdkmcp.ToolAnnotations{ReadOnlyHint: true},
	}, h.resolve)
}

// resolve возвращает исходный адрес по идентификатору.
func (h *URLHandler) resolve(ctx context.Context, _ *sdkmcp.CallToolRequest, input dto.IDInput) (*sdkmcp.CallToolResult, dto.URLItem, error) {
	id, err := input.URLID()
	if err != nil {
		return nil, dto.URLItem{}, err
	}

	originalURL, err := h.service.Resolve(ctx, id)
	if err != nil {
		return nil, dto.URLItem{}, err
	}

	return nil, dto.NewURLItem(h.baseURL, id, originalURL), nil
}
