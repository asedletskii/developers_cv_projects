package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(handler *Handler, api *gin.RouterGroup, lk *gin.RouterGroup) {
	lk.GET("/resumes", handler.ListResumes)
	lk.GET("/resume/:resumeId", handler.GetResume)
	lk.POST("/resume", handler.CreateResume)
	lk.PUT("/resume/:resumeId", handler.UpdateResume)
	lk.PATCH("/resume/:resumeId", handler.PatchResume)
	api.GET("/cv/:resumeId", handler.GetPublicResume)
}
