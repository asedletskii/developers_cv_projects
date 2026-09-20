package config

func parseAppConfig() AppConfig {
	return AppConfig{
		Env:     envString("APP_ENV", "development"),
		Version: envString("APP_VERSION", "dev"),
	}
}

func parseLogConfig() LogConfig {
	return LogConfig{Level: envString("LOG_LEVEL", "info"), Format: envString("LOG_FORMAT", "json")}
}
