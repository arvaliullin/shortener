package handlers

import (
	"context"

	"github.com/arvaliullin/shortener/internal/api/mcp/dto"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// addDeleteTool регистрирует удаление короткой ссылки.
func (h *URLHandler) addDeleteTool(server *sdkmcp.Server) {
	destructive := true
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "delete_url",
		Description: "Удаляет короткую ссылку.",
		Annotations: &sdkmcp.ToolAnnotations{DestructiveHint: &destructive},
	}, h.delete)
}

// delete удаляет короткую ссылку.
func (h *URLHandler) delete(ctx context.Context, _ *sdkmcp.CallToolRequest, input dto.IDInput) (*sdkmcp.CallToolResult, dto.DeleteOutput, error) {
	id, err := input.URLID()
	if err != nil {
		return nil, dto.DeleteOutput{}, err
	}

	if err := h.service.Delete(ctx, id); err != nil {
		return nil, dto.DeleteOutput{}, err
	}

	return nil, dto.DeleteOutput{ID: id}, nil
}
