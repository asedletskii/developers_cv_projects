package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(handler *Handler, lk *gin.RouterGroup) {
	lk.GET("/profile", handler.GetProfile)
	lk.PUT("/profile", handler.ReplaceProfile)
}
