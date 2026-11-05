// Package server принимает POST /api/chat и отвечает текстом агента.
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arvaliullin/shortener/internal/agent/chat"
	"github.com/rs/zerolog"
)

const chatTimeout = 180 * time.Second

// Session закрывает набор инструментов после ответа.
type Session interface {
	chat.Tools
	Close() error
}

// Opener открывает MCP на время одного сообщения.
type Opener interface {
	Open(context.Context) (Session, error)
}

// Handler отвечает на сообщение чата.
type Handler struct {
	model  chat.Model
	tools  Opener
	logger zerolog.Logger
}

// NewHandler собирает обработчик чата.
func NewHandler(model chat.Model, tools Opener, logger zerolog.Logger) *Handler {
	return &Handler{model: model, tools: tools, logger: logger}
}

// ServeHTTP читает историю, спрашивает модель и возвращает текст.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	history, err := readHistory(w, r)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), chatTimeout)
	defer cancel()

	session, err := h.tools.Open(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("open mcp")
		http.Error(w, "chat failed", http.StatusBadGateway)
		return
	}
	defer func() {
		if err := session.Close(); err != nil {
			h.logger.Error().Err(err).Msg("close mcp")
		}
	}()

	reply, err := chat.NewLoop(h.model, session).Reply(ctx, history)
	if err != nil {
		h.logger.Error().Err(err).Msg("chat reply")
		http.Error(w, "chat failed", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse{Content: reply})
}

type chatResponse struct {
	Content string `json:"content"`
}
