// Package handlers публикует операции коротких ссылок по Model Context Protocol.
package handlers

import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

const serverVersion = "v1.0.0"

// mcpServer собирает MCP-сервер и регистрирует инструменты и ресурсы.
func (h *URLHandler) mcpServer() *sdkmcp.Server {
	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "shortener",
		Version: serverVersion,
	}, nil)
	h.addShortenTool(server)
	h.addListTool(server)
	h.addResolveTool(server)
	h.addUpdateTool(server)
	h.addDeleteTool(server)
	h.addURLListResource(server)
	h.addURLResource(server)
	return server
}
