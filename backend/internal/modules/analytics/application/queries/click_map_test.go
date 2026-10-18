package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestAnalyticsClickQueryUsesCurrentLayoutHash(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	layout, err := domain.NewLayoutSnapshot([]domain.AnalyticsSectionID{domain.AnalyticsSectionAbout}, map[domain.AnalyticsSectionID]string{domain.AnalyticsSectionAbout: "0123456789abcdef"})
	require.NoError(t, err)
	reader := &clickMapReaderStub{}
	h := NewAnalyticsClickQueryHandler(
		&analyticsAccessReaderStub{resumeID: resumeID, layout: ports.ResumeAnalyticsLayout{Sections: []ports.ResumeAnalyticsSection{{ID: "about", Hash: "0123456789abcdef"}}}},
		layoutSnapshotBuilderStub{snapshot: layout},
		reader,
	)

	result, err := h.Handle(context.Background(), LinkClicksQuery{UserID: sharedvalueobjects.NewID(), LinkID: linkID, Section: "about", Breakpoint: "md"})
	require.NoError(t, err)
	require.Equal(t, int64(3), result.Total)
	require.Equal(t, "0123456789abcdef", reader.query.CurrentHash)
}

func TestAnalyticsClickQueryRejectsInvalidBreakpoint(t *testing.T) {
	h := NewAnalyticsClickQueryHandler(&analyticsAccessReaderStub{}, layoutSnapshotBuilderStub{}, &clickMapReaderStub{})
	_, err := h.Handle(context.Background(), LinkClicksQuery{UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Section: "about", Breakpoint: "wide"})
	require.ErrorIs(t, err, domain.ErrInvalidAnalyticsClick)
}
