package queries

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResumeOverviewViewKeepsRootFields(t *testing.T) {
	view := ResumeOverviewView{
		Theme:        &OverviewThemeSettings{Primary: "#E8785A"},
		Greeting:     OverviewGreeting{First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Engineer"},
		Availability: &OverviewAvailability{Status: "open"},
		Sections:     []OverviewSectionSettings{{ID: "about", Enabled: true, Collapsed: false, Style: stringPtr("second")}},
	}

	require.Equal(t, "#E8785A", view.Theme.Primary)
	require.Equal(t, "Hello", view.Greeting.First)
	require.Equal(t, "open", view.Availability.Status)
	require.Len(t, view.Sections, 1)
	require.Equal(t, "about", view.Sections[0].ID)
	require.Equal(t, stringPtr("second"), view.Sections[0].Style)

	viewType := reflect.TypeOf(view)
	for _, fieldName := range []string{"Presentation", "ResumeLink"} {
		_, exists := viewType.FieldByName(fieldName)
		require.False(t, exists, "overview contract must not contain legacy field %s", fieldName)
	}
}

func TestResumeOverviewViewKeepsEmptySectionsNonNil(t *testing.T) {
	view := ResumeOverviewView{Sections: make([]OverviewSectionSettings, 0)}

	require.NotNil(t, view.Sections)
	require.Empty(t, view.Sections)
}

func TestResumeOverviewViewJSONUsesLowerCamelCaseFields(t *testing.T) {
	from := "2026-10-01"
	note := "Looking for a senior role"
	view := ResumeOverviewView{
		Theme:        &OverviewThemeSettings{Primary: "#E8785A", Secondary: "#428BF9"},
		Greeting:     OverviewGreeting{First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Engineer"},
		Availability: &OverviewAvailability{Status: "open", From: &from, Note: &note},
		Sections:     []OverviewSectionSettings{{ID: "about", Enabled: true, Collapsed: false, Style: stringPtr("second")}},
	}

	payload, err := json.Marshal(view)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"theme": {"primary": "#E8785A", "secondary": "#428BF9"},
		"greeting": {"first": "Hello", "second": "I am", "secondHighlight": "Ada", "third": "Engineer"},
		"availability": {"status": "open", "from": "2026-10-01", "note": "Looking for a senior role"},
		"sections": [{"id": "about", "enabled": true, "collapsed": false, "style": "second"}]
	}`, string(payload))
}

func TestResumeOverviewViewJSONOmitsEmptyOptionalFields(t *testing.T) {
	view := ResumeOverviewView{
		Greeting: OverviewGreeting{},
		Sections: []OverviewSectionSettings{{ID: "courses", Enabled: true}},
	}

	payload, err := json.Marshal(view)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"greeting": {"first": "", "second": "", "secondHighlight": "", "third": ""},
		"sections": [{"id": "courses", "enabled": true, "collapsed": false}]
	}`, string(payload))
}
