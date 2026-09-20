package config

import "time"

func parseAnalyticsConfig() (AnalyticsConfig, error) {
	burst, err := envInt("ANALYTICS_SUSPECT_BURST", 20)
	if err != nil {
		return AnalyticsConfig{}, err
	}
	window, err := envDuration("ANALYTICS_SUSPECT_WINDOW", time.Minute)
	if err != nil {
		return AnalyticsConfig{}, err
	}
	clickRetention, err := envInt("ANALYTICS_CLICK_RETENTION_DAYS", 90)
	if err != nil {
		return AnalyticsConfig{}, err
	}
	viewRetention, err := envInt("ANALYTICS_VIEW_RETENTION_DAYS", 365)
	if err != nil {
		return AnalyticsConfig{}, err
	}

	return AnalyticsConfig{
		VisitorSecret:      envString("ANALYTICS_VISITOR_SECRET", "development-only-analytics-secret-change-me-32"),
		SuspectBurst:       burst,
		SuspectWindow:      window,
		ClickRetentionDays: clickRetention,
		ViewRetentionDays:  viewRetention,
	}, nil
}
