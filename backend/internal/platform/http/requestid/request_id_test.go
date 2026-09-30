package requestid

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetReturnsRequestIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(nil)
	c.Set(Key, "request-123")

	require.Equal(t, "request-123", Get(c))
}

func TestGetReturnsEmptyForNilContext(t *testing.T) {
	require.Empty(t, Get(nil))
}

func TestGetReturnsEmptyForMissingOrInvalidValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(nil)
	require.Empty(t, Get(c))

	c.Set(Key, 123)
	require.Empty(t, Get(c))
}

func TestWithContextRoundTripsRequestID(t *testing.T) {
	ctx := WithContext(context.Background(), "request-123")
	require.Equal(t, "request-123", FromContext(ctx))
}
