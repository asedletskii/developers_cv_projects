package middleware

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func requestIDFromBody(t *testing.T, body []byte) string {
	t.Helper()
	var response struct {
		Error struct {
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(body, &response))
	require.NotEmpty(t, response.Error.RequestID)
	return response.Error.RequestID
}
