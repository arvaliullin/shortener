package chat

import (
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	toolCallMarkup     = regexp.MustCompile(`(?s)<tool_call>\s*(.*?)\s*</tool_call>`)
	functionNameMarkup = regexp.MustCompile(`<function_name>\s*([A-Za-z0-9_]+)\s*</function_name>`)
)

type markupSpan struct {
	start int
	end   int
	call  ToolCall
}

// toolCallsFromContent достает вызовы из текста и убирает их разметку.
func toolCallsFromContent(content string, offered []Tool) ([]ToolCall, string) {
	spans := acceptedSpans(append(toolCallSpans(content, offered), functionNameSpans(content, offered)...))
	if len(spans) == 0 {
		return nil, content
	}

	calls := make([]ToolCall, 0, len(spans))
	for i, span := range spans {
		span.call.ID = "call_" + strconv.Itoa(i+1)
		calls = append(calls, span.call)
	}
	return calls, withoutSpans(content, spans)
}

// toolCallSpans читает JSON внутри <tool_call>.
func toolCallSpans(content string, offered []Tool) []markupSpan {
	var spans []markupSpan
	for _, match := range toolCallMarkup.FindAllStringSubmatchIndex(content, -1) {
		call, ok := toolCallFromJSON(content[match[2]:match[3]], offered)
		if !ok {
			continue
		}
		spans = append(spans, markupSpan{start: match[0], end: match[1], call: call})
	}
	return spans
}

// functionNameSpans читает имя внутри <function_name>.
func functionNameSpans(content string, offered []Tool) []markupSpan {
	var spans []markupSpan
	for _, match := range functionNameMarkup.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		if !knownTool(offered, name) {
			continue
		}
		spans = append(spans, markupSpan{
			start: match[0],
			end:   match[1],
			call:  ToolCall{Name: name, Arguments: "{}"},
		})
	}
	return spans
}

// toolCallFromJSON разбирает объект name и arguments.
func toolCallFromJSON(body string, offered []Tool) (ToolCall, bool) {
	var parsed struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return ToolCall{}, false
	}
	if !knownTool(offered, parsed.Name) {
		return ToolCall{}, false
	}
	return ToolCall{Name: parsed.Name, Arguments: callArguments(parsed.Arguments)}, true
}

// callArguments оставляет аргументы одной JSON-строкой.
func callArguments(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return "{}"
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		if strings.TrimSpace(text) == "" {
			return "{}"
		}
		return text
	}
	return strings.TrimSpace(string(raw))
}

// acceptedSpans оставляет непересекающиеся фрагменты в порядке текста.
func acceptedSpans(spans []markupSpan) []markupSpan {
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].start < spans[j].start
	})

	kept := make([]markupSpan, 0, len(spans))
	cursor := 0
	for _, span := range spans {
		if span.start < cursor {
			continue
		}
		kept = append(kept, span)
		cursor = span.end
	}
	return kept
}

// withoutSpans убирает распознанную разметку из ответа.
func withoutSpans(content string, spans []markupSpan) string {
	var builder strings.Builder
	cursor := 0
	for _, span := range spans {
		builder.WriteString(content[cursor:span.start])
		cursor = span.end
	}
	builder.WriteString(content[cursor:])
	return strings.TrimSpace(builder.String())
}
