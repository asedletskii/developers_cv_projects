package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	activitiesHandler "developers-cv/backend/internal/modules/resume/transport/http/activities"
	contactshandler "developers-cv/backend/internal/modules/resume/transport/http/contacts"
	coursesHandler "developers-cv/backend/internal/modules/resume/transport/http/courses"
	educationHandler "developers-cv/backend/internal/modules/resume/transport/http/educations"
	experiencesHandler "developers-cv/backend/internal/modules/resume/transport/http/experiences"
	languagesHandler "developers-cv/backend/internal/modules/resume/transport/http/languages"
	overviewhandler "developers-cv/backend/internal/modules/resume/transport/http/overview"
	profilehandler "developers-cv/backend/internal/modules/resume/transport/http/profile"
	projectsHandler "developers-cv/backend/internal/modules/resume/transport/http/projects"
	publichandler "developers-cv/backend/internal/modules/resume/transport/http/public"
	skillsHandler "developers-cv/backend/internal/modules/resume/transport/http/skills"
	technologiesHandler "developers-cv/backend/internal/modules/resume/transport/http/technologies"
)

type Handler struct {
	resume       *resumeHandler
	overview     *overviewhandler.Handler
	profile      *profilehandler.Handler
	contacts     *contactshandler.Handler
	skills       *skillsHandler.Handler
	languages    *languagesHandler.Handler
	technologies *technologiesHandler.Handler
	experiences  *experiencesHandler.Handler
	education    *educationHandler.Handler
	courses      *coursesHandler.Handler
	activities   *activitiesHandler.Handler
	projects     *projectsHandler.Handler
	public       *publichandler.Handler
}

type HandlerDependencies struct {
	ListResumes    *queries.ListResumesQueryHandler
	CreateResume   *commands.CreateResumeHandler
	GetResume      *queries.GetResumeQueryHandler
	UpdateResume   *commands.UpdateResumeHandler
	DeleteResume   *commands.DeleteResumeHandler
	GetOverview    *queries.GetOverviewQueryHandler
	UpdateOverview *commands.UpdateOverviewHandler
	CreateProfile  *commands.CreateResumeProfileHandler
	GetProfile     *queries.GetResumeProfileQueryHandler
	UpdateProfile  *commands.UpdateResumeProfileHandler
	CreateContacts *commands.CreateResumeContactsHandler
	GetContacts    *queries.GetResumeContactsQueryHandler
	UpdateContacts *commands.UpdateResumeContactsHandler
	Skills         *skillsHandler.Handler
	Languages      *languagesHandler.Handler
	Technologies   *technologiesHandler.Handler
	Experiences    *experiencesHandler.Handler
	Education      *educationHandler.Handler
	Courses        *coursesHandler.Handler
	Activities     *activitiesHandler.Handler
	Projects       *projectsHandler.Handler
	Public         *publichandler.Handler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		resume:       newResumeHandler(deps.CreateResume, deps.GetResume, deps.UpdateResume, deps.DeleteResume, deps.ListResumes),
		overview:     overviewhandler.NewHandler(deps.GetOverview, deps.UpdateOverview),
		profile:      profilehandler.NewHandler(deps.CreateProfile, deps.GetProfile, deps.UpdateProfile),
		contacts:     contactshandler.NewHandler(deps.CreateContacts, deps.GetContacts, deps.UpdateContacts),
		skills:       deps.Skills,
		languages:    deps.Languages,
		technologies: deps.Technologies,
		experiences:  deps.Experiences,
		education:    deps.Education,
		courses:      deps.Courses,
		activities:   deps.Activities,
		projects:     deps.Projects,
		public:       deps.Public,
	}
}

func (h *Handler) ListResumes(c *gin.Context)          { h.resume.List(c) }
func (h *Handler) GetPublicResume(c *gin.Context)      { h.public.Handle(c) }
func (h *Handler) CreateResume(c *gin.Context)         { h.resume.Create(c) }
func (h *Handler) GetResume(c *gin.Context)            { h.resume.Get(c) }
func (h *Handler) UpdateResume(c *gin.Context)         { h.resume.Update(c) }
func (h *Handler) GetResumeOverview(c *gin.Context)    { h.overview.Get(c) }
func (h *Handler) UpdateResumeOverview(c *gin.Context) { h.overview.Update(c) }
func (h *Handler) DeleteResume(c *gin.Context)         { h.resume.Delete(c) }
func (h *Handler) CreateResumeProfile(c *gin.Context)  { h.profile.Create(c) }
func (h *Handler) GetResumeProfile(c *gin.Context)     { h.profile.Get(c) }
func (h *Handler) UpdateResumeProfile(c *gin.Context)  { h.profile.Update(c) }
func (h *Handler) CreateResumeContacts(c *gin.Context) { h.contacts.Create(c) }
func (h *Handler) GetResumeContacts(c *gin.Context)    { h.contacts.Get(c) }
func (h *Handler) UpdateResumeContacts(c *gin.Context) { h.contacts.Update(c) }
func (h *Handler) CreateSkills(c *gin.Context)         { h.skills.Create(c) }
func (h *Handler) GetSkills(c *gin.Context)            { h.skills.Get(c) }
func (h *Handler) UpdateSkills(c *gin.Context)         { h.skills.Update(c) }
func (h *Handler) CreateLanguage(c *gin.Context)       { h.languages.Create(c) }
func (h *Handler) ListLanguages(c *gin.Context)        { h.languages.List(c) }
func (h *Handler) GetLanguage(c *gin.Context)          { h.languages.Get(c) }
func (h *Handler) UpdateLanguage(c *gin.Context)       { h.languages.Update(c) }
func (h *Handler) DeleteLanguage(c *gin.Context)       { h.languages.Delete(c) }
func (h *Handler) CreateTechnology(c *gin.Context)     { h.technologies.Create(c) }
func (h *Handler) ListTechnologies(c *gin.Context)     { h.technologies.List(c) }
func (h *Handler) GetTechnology(c *gin.Context)        { h.technologies.Get(c) }
func (h *Handler) UpdateTechnology(c *gin.Context)     { h.technologies.Update(c) }
func (h *Handler) DeleteTechnology(c *gin.Context)     { h.technologies.Delete(c) }
func (h *Handler) CreateExperience(c *gin.Context)     { h.experiences.Create(c) }
func (h *Handler) ListExperiences(c *gin.Context)      { h.experiences.List(c) }
func (h *Handler) GetExperience(c *gin.Context)        { h.experiences.Get(c) }
func (h *Handler) UpdateExperience(c *gin.Context)     { h.experiences.Update(c) }
func (h *Handler) DeleteExperience(c *gin.Context)     { h.experiences.Delete(c) }
func (h *Handler) CreateEducation(c *gin.Context)      { h.education.Create(c) }
func (h *Handler) ListEducations(c *gin.Context)       { h.education.List(c) }
func (h *Handler) GetEducation(c *gin.Context)         { h.education.Get(c) }
func (h *Handler) UpdateEducation(c *gin.Context)      { h.education.Update(c) }
func (h *Handler) DeleteEducation(c *gin.Context)      { h.education.Delete(c) }
func (h *Handler) CreateCourse(c *gin.Context)         { h.courses.Create(c) }
func (h *Handler) ListCourses(c *gin.Context)          { h.courses.List(c) }
func (h *Handler) GetCourse(c *gin.Context)            { h.courses.Get(c) }
func (h *Handler) UpdateCourse(c *gin.Context)         { h.courses.Update(c) }
func (h *Handler) DeleteCourse(c *gin.Context)         { h.courses.Delete(c) }
func (h *Handler) CreateActivity(c *gin.Context)       { h.activities.Create(c) }
func (h *Handler) ListActivities(c *gin.Context)       { h.activities.List(c) }
func (h *Handler) GetActivity(c *gin.Context)          { h.activities.Get(c) }
func (h *Handler) UpdateActivity(c *gin.Context)       { h.activities.Update(c) }
func (h *Handler) DeleteActivity(c *gin.Context)       { h.activities.Delete(c) }
func (h *Handler) CreateProject(c *gin.Context)        { h.projects.Create(c) }
func (h *Handler) ListProjects(c *gin.Context)         { h.projects.List(c) }
func (h *Handler) GetProject(c *gin.Context)           { h.projects.Get(c) }
func (h *Handler) UpdateProject(c *gin.Context)        { h.projects.Update(c) }
func (h *Handler) DeleteProject(c *gin.Context)        { h.projects.Delete(c) }
