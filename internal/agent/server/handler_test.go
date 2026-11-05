package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arvaliullin/shortener/internal/agent/chat"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type textModel struct {
	content string
}

func (m textModel) Complete(context.Context, []chat.Message, []chat.Tool) (chat.Completion, error) {
	return chat.Completion{Content: m.content}, nil
}

type memorySession struct{}

func (memorySession) List(context.Context) ([]chat.Tool, error) {
	return nil, nil
}

func (memorySession) Call(context.Context, string, string) (string, error) {
	return "", nil
}

func (memorySession) Close() error {
	return nil
}

type memoryOpener struct{}

func (memoryOpener) Open(context.Context) (Session, error) {
	return memorySession{}, nil
}

func TestChatReturnsText(t *testing.T) {
	t.Parallel()

	handler := NewHandler(textModel{content: "готово"}, memoryOpener{}, zerolog.Nop())
	request := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"messages":[{"role":"user","content":"привет"}]}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"content":"готово"}`, recorder.Body.String())
}

func TestChatRejectsHistoryWithoutUser(t *testing.T) {
	t.Parallel()

	handler := NewHandler(textModel{content: "готово"}, memoryOpener{}, zerolog.Nop())
	request := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"messages":[{"role":"assistant","content":"привет"}]}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
