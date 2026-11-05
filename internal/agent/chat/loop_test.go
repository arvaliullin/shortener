package chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type scriptedModel struct {
	steps []Completion
	seen  [][]string
}

func (m *scriptedModel) Complete(_ context.Context, _ []Message, tools []Tool) (Completion, error) {
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.Name)
	}
	m.seen = append(m.seen, names)
	step := m.steps[0]
	m.steps = m.steps[1:]
	return step, nil
}

type fakeTools struct {
	tools  []Tool
	called []string
}

func (f *fakeTools) List(context.Context) ([]Tool, error) {
	return f.tools, nil
}

func (f *fakeTools) Call(_ context.Context, name, _ string) (string, error) {
	f.called = append(f.called, name)
	return `{"ok":true}`, nil
}

func TestReplyReturnsText(t *testing.T) {
	t.Parallel()

	model := &scriptedModel{steps: []Completion{{Content: "готово"}}}
	tools := &fakeTools{tools: []Tool{{Name: "list_urls"}, {Name: "delete_url"}}}
	reply, err := NewLoop(model, tools).Reply(context.Background(), []Message{{Role: "user", Content: "покажи ссылки"}})

	require.NoError(t, err)
	require.Equal(t, "готово", reply)
	require.Equal(t, [][]string{{"list_urls"}}, model.seen)
}

func TestReplyCallsToolThenReturnsText(t *testing.T) {
	t.Parallel()

	model := &scriptedModel{steps: []Completion{
		{ToolCalls: []ToolCall{{ID: "call_1", Name: "list_urls", Arguments: "{}"}}},
		{Content: "ссылок нет"},
	}}
	tools := &fakeTools{tools: []Tool{{Name: "list_urls"}}}
	reply, err := NewLoop(model, tools).Reply(context.Background(), []Message{{Role: "user", Content: "покажи ссылки"}})

	require.NoError(t, err)
	require.Equal(t, "ссылок нет", reply)
	require.Equal(t, []string{"list_urls"}, tools.called)
}

func TestReplyHidesMutationUntilAsked(t *testing.T) {
	t.Parallel()

	model := &scriptedModel{steps: []Completion{
		{ToolCalls: []ToolCall{{ID: "call_1", Name: "delete_url", Arguments: `{"id":"abc"}`}}},
		{Content: "не стал удалять"},
	}}
	tools := &fakeTools{tools: []Tool{{Name: "list_urls"}, {Name: "delete_url"}}}
	reply, err := NewLoop(model, tools).Reply(context.Background(), []Message{{Role: "user", Content: "покажи ссылки"}})

	require.NoError(t, err)
	require.Equal(t, "не стал удалять", reply)
	require.Empty(t, tools.called)
	require.Equal(t, [][]string{{"list_urls"}, {"list_urls"}}, model.seen)
}

func TestReplyAllowsDeleteWhenAsked(t *testing.T) {
	t.Parallel()

	model := &scriptedModel{steps: []Completion{
		{ToolCalls: []ToolCall{{ID: "call_1", Name: "delete_url", Arguments: `{"id":"abc"}`}}},
		{Content: "удалил"},
	}}
	tools := &fakeTools{tools: []Tool{{Name: "delete_url"}}}
	reply, err := NewLoop(model, tools).Reply(context.Background(), []Message{{Role: "user", Content: "удали ссылку abc"}})

	require.NoError(t, err)
	require.Equal(t, "удалил", reply)
	require.Equal(t, []string{"delete_url"}, tools.called)
	require.Equal(t, [][]string{{"delete_url"}, {"delete_url"}}, model.seen)
}

func TestReplyRejectsEmptyCompletion(t *testing.T) {
	t.Parallel()

	model := &scriptedModel{steps: []Completion{{}}}
	_, err := NewLoop(model, &fakeTools{}).Reply(context.Background(), []Message{{Role: "user", Content: "привет"}})

	require.ErrorIs(t, err, errEmptyResponse)
}

func TestReplyStopsAfterTooManySteps(t *testing.T) {
	t.Parallel()

	call := Completion{ToolCalls: []ToolCall{{ID: "call_1", Name: "list_urls", Arguments: "{}"}}}
	model := &scriptedModel{steps: []Completion{call, call, call}}
	tools := &fakeTools{tools: []Tool{{Name: "list_urls"}}}
	_, err := NewLoop(model, tools).Reply(context.Background(), []Message{{Role: "user", Content: "покажи ссылки"}})

	require.ErrorIs(t, err, errTooManySteps)
}
