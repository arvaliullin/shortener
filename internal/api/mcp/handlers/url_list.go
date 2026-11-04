package handlers

import (
	"context"

	"github.com/arvaliullin/shortener/internal/api/mcp/dto"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// addListTool регистрирует чтение сохраненных коротких ссылок.
func (h *URLHandler) addListTool(server *sdkmcp.Server) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_urls",
		Description: "Возвращает сохраненные короткие ссылки.",
		Annotations: &sdkmcp.ToolAnnotations{ReadOnlyHint: true},
	}, h.list)
}

// list возвращает сохраненные короткие ссылки.
func (h *URLHandler) list(ctx context.Context, _ *sdkmcp.CallToolRequest, _ dto.ListInput) (*sdkmcp.CallToolResult, dto.ListOutput, error) {
	stored, err := h.service.List(ctx)
	if err != nil {
		return nil, dto.ListOutput{}, err
	}

	items := make([]dto.URLItem, 0, len(stored))
	for _, item := range stored {
		items = append(items, dto.NewURLItem(h.baseURL, item.ID, item.Original))
	}
	return nil, dto.ListOutput{URLs: items}, nil
}
