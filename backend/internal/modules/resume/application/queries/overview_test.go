package queries

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResumeOverviewViewKeepsRootFields(t *testing.T) {
	view := ResumeOverviewView{
		Theme:        &OverviewThemeSettings{Primary: "#E8785A"},
		Greeting:     OverviewGreeting{First: "Hello", Second: "I am", SecondHighlight: "Ada", Third: "Engineer"},
		Availability: &OverviewAvailability{Status: "open"},
		Sections:     []OverviewSectionSettings{{ID: "about", Enabled: true, Collapsed: false}},
	}

	require.Equal(t, "#E8785A", view.Theme.Primary)
	require.Equal(t, "Hello", view.Greeting.First)
	require.Equal(t, "open", view.Availability.Status)
	require.Len(t, view.Sections, 1)
	require.Equal(t, "about", view.Sections[0].ID)

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
