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

func TestTrafficClassifierDoesNotTreatBotSubstringAsCrawler(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 12; Cubot X20 Pro) AppleWebKit/537.36 Chrome/100.0.0.0 Mobile Safari/537.36"

	require.False(t, isBotUA(ua))

	classifier := NewTrafficClassifier(5, time.Minute)
	result, err := classifier.Classify(context.Background(), ports.TrafficClassificationInput{
		UserAgent: ua,
		IPHash:    make([]byte, 32),
		LinkID:    sharedvalueobjects.NewID(),
		Now:       time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	require.Equal(t, domain.TrafficHuman, result.Traffic)
}

func TestTrafficClassifierRecognizesKnownCrawlerTokens(t *testing.T) {
	for _, ua := range []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Example Bot/1.0",
		"curl/8.0.0",
	} {
		t.Run(ua, func(t *testing.T) {
			require.True(t, isBotUA(ua))
		})
	}
}

func TestClassifyDeviceTreatsAndroidWithoutMobileTokenAsTablet(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; SM-X700) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	require.Equal(t, domain.DeviceTablet, classifyDevice(ua))
}

func TestClassifyDeviceAndOSRecognizeDesktopClassIPadOSUserAgent(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"

	require.Equal(t, domain.DeviceTablet, classifyDevice(ua))
	require.Equal(t, domain.OSIOS, classifyOS(ua))
}
