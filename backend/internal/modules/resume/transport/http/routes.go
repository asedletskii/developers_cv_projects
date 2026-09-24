package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(handler *Handler, api *gin.RouterGroup, lk *gin.RouterGroup) {
	lk.GET("/resumes", handler.ListResumes)
	lk.POST("/resumes", handler.CreateResumeRoot)
	lk.GET("/resumes/:resumeId", handler.GetResumeRoot)
	lk.PUT("/resumes/:resumeId", handler.UpdateResumeRoot)
	lk.DELETE("/resumes/:resumeId", handler.DeleteResumeRoot)
	lk.POST("/resume/:resumeId/profile", handler.CreateResumeProfile)
	lk.GET("/resume/:resumeId/profile", handler.GetResumeProfile)
	lk.PUT("/resume/:resumeId/profile", handler.UpdateResumeProfile)
	lk.POST("/resume/:resumeId/contacts", handler.CreateResumeContacts)
	lk.GET("/resume/:resumeId/contacts", handler.GetResumeContacts)
	lk.PUT("/resume/:resumeId/contacts", handler.UpdateResumeContacts)
	lk.GET("/resume/:resumeId", handler.GetResume)
	lk.POST("/resume", handler.CreateResume)
	lk.PUT("/resume/:resumeId", handler.UpdateResume)
	lk.PATCH("/resume/:resumeId", handler.PatchResume)
	api.GET("/cv/:resumeId", handler.GetPublicResume)
}
