package telemetry

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"

	"developers-cv/backend/internal/platform/config"
)

const exporterInitTimeout = 5 * time.Second

type Service struct {
	provider *sdktrace.TracerProvider
}

func New(cfg config.TelemetryConfig, log *zap.Logger) (*Service, error) {
	if !cfg.Enabled {
		return &Service{}, nil
	}

	options := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint)}
	if cfg.Insecure {
		options = append(options, otlptracegrpc.WithInsecure())
	} else {
		options = append(options,
			otlptracegrpc.WithTLSCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), exporterInitTimeout)
	defer cancel()

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP exporter: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentName(cfg.Environment),
		),
	)
	if err != nil {
		if shutdownErr := exporter.Shutdown(context.Background()); shutdownErr != nil {
			return nil, fmt.Errorf("create telemetry resource: %w; shutdown exporter: %v", err, shutdownErr)
		}
		return nil, fmt.Errorf("create telemetry resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRatio))),
	)
	otel.SetTracerProvider(provider)
	SetGlobalPropagation()
	log.Info(
		"OpenTelemetry tracing enabled",
		zap.String("endpoint", cfg.OTLPEndpoint),
		zap.String("service_name", cfg.ServiceName),
		zap.String("service_version", cfg.ServiceVersion),
		zap.String("environment", cfg.Environment),
	)
	return &Service{provider: provider}, nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	if s == nil || s.provider == nil {
		return nil
	}
	return s.provider.Shutdown(ctx)
}
