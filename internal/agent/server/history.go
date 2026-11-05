package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/arvaliullin/shortener/internal/agent/chat"
)

const (
	maxMessages = 40
	maxRunes    = 4000
)

var errBadRequest = errors.New("bad request")

type chatRequest struct {
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// readHistory принимает историю страницы: последнее сообщение должно быть от пользователя.
func readHistory(w http.ResponseWriter, r *http.Request) ([]chat.Message, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var body chatRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, errBadRequest
	}
	if len(body.Messages) == 0 || len(body.Messages) > maxMessages {
		return nil, errBadRequest
	}

	history := make([]chat.Message, 0, len(body.Messages))
	for _, message := range body.Messages {
		content := strings.TrimSpace(message.Content)
		if (message.Role != "user" && message.Role != "assistant") || content == "" || utf8.RuneCountInString(content) > maxRunes {
			return nil, errBadRequest
		}
		history = append(history, chat.Message{Role: message.Role, Content: content})
	}
	if history[len(history)-1].Role != "user" {
		return nil, errBadRequest
	}
	return history, nil
}
