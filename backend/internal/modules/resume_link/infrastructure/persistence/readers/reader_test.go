package readers

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
)

func TestMapView(t *testing.T) {
	model := gormmodel.ResumeLink{
		ID:       testUUID("10000000-0000-0000-0000-000000000001"),
		ResumeID: testUUID("20000000-0000-0000-0000-000000000002"),
		Name:     "GitHub", Status: "active", Version: 2,
	}

	view := mapView(model)

	require.Equal(t, model.ID.String(), view.ID)
	require.Equal(t, model.ResumeID.String(), view.ResumeID)
	require.Equal(t, "active", view.Status)
}

func testUUID(value string) uuid.UUID {
	return uuid.MustParse(value)
}

func TestResumeLinkReaderDoesNotExposeAnalyticsQueries(t *testing.T) {
	typ := reflect.TypeOf(&ResumeLinkReader{})
	for _, method := range []string{"GetLinkStatistics", "GetDaily", "GetVisits", "GetSections", "GetTechnologies", "GetResumeStatistics", "GetClicks"} {
		if _, ok := typ.MethodByName(method); ok {
			t.Fatalf("link reader must not expose analytics query %s", method)
		}
	}
}
