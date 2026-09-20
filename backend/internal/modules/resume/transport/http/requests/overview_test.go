package requests

import (
	"bytes"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestOverviewRequestAcceptsRootFields(t *testing.T) {
	command, err := (ResumeOverviewRequest{
		Version:  7,
		Theme:    &OverviewThemeRequest{Primary: stringPtr("#E8785A")},
		Greeting: OverviewGreetingRequest{First: "A", Second: "B", SecondHighlight: "C", Third: "D"},
		Availability: &OverviewAvailabilityRequest{
			Status: "open",
			From:   stringPtr("2026-10-01"),
			Note:   stringPtr("Remote only"),
		},
		Sections: []OverviewSectionRequest{{ID: "about", Enabled: true, Collapsed: false, Style: stringPtr("second")}},
	}).ToCommand(sharedvalueobjects.NewID(), sharedvalueobjects.NewID())

	require.NoError(t, err)
	require.Equal(t, int64(7), command.Version)
	require.NotNil(t, command.Input.Theme)
	require.NotNil(t, command.Input.Availability)
	require.Equal(t, "open", string(command.Input.Availability.Status))
	require.Len(t, command.Input.Sections, 1)
	require.Equal(t, stringPtr("second"), command.Input.Sections[0].Style)
}

func TestOverviewRequestJSONUsesRootFields(t *testing.T) {
	payload := []byte(`{
		"version":7,
		"theme":{"primary":"#E8785A"},
		"greeting":{"first":"A","second":"B","secondHighlight":"C","third":"D"},
		"availability":{"status":"open"},
		"sections":[{"id":"about","enabled":true,"collapsed":false,"style":"second"}]
	}`)

	request, err := decodeStrictResumeOverview(payload)

	require.NoError(t, err)
	require.NotNil(t, request.Theme)
	require.Len(t, request.Sections, 1)
	require.Equal(t, stringPtr("second"), request.Sections[0].Style)
}

func TestOverviewRequestAllowsSectionWithoutStyle(t *testing.T) {
	payload := []byte(`{
		"version":7,
		"greeting":{"first":"A","second":"B","secondHighlight":"C","third":"D"},
		"sections":[{"id":"courses","enabled":true,"collapsed":false}]
	}`)

	request, err := decodeStrictResumeOverview(payload)

	require.NoError(t, err)
	require.Len(t, request.Sections, 1)
	require.Nil(t, request.Sections[0].Style)
}

func TestOverviewRequestRejectsRemovedResumeLinkField(t *testing.T) {
	payload := []byte(`{
		"version":7,
		"greeting":{"first":"A","second":"B","secondHighlight":"C","third":"D"},
		"resumeLink":"https://example.com/cv.pdf",
		"sections":[]
	}`)

	_, err := decodeStrictResumeOverview(payload)

	require.Error(t, err)
}

func TestOverviewRequestJSONRejectsLegacyPresentation(t *testing.T) {
	payload := []byte(`{
		"presentation":{
			"greeting":{"first":"A","second":"B","secondHighlight":"C","third":"D"},
			"sections":[]
		}
	}`)

	_, err := decodeStrictResumeOverview(payload)

	require.Error(t, err)
}

func decodeStrictResumeOverview(payload []byte) (ResumeOverviewRequest, error) {
	var request ResumeOverviewRequest
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return ResumeOverviewRequest{}, err
	}
	return request, nil
}

func stringPtr(value string) *string { return &value }
