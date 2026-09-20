package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestPublicResumeLinkQueryResolvesActiveLinkToResume(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive)}}
	resume := &publicResumeQueryStub{view: resumecontracts.PublicResumeView{Title: "CV"}}

	view, err := NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(links), resume).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.NoError(t, err)
	require.Equal(t, "CV", view.Title)
	require.Equal(t, linkID, links.query.LinkID)
	require.Equal(t, resumeID, resume.query.ID)
}

func TestPublicResumeLinkQueryRejectsInactiveLink(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusInactive),
	}}
	resume := &publicResumeQueryStub{}

	_, err := NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(links), resume).Handle(context.Background(), PublicResumeLinkQuery{LinkID: linkID})

	require.ErrorIs(t, err, domain.ErrResumeLinkInactive)
	require.True(t, resume.query.ID.IsZero(), "inactive links must not trigger a public resume read")
}

func TestPublicResumeLinkQueryDoesNotRecordFailedLookup(t *testing.T) {
	links := &resumeLinkReaderStub{err: application.ErrResumeLinkNotFound}
	resume := &publicResumeQueryStub{}

	_, err := NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(links), resume).Handle(context.Background(), PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, application.ErrResumeLinkNotFound)
	require.True(t, resume.query.ID.IsZero(), "failed link lookups must not trigger a public resume read")
}

func TestPublicResumeLinkQueryRejectsInvalidResolvedResumeID(t *testing.T) {
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: "not-a-uuid",
		Status:   string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{}

	_, err := NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(links), resume).Handle(
		context.Background(),
		PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()},
	)

	require.Error(t, err)
	require.True(t, resume.query.ID.IsZero())
}

func TestPublicResumeLinkQueryTranslatesDomainNotFound(t *testing.T) {
	links := &resumeLinkReaderStub{err: domain.ErrEntityNotFound}
	handler := NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(links), &publicResumeQueryStub{})
	_, err := handler.Handle(context.Background(), PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()})
	require.ErrorIs(t, err, application.ErrResumeLinkNotFound)
}

func TestPublicResumeLinkNormalizesProviderResumeNotFound(t *testing.T) {
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ResumeID: sharedvalueobjects.NewID().String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{err: resumecontracts.ErrResumeNotFound}
	handler := NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(links), resume)

	_, err := handler.Handle(context.Background(), PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()})

	require.ErrorIs(t, err, application.ErrResumeNotFound)
	require.ErrorIs(t, err, resumecontracts.ErrResumeNotFound)
}

func TestPublicResumeLinkClearsUnverifiedOwnerWhenOwnershipReaderIsAbsent(t *testing.T) {
	visit := &publicVisitCreatorStub{}
	handler, query, _ := newPublicAnalyticsTestHandler(nil, visit, nil)
	originalOwnerID := query.Visit.Visitor.OwnerID
	require.False(t, originalOwnerID.IsZero())

	_, _, err := handler.HandleWithVisit(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, 1, visit.calls)
	require.True(t, visit.input.Visitor.OwnerID.IsZero(), "unverified owner identity must not be used for analytics classification")
}

func TestPublicResumeLinkWithoutAnalyticsStillServesPublicResume(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}
	resume := &publicResumeQueryStub{view: resumecontracts.PublicResumeView{Title: "Public CV"}}
	handler := NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(links), resume)

	view, identity, err := handler.HandleWithVisit(context.Background(), PublicResumeLinkQuery{
		LinkID: linkID,
		Visit:  &analyticscontracts.CreatePublicVisitInput{LinkID: linkID},
	})

	require.NoError(t, err)
	require.Equal(t, "Public CV", view.Title)
	require.True(t, identity.ViewID.IsZero())
}

func TestPublicResumeLinkRejectsMissingRequiredDependenciesWithoutPanic(t *testing.T) {
	query := PublicResumeLinkQuery{LinkID: sharedvalueobjects.NewID()}

	tests := []struct {
		name    string
		handler *PublicResumeLinkQueryHandler
	}{
		{name: "nil handler", handler: nil},
		{name: "nil link resolver", handler: NewPublicResumeLinkQueryHandler(nil, &publicResumeQueryStub{})},
		{name: "nil resume query", handler: NewPublicResumeLinkQueryHandler(NewActivePublicResumeLinkResolver(&resumeLinkReaderStub{}), nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			require.NotPanics(t, func() {
				_, err = tt.handler.Handle(context.Background(), query)
			})
			require.ErrorContains(t, err, "dependencies are nil")
		})
	}
}

func TestPublicResumeLinkReportsUnexpectedOwnerResolutionFailureAndContinues(t *testing.T) {
	ownerErr := errors.New("owner storage unavailable")
	owner := &ownershipReaderStub{err: ownerErr}
	identity := analyticscontracts.VisitorIdentityResult{ViewID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(), SetCookie: true}
	visit := &publicVisitCreatorStub{identity: identity}
	var reported []PublicAnalyticsFailure
	handler, query, _ := newPublicAnalyticsTestHandler(owner, visit, func(failure PublicAnalyticsFailure) {
		reported = append(reported, failure)
	})
	linkID := query.LinkID
	originalVisitInput := *query.Visit

	view, actualIdentity, err := handler.HandleWithVisit(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, "Public CV", view.Title)
	require.Len(t, visit.input.ResumeLayout.Sections, 1)
	require.Equal(t, "about", visit.input.ResumeLayout.Sections[0].ID)
	require.Len(t, visit.input.ResumeLayout.Sections[0].Hash, 16)
	require.Equal(t, []string{"Go"}, visit.input.ResumeLayout.Technologies)
	require.Equal(t, identity, actualIdentity)
	require.Equal(t, 1, owner.calls)
	require.Equal(t, 1, visit.calls)
	require.True(t, visit.input.Visitor.OwnerID.IsZero(), "unverified owner identity must be cleared")
	require.Equal(t, originalVisitInput, *query.Visit, "the handler must not mutate caller-owned visit input")
	require.Equal(t, originalVisitInput.Visitor.OwnerID, query.Visit.Visitor.OwnerID)
	require.False(t, visit.input.ResumeID.IsZero())
	require.Len(t, reported, 1)
	require.Equal(t, PublicAnalyticsStageOwnerResolution, reported[0].Stage)
	require.Equal(t, linkID, reported[0].LinkID)
	require.Equal(t, visit.input.ResumeID, reported[0].ResumeID)
	require.ErrorIs(t, reported[0].Err, ownerErr)
}

func TestPublicResumeLinkClearsOwnerWhenOwnershipResultDoesNotMatchResume(t *testing.T) {
	owner := &ownershipReaderStub{viewID: sharedvalueobjects.NewID().String()}
	visit := &publicVisitCreatorStub{}
	var reportCount int
	handler, query, _ := newPublicAnalyticsTestHandler(owner, visit, func(PublicAnalyticsFailure) {
		reportCount++
	})

	_, _, err := handler.HandleWithVisit(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, 1, owner.calls)
	require.Zero(t, reportCount, "a negative ownership result is expected, not an analytics failure")
	require.True(t, visit.input.Visitor.OwnerID.IsZero(), "mismatched ownership result must not establish ownership")
}

func TestPublicResumeLinkVisitFailureIsReportedWithoutFailingResumeResponse(t *testing.T) {
	visitErr := errors.New("analytics database unavailable")
	visit := &publicVisitCreatorStub{err: visitErr, identity: analyticscontracts.VisitorIdentityResult{ViewID: sharedvalueobjects.NewID()}}
	var reported []PublicAnalyticsFailure
	handler, query, _ := newPublicAnalyticsTestHandler(nil, visit, func(failure PublicAnalyticsFailure) {
		reported = append(reported, failure)
	})
	originalVisitInput := *query.Visit

	view, identity, err := handler.HandleWithVisit(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, "Public CV", view.Title)
	require.True(t, identity.ViewID.IsZero(), "failed visit must not expose visitor identity")
	require.Equal(t, originalVisitInput, *query.Visit, "the handler must not mutate caller-owned visit input")
	require.Len(t, reported, 1)
	require.Equal(t, PublicAnalyticsStageVisitCreation, reported[0].Stage)
	require.Equal(t, query.LinkID, reported[0].LinkID)
	require.Equal(t, visit.input.ResumeID, reported[0].ResumeID)
	require.ErrorIs(t, reported[0].Err, visitErr)
}

func TestPublicResumeLinkSuccessfulVisitPreservesIdentity(t *testing.T) {
	identity := analyticscontracts.VisitorIdentityResult{ViewID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(), SetCookie: true}
	visit := &publicVisitCreatorStub{identity: identity}
	var reportCount int
	handler, query, _ := newPublicAnalyticsTestHandler(nil, visit, func(PublicAnalyticsFailure) { reportCount++ })

	view, actualIdentity, err := handler.HandleWithVisit(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, "Public CV", view.Title)
	require.Equal(t, identity, actualIdentity)
	require.Zero(t, reportCount)
	require.Equal(t, 1, visit.calls)
}

func TestPublicResumeLinkDoesNotReportResumeReadFailureAsAnalyticsFailure(t *testing.T) {
	resumeErr := errors.New("resume storage unavailable")
	visit := &publicVisitCreatorStub{}
	var reported []PublicAnalyticsFailure
	handler, query, resume := newPublicAnalyticsTestHandler(nil, visit, func(failure PublicAnalyticsFailure) {
		reported = append(reported, failure)
	})
	resume.err = resumeErr

	view, identity, err := handler.HandleWithVisit(context.Background(), query)

	require.ErrorIs(t, err, resumeErr)
	require.Empty(t, view.Title)
	require.True(t, identity.ViewID.IsZero())
	require.Empty(t, reported)
	require.Zero(t, visit.calls)
}
