package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestAnalyticsClickQueryUsesCurrentLayoutHash(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	layout, err := domain.NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"})
	require.NoError(t, err)
	reader := &clickMapReaderStub{}
	h := NewAnalyticsClickQueryHandler(&resumeLinkReaderStub{view: ResumeLinkPublicView{ResumeID: resumeID.String()}}, &ownershipReaderStub{}, &publicResumeQueryStub{}, layoutSnapshotBuilderStub{snapshot: layout}, reader)

	result, err := h.Handle(context.Background(), LinkClicksQuery{UserID: sharedvalueobjects.NewID(), LinkID: linkID, Section: "about", Breakpoint: "md"})
	require.NoError(t, err)
	require.Equal(t, int64(3), result.Total)
	require.Equal(t, "0123456789abcdef", reader.query.CurrentHash)
}

func TestAnalyticsClickQueryRejectsInvalidBreakpoint(t *testing.T) {
	h := NewAnalyticsClickQueryHandler(&resumeLinkReaderStub{}, &ownershipReaderStub{}, &publicResumeQueryStub{}, layoutSnapshotBuilderStub{}, &clickMapReaderStub{})
	_, err := h.Handle(context.Background(), LinkClicksQuery{UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Section: "about", Breakpoint: "wide"})
	require.ErrorIs(t, err, domain.ErrInvalidAnalyticsClick)
}
