package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyticsTypesValidateKnownValues(t *testing.T) {
	require.True(t, TrafficHuman.IsValid())
	require.True(t, DeviceDesktop.IsValid())
	require.True(t, BrowserChrome.IsValid())
	require.True(t, OSLinux.IsValid())
	require.False(t, TrafficType("unknown").IsValid())
	require.True(t, DeviceType("unknown").IsValid())
	require.False(t, BrowserType("unknown").IsValid())
	require.False(t, OSType("unknown").IsValid())
}
