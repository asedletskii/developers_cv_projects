package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"developers-cv/backend/internal/platform/config"
)

func New(cfg config.LogConfig) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	return newProductionConfig(cfg, level).Build()
}

func newProductionConfig(cfg config.LogConfig, level zapcore.Level) zap.Config {
	production := zap.NewProductionConfig()
	production.Level = zap.NewAtomicLevelAt(level)
	production.Encoding = cfg.Format
	production.OutputPaths = []string{"stdout"}
	production.ErrorOutputPaths = []string{"stderr"}
	production.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	production.Sampling = nil
	return production
}

func parseLevel(value string) (zapcore.Level, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(value)); err != nil {
		return level, fmt.Errorf("invalid log level %q: %w", value, err)
	}
	return level, nil
}
