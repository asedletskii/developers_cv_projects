package config

func parseTelemetryConfig(app AppConfig) (TelemetryConfig, error) {
	enabled, err := envBool("OTEL_ENABLED", false)
	if err != nil {
		return TelemetryConfig{}, err
	}
	insecure, err := envBool("OTEL_EXPORTER_OTLP_INSECURE", false)
	if err != nil {
		return TelemetryConfig{}, err
	}
	sampleRatio, err := envFloat("OTEL_TRACES_SAMPLE_RATIO", 0.1)
	if err != nil {
		return TelemetryConfig{}, err
	}
	return TelemetryConfig{
		Enabled:        enabled,
		ServiceName:    envString("OTEL_SERVICE_NAME", "developers-cv-api"),
		ServiceVersion: app.Version,
		Environment:    app.Env,
		OTLPEndpoint:   envString("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		Insecure:       insecure,
		SampleRatio:    sampleRatio,
	}, nil
}

func parsePrometheusConfig() (PrometheusConfig, error) {
	enabled, err := envBool("PROMETHEUS_ENABLED", true)
	if err != nil {
		return PrometheusConfig{}, err
	}
	return PrometheusConfig{Enabled: enabled}, nil
}
