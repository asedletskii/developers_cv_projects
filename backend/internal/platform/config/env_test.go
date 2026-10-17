package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEnvironmentHelpersUseFallbackForUnsetOrBlankValues(t *testing.T) {
	t.Setenv("CONFIG_TEST_STRING", " \t")
	t.Setenv("CONFIG_TEST_INT", "")
	t.Setenv("CONFIG_TEST_BOOL", " ")
	t.Setenv("CONFIG_TEST_FLOAT", "")
	t.Setenv("CONFIG_TEST_DURATION", "")
	t.Setenv("CONFIG_TEST_INT64", "")

	require.Equal(t, "fallback", envString("CONFIG_TEST_STRING", "fallback"))
	integer, err := envInt("CONFIG_TEST_INT", 17)
	require.NoError(t, err)
	require.Equal(t, 17, integer)
	integer64, err := envInt64("CONFIG_TEST_INT64", 1<<40)
	require.NoError(t, err)
	require.Equal(t, int64(1<<40), integer64)
	boolean, err := envBool("CONFIG_TEST_BOOL", true)
	require.NoError(t, err)
	require.True(t, boolean)
	floating, err := envFloat("CONFIG_TEST_FLOAT", 1.25)
	require.NoError(t, err)
	require.Equal(t, 1.25, floating)
	duration, err := envDuration("CONFIG_TEST_DURATION", 3*time.Second)
	require.NoError(t, err)
	require.Equal(t, 3*time.Second, duration)
}

func TestEnvironmentHelpersRejectMalformedValuesWithVariableName(t *testing.T) {
	tests := []struct {
		name  string
		value string
		parse func(string) error
	}{
		{name: "integer", value: "1.5", parse: func(key string) error { _, err := envInt(key, 0); return err }},
		{name: "int64", value: "9223372036854775808", parse: func(key string) error { _, err := envInt64(key, 0); return err }},
		{name: "boolean", value: "yes", parse: func(key string) error { _, err := envBool(key, false); return err }},
		{name: "float", value: "not-a-number", parse: func(key string) error { _, err := envFloat(key, 0); return err }},
		{name: "duration", value: "tomorrow", parse: func(key string) error { _, err := envDuration(key, time.Second); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "CONFIG_TEST_" + tt.name
			t.Setenv(key, tt.value)
			err := tt.parse(key)
			require.Error(t, err)
			require.Contains(t, err.Error(), key)
		})
	}
}

func TestSplitCSVTrimsAndDropsEmptyEntries(t *testing.T) {
	require.Equal(t, []string{"first", "second", "third"}, splitCSV(" first, ,second,third ,, "))
	require.Empty(t, splitCSV(" , , "))
}
