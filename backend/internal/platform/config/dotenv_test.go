package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDotEnvDoesNotOverrideEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("TEST_DOTENV_VALUE", "from-environment")

	err := loadDotEnvContent([]byte(
		"TEST_DOTENV_VALUE=from-file\n" +
			"TEST_DOTENV_QUOTED_CONFIG_2_9_TEST=\"hello world\"\n",
	))

	require.NoError(t, err)
	require.Equal(t, "from-environment", os.Getenv("TEST_DOTENV_VALUE"))
	require.Equal(t, "hello world", os.Getenv("TEST_DOTENV_QUOTED_CONFIG_2_9_TEST"))
}

func TestLoadDotEnvDoesNotLoadProductionContentWithoutProcessEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("POSTGRES_PASSWORD", "")

	err := loadDotEnvContent([]byte("APP_ENV=production\nPOSTGRES_PASSWORD=from-file\n"))

	require.Error(t, err)
	require.Contains(t, err.Error(), "must be set in the process environment")
	require.Empty(t, os.Getenv("POSTGRES_PASSWORD"))
}

func TestLoadDotEnvSkipsContentInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("TEST_DOTENV_PRODUCTION_ONLY", "")

	require.NoError(t, loadDotEnvContent([]byte("TEST_DOTENV_PRODUCTION_ONLY=from-file\n")))
	require.Empty(t, os.Getenv("TEST_DOTENV_PRODUCTION_ONLY"))
}

func TestDotEnvParserSupportsCommentsAndEscapes(t *testing.T) {
	values := map[string]string{}
	input := "# whole line comment\nexport TEST_DOTENV_INLINE=plain-value # comment\nTEST_DOTENV_HASH=abc#123\nTEST_DOTENV_DOUBLE=\"hello\\nworld\" # comment\nTEST_DOTENV_SINGLE='literal\\ntext'\n"
	for _, line := range strings.Split(input, "\n") {
		key, value, ok, err := parseDotEnvLine(line)
		require.NoError(t, err)
		if ok {
			values[key] = value
		}
	}
	require.Equal(t, "plain-value", values["TEST_DOTENV_INLINE"])
	require.Equal(t, "abc#123", values["TEST_DOTENV_HASH"])
	require.Equal(t, "hello\nworld", values["TEST_DOTENV_DOUBLE"])
	require.Equal(t, "literal\\ntext", values["TEST_DOTENV_SINGLE"])
}

func TestDotEnvParserRejectsMalformedAssignments(t *testing.T) {
	_, _, _, err := parseDotEnvLine("NOT_AN_ASSIGNMENT")
	require.Error(t, err)

	_, _, _, err = parseDotEnvLine("1INVALID=value")
	require.Error(t, err)

	_, _, _, err = parseDotEnvLine("VALUE=\"unterminated")
	require.Error(t, err)
}
