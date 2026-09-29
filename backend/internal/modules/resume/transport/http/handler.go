package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	activitieshandler "developers-cv/backend/internal/modules/resume/transport/http/activities"
	contactshandler "developers-cv/backend/internal/modules/resume/transport/http/contacts"
	courseshandler "developers-cv/backend/internal/modules/resume/transport/http/courses"
	educationshandler "developers-cv/backend/internal/modules/resume/transport/http/educations"
	experienceshandler "developers-cv/backend/internal/modules/resume/transport/http/experiences"
	languageshandler "developers-cv/backend/internal/modules/resume/transport/http/languages"
	overviewhandler "developers-cv/backend/internal/modules/resume/transport/http/overview"
	profilehandler "developers-cv/backend/internal/modules/resume/transport/http/profile"
	projectshandler "developers-cv/backend/internal/modules/resume/transport/http/projects"
	skillshandler "developers-cv/backend/internal/modules/resume/transport/http/skills"
	technologieshandler "developers-cv/backend/internal/modules/resume/transport/http/technologies"
)

type Handler struct {
	resume       *resumeHandler
	overview     *overviewhandler.Handler
	profile      *profilehandler.Handler
	contacts     *contactshandler.Handler
	skills       *skillshandler.Handler
	languages    *languageshandler.Handler
	technologies *technologieshandler.Handler
	experiences  *experienceshandler.Handler
	educations   *educationshandler.Handler
	courses      *courseshandler.Handler
	activities   *activitieshandler.Handler
	projects     *projectshandler.Handler
}

type HandlerDependencies struct {
	ListResumes        *queries.ListResumesQueryHandler
	CreateResume       *commands.CreateResumeHandler
	GetResume          *queries.GetResumeQueryHandler
	ReplaceResumeTitle *commands.ReplaceResumeTitleHandler
	DeleteResume       *commands.DeleteResumeHandler
	ChangeResumeStatus *commands.ChangeResumeStatusHandler
	ReplaceOverview    *commands.ReplaceOverviewHandler
	ReplaceProfile     *commands.ReplaceProfileHandler
	ReplaceContacts    *commands.ReplaceContactsHandler
	Overview           *overviewhandler.Handler
	Profile            *profilehandler.Handler
	Contacts           *contactshandler.Handler
	Skills             *skillshandler.Handler
	Languages          *languageshandler.Handler
	Technologies       *technologieshandler.Handler
	Experiences        *experienceshandler.Handler
	Educations         *educationshandler.Handler
	Courses            *courseshandler.Handler
	Activities         *activitieshandler.Handler
	Projects           *projectshandler.Handler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		resume:       newResumeHandler(deps.CreateResume, deps.GetResume, deps.ReplaceResumeTitle, deps.DeleteResume, deps.ChangeResumeStatus, deps.ListResumes),
		overview:     deps.Overview,
		profile:      deps.Profile,
		contacts:     deps.Contacts,
		skills:       deps.Skills,
		languages:    deps.Languages,
		technologies: deps.Technologies,
		experiences:  deps.Experiences,
		educations:   deps.Educations,
		courses:      deps.Courses,
		activities:   deps.Activities,
		projects:     deps.Projects,
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
