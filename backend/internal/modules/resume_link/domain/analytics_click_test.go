package domain

import (
	"testing"
	"time"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"

	"github.com/stretchr/testify/require"
)

func TestAnalyticsClickValidation(t *testing.T) {
	now := time.Now().UTC()
	click, err := NewAnalyticsClick(AnalyticsClickParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Section: resumedomain.SectionAbout, XPct: 4210, YPct: 1733, ViewportW: 1440,
		LayoutHash: "0123456789abcdef", Target: "download", OccurredAt: now,
	})
	require.NoError(t, err)
	require.False(t, click.ID().IsZero())
	require.Equal(t, 4210, click.XPct())
	require.Equal(t, "0123456789abcdef", click.LayoutHash())

	_, err = NewAnalyticsClick(AnalyticsClickParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Section: resumedomain.SectionAbout, XPct: 10001, YPct: 0, ViewportW: 1440,
		LayoutHash: "0123456789abcdef", OccurredAt: now,
	})
	require.ErrorIs(t, err, ErrInvalidAnalyticsClick)
}

func TestAnalyticsClickCopiesOccurredAtToUTC(t *testing.T) {
	localTime := time.Date(2026, 10, 5, 10, 0, 0, 0, time.FixedZone("test", 2*60*60))
	click, err := NewAnalyticsClick(AnalyticsClickParams{
		ViewID: sharedvalueobjects.NewID(), LinkID: sharedvalueobjects.NewID(),
		Section: resumedomain.SectionAbout, XPct: 0, YPct: 0, ViewportW: 1,
		LayoutHash: "0123456789abcdef", OccurredAt: localTime,
	})
	require.NoError(t, err)
	require.Equal(t, localTime.UTC(), click.OccurredAt())
}
