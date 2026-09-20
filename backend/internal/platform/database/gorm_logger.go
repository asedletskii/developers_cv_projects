package database

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func otelTracer() trace.Tracer { return otel.Tracer("developers-cv/database") }

type zapGORMLogger struct {
	log            *zap.Logger
	level          gormlogger.LogLevel
	slow           time.Duration
	ignoreNotFound bool
}

func newGORMLogger(log *zap.Logger, slow time.Duration) gormlogger.Interface {
	if log == nil {
		log = zap.NewNop()
	}
	if slow <= 0 {
		slow = 500 * time.Millisecond
	}
	return &zapGORMLogger{log: log, level: gormlogger.Warn, slow: slow, ignoreNotFound: true}
}

func (l *zapGORMLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *l
	clone.level = level
	return &clone
}

func (l *zapGORMLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Info {
		l.log.Info("gorm", zap.String("message", msg))
	}
}

func (l *zapGORMLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Warn {
		l.log.Warn("gorm", zap.String("message", msg))
	}
}

func (l *zapGORMLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Error {
		l.log.Error("gorm", zap.String("message", msg))
	}
}

func (l *zapGORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	duration := time.Since(begin)
	_, span := otelTracer().Start(ctx, "postgres.query", trace.WithTimestamp(begin))
	defer span.End(trace.WithTimestamp(time.Now()))
	span.SetAttributes(attribute.String("db.system.name", "postgresql"), attribute.Float64("db.client.operation.duration", duration.Seconds()))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "database operation failed")
	}

	if l.level == gormlogger.Silent {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) && l.ignoreNotFound {
		return
	}
	if err != nil && l.level >= gormlogger.Error {
		l.log.Error("database operation failed", zap.Duration("duration", duration), zap.Error(err))
		return
	}
	if duration >= l.slow && l.level >= gormlogger.Warn {
		l.log.Warn("slow database operation", zap.Duration("duration", duration), zap.Duration("threshold", l.slow))
	}
}
