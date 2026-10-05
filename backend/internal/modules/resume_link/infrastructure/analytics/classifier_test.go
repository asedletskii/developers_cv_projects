package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestTrafficClassifierOrderAndDistinctLinkBurst(t *testing.T) {
	classifier := NewTrafficClassifier(5, time.Minute)
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	ip := make([]byte, 32)

	owner, err := classifier.Classify(context.Background(), ports.TrafficClassificationInput{OwnerID: sharedvalueobjects.NewID(), UserAgent: "Googlebot", IPHash: ip, LinkID: sharedvalueobjects.NewID(), Now: now})
	require.NoError(t, err)
	require.Equal(t, domain.TrafficOwner, owner.Traffic)

	bot, err := classifier.Classify(context.Background(), ports.TrafficClassificationInput{UserAgent: "Googlebot", IPHash: ip, LinkID: sharedvalueobjects.NewID(), Now: now})
	require.NoError(t, err)
	require.Equal(t, domain.TrafficBot, bot.Traffic)

	for i := 0; i < 5; i++ {
		result, err := classifier.Classify(context.Background(), ports.TrafficClassificationInput{UserAgent: "Mozilla/5.0", IPHash: ip, LinkID: sharedvalueobjects.NewID(), Now: now.Add(time.Duration(i) * time.Second)})
		require.NoError(t, err)
		require.Equal(t, domain.TrafficHuman, result.Traffic)
	}

	result, err := classifier.Classify(context.Background(), ports.TrafficClassificationInput{UserAgent: "Mozilla/5.0", IPHash: ip, LinkID: sharedvalueobjects.NewID(), Now: now.Add(5 * time.Second)})
	require.NoError(t, err)
	require.Equal(t, domain.TrafficSuspect, result.Traffic)
	require.Equal(t, "burst", result.Reason)
}

func TestTrafficClassifierRepeatedSameLinkDoesNotTriggerBurst(t *testing.T) {
	classifier := NewTrafficClassifier(5, time.Minute)
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	ip := make([]byte, 32)
	linkID := sharedvalueobjects.NewID()

	for i := 0; i < 20; i++ {
		result, err := classifier.Classify(context.Background(), ports.TrafficClassificationInput{UserAgent: "Mozilla/5.0", IPHash: ip, LinkID: linkID, Now: now.Add(time.Duration(i) * time.Second)})
		require.NoError(t, err)
		require.Equal(t, domain.TrafficHuman, result.Traffic)
	}
}

func TestTrafficClassifierEmptyUserAgentIsBot(t *testing.T) {
	classifier := NewTrafficClassifier(5, time.Minute)
	result, err := classifier.Classify(context.Background(), ports.TrafficClassificationInput{Now: time.Now().UTC(), LinkID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32)})
	require.NoError(t, err)
	require.Equal(t, domain.TrafficBot, result.Traffic)
}
