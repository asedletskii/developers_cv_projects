package http

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(handler *Handler, auth *gin.RouterGroup, startRateLimiter, callbackRateLimiter gin.HandlerFunc) {
	start := auth.Group("")
	if startRateLimiter != nil {
		start.Use(startRateLimiter)
	}
	start.GET("/telegram", handler.StartTelegram)

	callback := auth.Group("")
	if callbackRateLimiter != nil {
		callback.Use(callbackRateLimiter)
	}
	callback.GET("/telegram/callback", handler.CompleteTelegram)
}

func RegisterProtectedRoutes(handler *Handler, auth *gin.RouterGroup) {
	auth.GET("/me", handler.CurrentSession)
	auth.POST("/logout", handler.Logout)
}
