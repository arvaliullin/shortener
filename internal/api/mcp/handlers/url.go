package handlers

import (
	"net/http"

	"github.com/arvaliullin/shortener/internal/core/ports"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// URLHandler обрабатывает запросы к коротким ссылкам.
type URLHandler struct {
	service ports.URLService
	baseURL string
	http    http.Handler
}

// NewURLHandler собирает обработчик коротких ссылок.
func NewURLHandler(service ports.URLService, baseURL string) *URLHandler {
	h := &URLHandler{
		service: service,
		baseURL: baseURL,
	}
	server := h.mcpServer()
	h.http = sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server {
		return server
	}, &sdkmcp.StreamableHTTPOptions{Stateless: true})
	return h
}

// ServeHTTP отдает обработчик по stateless Streamable HTTP.
func (h *URLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.http.ServeHTTP(w, r)
}
