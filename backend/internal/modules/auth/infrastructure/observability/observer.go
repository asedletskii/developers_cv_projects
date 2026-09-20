package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/platform/http/requestid"
	"developers-cv/backend/internal/platform/telemetry"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const ProviderTelegram = "telegram"

type Observer struct {
	logger  *zap.Logger
	metrics *telemetry.Metrics
	tracer  trace.Tracer
}

func NewObserver(logger *zap.Logger, metrics *telemetry.Metrics, serviceName string) *Observer {
	if serviceName == "" {
		serviceName = "developers-cv-api"
	}
	return &Observer{logger: logger, metrics: metrics, tracer: otel.Tracer(serviceName)}
}

func (o *Observer) LoginStarted(ctx context.Context, provider domain.Provider) {
	providerName := provider.String()
	if o.metrics != nil {
		o.metrics.RecordAuthLoginAttempt(providerName)
	}
	o.log(ctx, zap.InfoLevel, "auth.login.started", zap.String("provider", providerName))
}

func (o *Observer) LoginCallback(ctx context.Context, provider domain.Provider) {
	o.log(ctx, zap.InfoLevel, "auth.login.callback", zap.String("provider", provider.String()))
}

func (o *Observer) LoginSucceeded(ctx context.Context, provider domain.Provider, userID sharedvalueobjects.ID, duration time.Duration) {
	providerName := provider.String()
	if o.metrics != nil {
		o.metrics.RecordAuthLoginSuccess(providerName, duration)
	}
	o.log(ctx, zap.InfoLevel, "auth.login.succeeded",
		zap.String("provider", providerName),
		zap.String("user_id", userID.String()),
		zap.Int64("duration_ms", duration.Milliseconds()),
	)
}

func (o *Observer) LoginFailed(ctx context.Context, provider domain.Provider, reason string, duration time.Duration) {
	providerName := provider.String()
	if o.metrics != nil {
		o.metrics.RecordAuthLoginFailure(providerName, reason, duration)
	}
	o.log(ctx, zap.WarnLevel, "auth.login.failed",
		zap.String("provider", providerName),
		zap.String("reason", reason),
		zap.Int64("duration_ms", duration.Milliseconds()),
	)
}

func (o *Observer) Logout(ctx context.Context, userID sharedvalueobjects.ID, success bool) {
	reason := "success"
	if !success {
		reason = "failure"
	}
	if o.metrics != nil {
		o.metrics.RecordAuthLogout(ProviderTelegram, reason)
	}
	fields := []zap.Field{zap.String("provider", ProviderTelegram), zap.Bool("success", success)}
	if !userID.IsZero() {
		fields = append(fields, zap.String("user_id", userID.String()))
	}
	o.log(ctx, zap.InfoLevel, "auth.logout", fields...)
}

func (o *Observer) SessionRejected(ctx context.Context, reason string) {
	if o.metrics != nil {
		o.metrics.RecordAuthSessionRejection(ProviderTelegram, reason)
	}
	o.log(ctx, zap.WarnLevel, "auth.session.rejected",
		zap.String("provider", ProviderTelegram),
		zap.String("reason", reason),
	)
}

func (o *Observer) StartSpan(ctx context.Context, name string, provider domain.Provider) (context.Context, func(error)) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, span := o.tracer.Start(ctx, name)
	span.SetAttributes(attribute.String("auth.provider", provider.String()))
	return ctx, func(err error) {
		if err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.String("auth.result", "failure"))
		} else {
			span.SetAttributes(attribute.String("auth.result", "success"))
		}
		span.End()
	}
}

func (o *Observer) log(ctx context.Context, level zapcore.Level, message string, fields ...zap.Field) {
	if o == nil || o.logger == nil {
		return
	}
	fields = append(fields, zap.String("request_id", requestid.FromContext(ctx)))
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		fields = append(fields, zap.String("trace_id", spanContext.TraceID().String()))
	}
	switch level {
	case zap.InfoLevel:
		o.logger.Info(message, fields...)
	case zap.WarnLevel:
		o.logger.Warn(message, fields...)
	default:
		o.logger.Debug(message, fields...)
	}
}
