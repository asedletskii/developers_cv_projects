package readers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
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

func TestGetStatisticsReturnsErrorWhenReaderDatabaseIsUnavailable(t *testing.T) {
	reader := NewResumeLinkReader(nil)

	_, err := reader.GetLinkStatistics(context.Background(), queries.LinkAnalyticsQuery{
		LinkID: sharedvalueobjects.NewID(),
	})
	require.Error(t, err)
}
