package overview

import (
	"bytes"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestOverviewRequestAcceptsRootFields(t *testing.T) {
	command, err := (overviewRequest{
		Version:  7,
		Theme:    &themeRequest{Primary: stringPtr("#E8785A")},
		Greeting: greetingRequest{First: "A", Second: "B", SecondHighlight: "C", Third: "D"},
		Availability: &availabilityRequest{
			Status: "open",
			From:   stringPtr("2026-10-01"),
			Note:   stringPtr("Remote only"),
		},
		ResumeLink: stringPtr("https://example.com/cv.pdf"),
		Sections:   []sectionRequest{{ID: "about", Enabled: true, Collapsed: false}},
	}).ToCommand(identity.NewID(), identity.NewID())

	require.NoError(t, err)
	require.Equal(t, int64(7), command.Version)
	require.NotNil(t, command.Input.Theme)
	require.NotNil(t, command.Input.Availability)
	require.NotNil(t, command.Input.ResumeLink)
	require.Equal(t, "https://example.com/cv.pdf", command.Input.ResumeLink.String())
	require.Equal(t, "open", string(command.Input.Availability.Status))
	require.Len(t, command.Input.Sections, 1)
}

func TestOverviewRequestRejectsRelativeResumePath(t *testing.T) {
	_, err := (overviewRequest{
		Greeting:   greetingRequest{First: "A", Second: "B", SecondHighlight: "C", Third: "D"},
		ResumeLink: stringPtr("/cv.pdf"),
		Sections:   []sectionRequest{},
	}).ToCommand(identity.NewID(), identity.NewID())

	require.Error(t, err)
}

func TestOverviewRequestJSONUsesRootFields(t *testing.T) {
	payload := []byte(`{
		"version":7,
		"theme":{"primary":"#E8785A"},
		"greeting":{"first":"A","second":"B","secondHighlight":"C","third":"D"},
		"availability":{"status":"open"},
		"resumeLink":"https://example.com/cv.pdf",
		"sections":[]
	}`)

	request, err := decodeStrict(payload)

	require.NoError(t, err)
	require.NotNil(t, request.Theme)
	require.Equal(t, "https://example.com/cv.pdf", *request.ResumeLink)
	require.Empty(t, request.Sections)
}

func TestOverviewRequestJSONRejectsLegacyPresentation(t *testing.T) {
	payload := []byte(`{
		"presentation":{
			"greeting":{"first":"A","second":"B","secondHighlight":"C","third":"D"},
			"sections":[]
		}
	}`)

	_, err := decodeStrict(payload)

	require.Error(t, err)
}

func decodeStrict(payload []byte) (overviewRequest, error) {
	var request overviewRequest
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return overviewRequest{}, err
	}
	return request, nil
}

func stringPtr(value string) *string {
	return &value
}
