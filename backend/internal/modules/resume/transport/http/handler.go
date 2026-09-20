package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/handlers"
)

type Handler struct {
	resume       *handlers.ResumeHandler
	overview     *handlers.ResumeOverviewHandler
	profile      *handlers.ProfileHandler
	contacts     *handlers.ContactsHandler
	skills       *handlers.SkillsHandler
	languages    *handlers.ResumeSpokenLanguagesHandler
	technologies *handlers.ResumeTechnologiesHandler
	experiences  *handlers.ExperienceHandler
	educations   *handlers.EducationHandler
	courses      *handlers.CourseHandler
	activities   *handlers.ActivityHandler
	projects     *handlers.ProjectHandler
}

type HandlerDependencies struct {
	ListResumes         *queries.ListResumesQueryHandler
	CreateResume        *commands.CreateResumeHandler
	GetResume           *queries.GetResumeQueryHandler
	ReplaceResumeTitle  *commands.ReplaceResumeTitleHandler
	DeleteResume        *commands.DeleteResumeHandler
	ChangeResumeStatus  *commands.ChangeResumeStatusHandler
	ReplaceOverview     *commands.ReplaceOverviewHandler
	ReplaceProfile      *commands.ReplaceProfileHandler
	ReplaceContacts     *commands.ReplaceContactsHandler
	ReplaceSkills       *commands.ReplaceSkillsHandler
	ReplaceLanguages    *commands.ReplaceLanguagesHandler
	ReplaceTechnologies *commands.ReplaceTechnologiesHandler
	ReplaceExperiences  *commands.ReplaceExperiencesHandler
	ReplaceEducations   *commands.ReplaceEducationsHandler
	ReplaceCourses      *commands.ReplaceCoursesHandler
	ReplaceActivities   *commands.ReplaceActivitiesHandler
	ReplaceProjects     *commands.ReplaceProjectsHandler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		resume:       handlers.NewResumeHandler(deps.CreateResume, deps.GetResume, deps.ReplaceResumeTitle, deps.DeleteResume, deps.ChangeResumeStatus, deps.ListResumes),
		overview:     handlers.NewResumeOverviewHandler(deps.ReplaceOverview),
		profile:      handlers.NewProfileHandler(deps.ReplaceProfile),
		contacts:     handlers.NewContactsHandler(deps.ReplaceContacts),
		skills:       handlers.NewSkillsHandler(deps.ReplaceSkills),
		languages:    handlers.NewResumeSpokenLanguagesHandler(deps.ReplaceLanguages),
		technologies: handlers.NewResumeTechnologiesHandler(deps.ReplaceTechnologies),
		experiences:  handlers.NewExperienceHandler(deps.ReplaceExperiences),
		educations:   handlers.NewEducationHandler(deps.ReplaceEducations),
		courses:      handlers.NewCourseHandler(deps.ReplaceCourses),
		activities:   handlers.NewActivityHandler(deps.ReplaceActivities),
		projects:     handlers.NewProjectHandler(deps.ReplaceProjects),
	}
}

func (h *Handler) ListResumes(c *gin.Context)           { h.resume.List(c) }
func (h *Handler) CreateResume(c *gin.Context)          { h.resume.Create(c) }
func (h *Handler) GetResume(c *gin.Context)             { h.resume.Get(c) }
func (h *Handler) ReplaceResumeTitle(c *gin.Context)    { h.resume.ReplaceTitle(c) }
func (h *Handler) ChangeResumeStatus(c *gin.Context)    { h.resume.ChangeStatus(c) }
func (h *Handler) ReplaceResumeOverview(c *gin.Context) { h.overview.Replace(c) }
func (h *Handler) DeleteResume(c *gin.Context)          { h.resume.Delete(c) }
func (h *Handler) ReplaceResumeProfile(c *gin.Context)  { h.profile.Replace(c) }
func (h *Handler) ReplaceResumeContacts(c *gin.Context) { h.contacts.Replace(c) }
func (h *Handler) ReplaceSkills(c *gin.Context)         { h.skills.Replace(c) }
func (h *Handler) ReplaceLanguages(c *gin.Context)      { h.languages.Replace(c) }
func (h *Handler) ReplaceTechnologies(c *gin.Context)   { h.technologies.Replace(c) }
func (h *Handler) ReplaceExperiences(c *gin.Context)    { h.experiences.Replace(c) }
func (h *Handler) ReplaceEducations(c *gin.Context)     { h.educations.Replace(c) }
func (h *Handler) ReplaceCourses(c *gin.Context)        { h.courses.Replace(c) }
func (h *Handler) ReplaceActivities(c *gin.Context)     { h.activities.Replace(c) }
func (h *Handler) ReplaceProjects(c *gin.Context)       { h.projects.Replace(c) }
