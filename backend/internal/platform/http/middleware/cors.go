package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/platform/config"
)

func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Mock-User-ID", RequestIDHeader},
		ExposeHeaders:    []string{RequestIDHeader},
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           12 * 60 * 60,
	})
}
