package telemetry

import (
	"context"
	"crypto/tls"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"

	"developers-cv/backend/internal/platform/config"
)

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
	exporter, err := otlptracegrpc.New(context.Background(), options...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP exporter: %w", err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceName(cfg.ServiceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("create telemetry resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
	otel.SetTracerProvider(provider)
	SetGlobalPropagation()
	log.Info("OpenTelemetry tracing enabled", zap.String("endpoint", cfg.OTLPEndpoint))
	return &Service{provider: provider}, nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	if s.provider == nil {
		return nil
	}
	return s.provider.Shutdown(ctx)
}
