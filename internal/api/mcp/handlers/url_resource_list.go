package handlers

import (
	"context"
	"encoding/json"

	"github.com/arvaliullin/shortener/internal/api/mcp/dto"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const urlsResourceURI = "shortener://urls"

// addURLListResource регистрирует список сохраненных коротких ссылок.
func (h *URLHandler) addURLListResource(server *sdkmcp.Server) {
	server.AddResource(&sdkmcp.Resource{
		URI:         urlsResourceURI,
		Name:        "urls",
		Description: "Список сохраненных коротких ссылок.",
		MIMEType:    "application/json",
	}, h.readURLList)
}

// readURLList возвращает JSON списка коротких ссылок.
func (h *URLHandler) readURLList(ctx context.Context, req *sdkmcp.ReadResourceRequest) (*sdkmcp.ReadResourceResult, error) {
	stored, err := h.service.List(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.URLItem, 0, len(stored))
	for _, item := range stored {
		items = append(items, dto.NewURLItem(h.baseURL, item.ID, item.Original))
	}

	body, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return textResource(req.Params.URI, string(body)), nil
}
