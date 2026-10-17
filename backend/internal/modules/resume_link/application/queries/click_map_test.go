package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type clickMapLinkReaderStub struct{ view ResumeLinkPublicView }

func (s clickMapLinkReaderStub) GetByID(context.Context, PublicResumeLinkQuery) (ResumeLinkPublicView, error) {
	return s.view, nil
}

type clickMapOwnershipStub struct{}

func (clickMapOwnershipStub) GetOwned(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeView, error) {
	return resumequeries.ResumeView{ID: "resume"}, nil
}

func (clickMapOwnershipStub) OwnsResume(context.Context, resumequeries.GetResumeQuery) (bool, error) {
	return true, nil
}

type clickMapPublicStub struct {
	view resumequeries.PublicResumeView
}

func (s clickMapPublicStub) Handle(context.Context, resumequeries.PublicResumeQuery) (resumequeries.PublicResumeView, error) {
	return s.view, nil
}

type clickMapLayoutStub struct{ snapshot domain.LayoutSnapshot }

func (s clickMapLayoutStub) Build(resumequeries.PublicResumeView) (domain.LayoutSnapshot, error) {
	return s.snapshot, nil
}

type clickMapReaderStub struct{ query LinkClicksQuery }

func (s *clickMapReaderStub) GetClicks(_ context.Context, q LinkClicksQuery) (ClickGrid, error) {
	s.query = q
	return ClickGrid{Section: q.Section, Grid: ClickGridSize, Cells: []ClickCell{{X: 1, Y: 2, Count: 3}}, Total: 3}, nil
}

func TestAnalyticsClickQueryUsesCurrentLayoutHash(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	layout, err := domain.NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"})
	require.NoError(t, err)
	reader := &clickMapReaderStub{}
	h := NewAnalyticsClickQueryHandler(clickMapLinkReaderStub{view: ResumeLinkPublicView{ResumeID: resumeID.String()}}, clickMapOwnershipStub{}, clickMapPublicStub{}, clickMapLayoutStub{snapshot: layout}, reader)

	result, err := h.Handle(context.Background(), LinkClicksQuery{UserID: sharedvalueobjects.NewID(), LinkID: linkID, Section: "about", Breakpoint: "md"})
	require.NoError(t, err)
	require.Equal(t, int64(3), result.Total)
	require.Equal(t, "0123456789abcdef", reader.query.CurrentHash)
}

func TestAnalyticsClickQueryRejectsInvalidBreakpoint(t *testing.T) {
	h := NewAnalyticsClickQueryHandler(clickMapLinkReaderStub{}, clickMapOwnershipStub{}, clickMapPublicStub{}, clickMapLayoutStub{}, &clickMapReaderStub{})
	_, err := h.Handle(context.Background(), LinkClicksQuery{UserID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), Section: "about", Breakpoint: "wide"})
	require.ErrorIs(t, err, domain.ErrInvalidAnalyticsClick)
}
