package handlers_test

import (
	"testing"

	"github.com/arvaliullin/shortener/internal/core/ports/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteURL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := mocks.NewMockURLService(ctrl)
	svc.EXPECT().
		Delete(gomock.Any(), "abc12345").
		Return(nil)

	session := connect(t, svc)
	output := decodeOutput[struct {
		ID string `json:"id"`
	}](t, callTool(t, session, "delete_url", map[string]any{"id": "abc12345"}))

	require.Equal(t, "abc12345", output.ID)
}
