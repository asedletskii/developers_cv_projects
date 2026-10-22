package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestViewMappingPreservesStateAndExistingTable(t *testing.T) {
	now := time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)
	layout, err := domain.NewLayoutSnapshot(
		[]domain.AnalyticsSectionID{domain.AnalyticsSectionAbout, domain.AnalyticsSectionExperience},
		map[domain.AnalyticsSectionID]string{
			domain.AnalyticsSectionAbout:      "0123456789abcdef",
			domain.AnalyticsSectionExperience: "fedcba9876543210",
		},
	)
	require.NoError(t, err)

	view, err := domain.RestoreView(domain.ViewSnapshot{
		ID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(), ResumeID: sharedvalueobjects.NewID(), VisitorID: sharedvalueobjects.NewID(),
		IPHash: make([]byte, 32), Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSLinux,
		Traffic: domain.TrafficHuman, Layout: layout, StartedAt: now, CreatedAt: now,
		State: domain.ViewStateSnapshot{
			LastSeq: 7, ConfirmedAt: now.Add(time.Second), LastSeenAt: now.Add(2 * time.Second),
			ActiveMS: 8000, MaxSection: domain.AnalyticsSectionExperience,
			ExitSection: domain.AnalyticsSectionAbout, EndedAt: now.Add(3 * time.Second),
		},
	})
	require.NoError(t, err)

	model, err := ViewToDB(view)
	require.NoError(t, err)
	require.Equal(t, "analytics.analytics_views", model.TableName())

	restored, err := ViewToDomain(model)
	require.NoError(t, err)
	require.Equal(t, view.ID(), restored.ID())
	require.Equal(t, view.LinkID(), restored.LinkID())
	require.Equal(t, view.ResumeID(), restored.ResumeID())
	require.Equal(t, view.State().Snapshot(), restored.State().Snapshot())
}
