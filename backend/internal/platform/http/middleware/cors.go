package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/platform/config"
	"developers-cv/backend/internal/platform/http/requestid"
)

func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	allowHeaders := []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", requestid.Header}
	if cfg.AllowMockAuthHeader {
		allowHeaders = append(allowHeaders, "X-Mock-User-ID")
	}
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE"},
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    []string{requestid.Header, "X-View-Id"},
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           12 * 60 * 60,
	})
}
