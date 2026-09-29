package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(handler *Handler, api *gin.RouterGroup, lk *gin.RouterGroup) {
	lk.GET("/resumes", handler.ListResumes)
	lk.POST("/resume", handler.CreateResume)
	lk.GET("/resume/:resumeId", handler.GetResume)
	lk.PUT("/resume/:resumeId", handler.ReplaceResumeTitle)
	lk.PUT("/resume/:resumeId/status", handler.ChangeResumeStatus)
	lk.PUT("/resume/:resumeId/overview", handler.ReplaceResumeOverview)
	lk.DELETE("/resume/:resumeId", handler.DeleteResume)
	lk.PUT("/resume/:resumeId/profile", handler.ReplaceResumeProfile)
	lk.PUT("/resume/:resumeId/contacts", handler.ReplaceResumeContacts)
	lk.PUT("/resume/:resumeId/skills", handler.ReplaceSkills)
	lk.PUT("/resume/:resumeId/languages", handler.ReplaceLanguages)
	lk.PUT("/resume/:resumeId/technologies", handler.ReplaceTechnologies)
	lk.PUT("/resume/:resumeId/experiences", handler.ReplaceExperiences)
	lk.PUT("/resume/:resumeId/educations", handler.ReplaceEducations)
	lk.PUT("/resume/:resumeId/courses", handler.ReplaceCourses)
	lk.PUT("/resume/:resumeId/activities", handler.ReplaceActivities)
	lk.PUT("/resume/:resumeId/projects", handler.ReplaceProjects)
}
