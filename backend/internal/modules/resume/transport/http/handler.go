package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	contactshandler "developers-cv/backend/internal/modules/resume/transport/http/contacts"
	createhandler "developers-cv/backend/internal/modules/resume/transport/http/create"
	experiencesHandler "developers-cv/backend/internal/modules/resume/transport/http/experiences"
	gethandler "developers-cv/backend/internal/modules/resume/transport/http/get"
	languagesHandler "developers-cv/backend/internal/modules/resume/transport/http/languages"
	listhandler "developers-cv/backend/internal/modules/resume/transport/http/list"
	patchhandler "developers-cv/backend/internal/modules/resume/transport/http/patch"
	profilehandler "developers-cv/backend/internal/modules/resume/transport/http/profile"
	resumeroot "developers-cv/backend/internal/modules/resume/transport/http/root"
	skillsHandler "developers-cv/backend/internal/modules/resume/transport/http/skills"
	technologiesHandler "developers-cv/backend/internal/modules/resume/transport/http/technologies"
	updatehandler "developers-cv/backend/internal/modules/resume/transport/http/update"
)

type Handler struct {
	create       *createhandler.Handler
	get          *gethandler.Handler
	list         *listhandler.Handler
	patch        *patchhandler.Handler
	update       *updatehandler.Handler
	root         *resumeroot.Handler
	profile      *profilehandler.Handler
	contacts     *contactshandler.Handler
	skills       *skillsHandler.Handler
	languages    *languagesHandler.Handler
	technologies *technologiesHandler.Handler
	experiences  *experiencesHandler.Handler
}

func NewHandler(
	getResume *queries.GetResumeHandler,
	listResumes *queries.ListResumesHandler,
	createResume *commands.CreateResumeHandler,
	updateResume *commands.UpdateResumeHandler,
	patchResume *commands.PatchResumeHandler,
	createRoot *commands.CreateResumeRootHandler,
	getRoot *commands.GetResumeRootHandler,
	updateRoot *commands.UpdateResumeRootHandler,
	deleteRoot *commands.DeleteResumeRootHandler,
	createProfile *commands.CreateResumeProfileHandler,
	getProfile *commands.GetResumeProfileHandler,
	updateProfile *commands.UpdateResumeProfileHandler,
	createContacts *commands.CreateResumeContactsHandler,
	getContacts *commands.GetResumeContactsHandler,
	updateContacts *commands.UpdateResumeContactsHandler,
	skills *skillsHandler.Handler,
	languages *languagesHandler.Handler,
	technologies *technologiesHandler.Handler,
	experiences *experiencesHandler.Handler,
) *Handler {
	return &Handler{
		create:       createhandler.NewHandler(createResume, getResume),
		get:          gethandler.NewHandler(getResume),
		list:         listhandler.NewHandler(listResumes),
		patch:        patchhandler.NewHandler(patchResume),
		update:       updatehandler.NewHandler(updateResume, getResume),
		root:         resumeroot.NewHandler(createRoot, getRoot, updateRoot, deleteRoot),
		profile:      profilehandler.NewHandler(createProfile, getProfile, updateProfile),
		contacts:     contactshandler.NewHandler(createContacts, getContacts, updateContacts),
		skills:       skills,
		languages:    languages,
		technologies: technologies,
		experiences:  experiences,
	}
}

func (h *Handler) CreateResume(c *gin.Context) {
	h.create.Handle(c)
}

func (h *Handler) GetResume(c *gin.Context) {
	h.get.Private(c)
}

func (h *Handler) GetPublicResume(c *gin.Context) {
	h.get.Public(c)
}

func (h *Handler) ListResumes(c *gin.Context) {
	h.list.Handle(c)
}

func (h *Handler) UpdateResume(c *gin.Context) {
	h.update.Handle(c)
}

func (h *Handler) PatchResume(c *gin.Context) {
	h.patch.Handle(c)
}

func (h *Handler) CreateResumeRoot(c *gin.Context) { h.root.Create(c) }
func (h *Handler) GetResumeRoot(c *gin.Context)    { h.root.Get(c) }
func (h *Handler) UpdateResumeRoot(c *gin.Context) { h.root.Update(c) }
func (h *Handler) DeleteResumeRoot(c *gin.Context) { h.root.Delete(c) }

func (h *Handler) CreateResumeProfile(c *gin.Context)  { h.profile.Create(c) }
func (h *Handler) GetResumeProfile(c *gin.Context)     { h.profile.Get(c) }
func (h *Handler) UpdateResumeProfile(c *gin.Context)  { h.profile.Update(c) }
func (h *Handler) CreateResumeContacts(c *gin.Context) { h.contacts.Create(c) }
func (h *Handler) GetResumeContacts(c *gin.Context)    { h.contacts.Get(c) }
func (h *Handler) UpdateResumeContacts(c *gin.Context) { h.contacts.Update(c) }

func (h *Handler) CreateSkills(c *gin.Context)     { h.skills.Create(c) }
func (h *Handler) GetSkills(c *gin.Context)        { h.skills.Get(c) }
func (h *Handler) UpdateSkills(c *gin.Context)     { h.skills.Update(c) }
func (h *Handler) CreateLanguage(c *gin.Context)   { h.languages.Create(c) }
func (h *Handler) ListLanguages(c *gin.Context)    { h.languages.List(c) }
func (h *Handler) GetLanguage(c *gin.Context)      { h.languages.Get(c) }
func (h *Handler) UpdateLanguage(c *gin.Context)   { h.languages.Update(c) }
func (h *Handler) DeleteLanguage(c *gin.Context)   { h.languages.Delete(c) }
func (h *Handler) CreateTechnology(c *gin.Context) { h.technologies.Create(c) }
func (h *Handler) ListTechnologies(c *gin.Context) { h.technologies.List(c) }
func (h *Handler) GetTechnology(c *gin.Context)    { h.technologies.Get(c) }
func (h *Handler) UpdateTechnology(c *gin.Context) { h.technologies.Update(c) }
func (h *Handler) DeleteTechnology(c *gin.Context) { h.technologies.Delete(c) }
func (h *Handler) CreateExperience(c *gin.Context) { h.experiences.Create(c) }
func (h *Handler) ListExperiences(c *gin.Context)  { h.experiences.List(c) }
func (h *Handler) GetExperience(c *gin.Context)    { h.experiences.Get(c) }
func (h *Handler) UpdateExperience(c *gin.Context) { h.experiences.Update(c) }
func (h *Handler) DeleteExperience(c *gin.Context) { h.experiences.Delete(c) }
