package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/arvaliullin/shortener/internal/api/mcp/dto"
	"github.com/arvaliullin/shortener/internal/core/domain"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	urlResourcePrefix   = "shortener://urls/"
	urlResourceTemplate = "shortener://urls/{id}"
)

// addURLResource регистрирует чтение одной короткой ссылки.
func (h *URLHandler) addURLResource(server *sdkmcp.Server) {
	server.AddResourceTemplate(&sdkmcp.ResourceTemplate{
		URITemplate: urlResourceTemplate,
		Name:        "url",
		Description: "Сохраненная короткая ссылка по идентификатору.",
		MIMEType:    "application/json",
	}, h.readURL)
}

// readURL возвращает JSON одной короткой ссылки.
func (h *URLHandler) readURL(ctx context.Context, req *sdkmcp.ReadResourceRequest) (*sdkmcp.ReadResourceResult, error) {
	id, ok := urlID(req.Params.URI)
	if !ok {
		return nil, sdkmcp.ResourceNotFoundError(req.Params.URI)
	}

	originalURL, err := h.service.Resolve(ctx, id)
	if errors.Is(err, domain.ErrURLNotFound) {
		return nil, sdkmcp.ResourceNotFoundError(req.Params.URI)
	}
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(dto.NewURLItem(h.baseURL, id, originalURL))
	if err != nil {
		return nil, err
	}

	return textResource(req.Params.URI, string(body)), nil
}

// urlID извлекает идентификатор из URI шаблона shortener://urls/{id}.
func urlID(uri string) (string, bool) {
	id, ok := strings.CutPrefix(uri, urlResourcePrefix)
	if !ok || id == "" || strings.Contains(id, "/") {
		return "", false
	}
	return id, true
}

// textResource собирает текстовое JSON-содержимое ресурса.
func textResource(uri, body string) *sdkmcp.ReadResourceResult {
	return &sdkmcp.ReadResourceResult{
		Contents: []*sdkmcp.ResourceContents{{
			URI:      uri,
			MIMEType: "application/json",
			Text:     body,
		}},
	}
}
