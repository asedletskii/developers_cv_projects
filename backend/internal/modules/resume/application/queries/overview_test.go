package queries

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResumeOverviewViewUsesRootFields(t *testing.T) {
	view := ResumeOverviewView{
		Theme:        &OverviewThemeSettings{Primary: "#E8785A"},
		Greeting:     OverviewGreeting{First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Engineer"},
		Availability: &OverviewAvailability{Status: "open"},
		ResumeLink:   stringPtr("https://example.com/cv.pdf"),
		Sections:     []OverviewSectionSettings{{ID: "about", Enabled: true, Collapsed: false}},
	}

	data, err := json.Marshal(view)
	require.NoError(t, err)

	require.JSONEq(t, `{
		"theme":{"primary":"#E8785A"},
		"greeting":{"first":"Hello","second":"I am","secondHighlight":"Ada","third":"Engineer"},
		"availability":{"status":"open"},
		"resumeLink":"https://example.com/cv.pdf",
		"sections":[{"id":"about","enabled":true,"collapsed":false}]
	}`, string(data))
	require.NotContains(t, string(data), `"presentation"`)
}

func TestResumeOverviewViewEmptySectionsSerializeAsArray(t *testing.T) {
	view := ResumeOverviewView{Sections: make([]OverviewSectionSettings, 0)}

	data, err := json.Marshal(view)
	require.NoError(t, err)
	require.JSONEq(t, `{"greeting":{"first":"","second":"","secondHighlight":"","third":""},"sections":[]}`, string(data))
}

func stringPtr(value string) *string {
	return &value
}
