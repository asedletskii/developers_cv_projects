package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatedIDResponseJSON(t *testing.T) {
	payload, err := json.Marshal(CreatedIDResponse{ID: "resume-123"})
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"resume-123"}`, string(payload))
}
