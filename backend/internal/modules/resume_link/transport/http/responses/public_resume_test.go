package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/application/queries"
)

func TestNewPublicResumeResponseNormalizesRequiredSlices(t *testing.T) {
	payload, err := json.Marshal(NewPublicResumeResponse(queries.PublicResumeResult{}))
	require.NoError(t, err)
	require.NotContains(t, string(payload), "null")

	var object map[string]any
	require.NoError(t, json.Unmarshal(payload, &object))
	for _, field := range []string{"languages", "technologies", "experiences", "educations", "courses", "activities", "projects"} {
		require.Equal(t, []any{}, object[field], "field %s must be an empty array", field)
	}
	overview, ok := object["overview"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, []any{}, overview["sections"])
	require.NotContains(t, object, "profile")
	require.NotContains(t, object, "contacts")
	require.NotContains(t, object, "skills")
}
