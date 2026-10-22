package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"developers-cv/backend/internal/platform/config"
)

func TestNewProductionConfigDisablesSampling(t *testing.T) {
	cfg := newProductionConfig(config.LogConfig{Level: "info", Format: "json"}, zapcore.InfoLevel)

	require.Nil(t, cfg.Sampling)
}
