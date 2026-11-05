package mcp

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBriefURLListKeepsCountAndTwenty(t *testing.T) {
	t.Parallel()

	urls := make([]map[string]string, 21)
	for i := range urls {
		urls[i] = map[string]string{"id": fmt.Sprintf("id%02d", i)}
	}
	body, err := json.Marshal(map[string]any{"urls": urls})
	require.NoError(t, err)

	var parsed struct {
		Count     int               `json:"count"`
		Truncated bool              `json:"truncated"`
		URLs      []json.RawMessage `json:"urls"`
	}
	require.NoError(t, json.Unmarshal([]byte(briefURLList(string(body))), &parsed))
	require.Equal(t, 21, parsed.Count)
	require.True(t, parsed.Truncated)
	require.Len(t, parsed.URLs, 20)
	kept, err := json.Marshal(parsed.URLs)
	require.NoError(t, err)
	require.Contains(t, string(kept), "id19")
	require.NotContains(t, string(kept), "id20")
}

func TestBriefURLListLeavesOtherTools(t *testing.T) {
	t.Parallel()

	raw := `{"id":"abc","result":"http://localhost/abc"}`
	require.Equal(t, raw, briefURLList(raw))
}
