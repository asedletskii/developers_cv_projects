package requestid

import "github.com/gin-gonic/gin"

const (
	Header = "X-Request-ID"
	Key    = "request_id"
)

func Get(c *gin.Context) string {
	if c == nil {
		return ""
	}

	value, exists := c.Get(Key)
	if !exists {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
