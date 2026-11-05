package mcp

import (
	"encoding/json"
	"strconv"
)

const maxListedURLs = 20

// briefURLList добавляет число ссылок и оставляет модели не больше двадцати элементов.
func briefURLList(text string) string {
	var body map[string]json.RawMessage
	if err := json.Unmarshal([]byte(text), &body); err != nil {
		return text
	}
	raw, ok := body["urls"]
	if !ok {
		return text
	}

	var urls []json.RawMessage
	if err := json.Unmarshal(raw, &urls); err != nil {
		return text
	}

	count := len(urls)
	truncated := count > maxListedURLs
	if truncated {
		urls = urls[:maxListedURLs]
	}
	limited, err := json.Marshal(urls)
	if err != nil {
		return text
	}

	body["urls"] = limited
	body["count"] = json.RawMessage(strconv.Itoa(count))
	if truncated {
		body["truncated"] = json.RawMessage("true")
	}
	out, err := json.Marshal(body)
	if err != nil {
		return text
	}
	return string(out)
}
