package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestAnalyticsUseCaseHandlersDelegateToDedicatedReaders(t *testing.T) {
	ctx := context.Background()
	userID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{ID: linkID.String(), ResumeID: resumeID.String()}}
	ownership := &ownershipReaderStub{}
	reader := newAnalyticsUseCaseReaderStub()
	period := validAnalyticsTestPeriod(t)

	stats, err := NewLinkStatisticsHandler(links, ownership, reader).Handle(ctx, LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.NoError(t, err)
	require.Equal(t, int64(7), stats.Views)
	require.Contains(t, stats.Devices, "desktop")

	_, err = NewDailyStatisticsHandler(links, ownership, reader).Handle(ctx, LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.NoError(t, err)
	_, err = NewVisitsHandler(links, ownership, reader).Handle(ctx, LinkVisitsQuery{UserID: userID, LinkID: linkID})
	require.NoError(t, err)
	sections, err := NewSectionsHandler(links, ownership, reader).Handle(ctx, LinkSectionsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.NoError(t, err)
	require.Equal(t, "about", sections.Sections[0].Section)
	_, err = NewTechnologiesHandler(links, ownership, reader).Handle(ctx, LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.NoError(t, err)
	_, err = NewResumeStatisticsHandler(ownership, reader).Handle(ctx, ResumeAnalyticsQuery{UserID: userID, ResumeID: resumeID, Period: period})
	require.NoError(t, err)

	require.Equal(t, map[string]int{"statistics": 1, "daily": 1, "visits": 1, "sections": 1, "technologies": 1, "resume_statistics": 1}, reader.calls)
}

func TestAnalyticsUseCaseHandlersAuthorizeBeforeReading(t *testing.T) {
	ctx := context.Background()
	userID := sharedvalueobjects.NewID()
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{ID: linkID.String(), ResumeID: resumeID.String()}}
	ownership := &ownershipReaderStub{err: resumequeries.ErrResumeNotFound}
	reader := newAnalyticsUseCaseReaderStub()
	period := validAnalyticsTestPeriod(t)

	_, err := NewLinkStatisticsHandler(links, ownership, reader).Handle(ctx, LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	_, err = NewDailyStatisticsHandler(links, ownership, reader).Handle(ctx, LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	_, err = NewVisitsHandler(links, ownership, reader).Handle(ctx, LinkVisitsQuery{UserID: userID, LinkID: linkID})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	_, err = NewSectionsHandler(links, ownership, reader).Handle(ctx, LinkSectionsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	_, err = NewTechnologiesHandler(links, ownership, reader).Handle(ctx, LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)
	_, err = NewResumeStatisticsHandler(ownership, reader).Handle(ctx, ResumeAnalyticsQuery{UserID: userID, ResumeID: resumeID, Period: period})
	require.ErrorIs(t, err, resumequeries.ErrResumeNotFound)

	require.Empty(t, reader.calls, "analytics readers must not run for a resume the user does not own")
}

func TestAnalyticsUseCaseHandlersRejectInvalidPeriodsBeforeReading(t *testing.T) {
	reader := newAnalyticsUseCaseReaderStub()
	linkID, userID, resumeID := sharedvalueobjects.NewID(), sharedvalueobjects.NewID(), sharedvalueobjects.NewID()
	links := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(),
	}}
	ownership := &ownershipReaderStub{}

	_, statsErr := NewLinkStatisticsHandler(links, ownership, reader).Handle(context.Background(), LinkAnalyticsQuery{UserID: userID, LinkID: linkID})
	_, dailyErr := NewDailyStatisticsHandler(links, ownership, reader).Handle(context.Background(), LinkAnalyticsQuery{UserID: userID, LinkID: linkID})
	_, sectionsErr := NewSectionsHandler(links, ownership, reader).Handle(context.Background(), LinkSectionsQuery{UserID: userID, LinkID: linkID})
	_, technologiesErr := NewTechnologiesHandler(links, ownership, reader).Handle(context.Background(), LinkAnalyticsQuery{UserID: userID, LinkID: linkID})
	_, resumeErr := NewResumeStatisticsHandler(ownership, reader).Handle(context.Background(), ResumeAnalyticsQuery{UserID: userID, ResumeID: resumeID})

	require.ErrorIs(t, statsErr, ErrInvalidAnalyticsPeriod)
	require.ErrorIs(t, dailyErr, ErrInvalidAnalyticsPeriod)
	require.ErrorIs(t, sectionsErr, ErrInvalidAnalyticsPeriod)
	require.ErrorIs(t, technologiesErr, ErrInvalidAnalyticsPeriod)
	require.ErrorIs(t, resumeErr, ErrInvalidAnalyticsPeriod)
	require.Empty(t, reader.calls)
}

func TestLinkStatisticsHandlerRejectsMissingDependencies(t *testing.T) {
	_, err := NewLinkStatisticsHandler(nil, nil, nil).Handle(context.Background(), LinkAnalyticsQuery{})
	require.Error(t, err)
}
