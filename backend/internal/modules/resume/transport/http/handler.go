package http

import (
	"errors"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	createhandler "developers-cv/backend/internal/modules/resume/transport/http/create"
	gethandler "developers-cv/backend/internal/modules/resume/transport/http/get"
	listhandler "developers-cv/backend/internal/modules/resume/transport/http/list"
	patchhandler "developers-cv/backend/internal/modules/resume/transport/http/patch"
	updatehandler "developers-cv/backend/internal/modules/resume/transport/http/update"
)

type Handler struct {
	get    *gethandler.Handler
	list   *listhandler.Handler
	create *createhandler.Handler
	update *updatehandler.Handler
	patch  *patchhandler.Handler
}

func NewHandler(
	getResume *queries.GetResumeHandler,
	listResumes *queries.ListResumesHandler,
	createResume *commands.CreateResumeHandler,
	updateResume *commands.UpdateResumeHandler,
	patchResume *commands.PatchResumeHandler,
) *Handler {
	return &Handler{
		get:    gethandler.NewHandler(getResume),
		list:   listhandler.NewHandler(listResumes),
		create: createhandler.NewHandler(createResume, getResume),
		update: updatehandler.NewHandler(updateResume, getResume),
		patch:  patchhandler.NewHandler(patchResume),
	}
}

func (h *Handler) GetResume(c *gin.Context) {
	if h == nil || h.get == nil {
		common.WriteError(c, errors.New("get resume handler dependencies are nil"))
		return
	}
	h.get.Private(c)
}

func (h *Handler) GetPublicResume(c *gin.Context) {
	if h == nil || h.get == nil {
		common.WriteError(c, errors.New("get resume handler dependencies are nil"))
		return
	}
	h.get.Public(c)
}

func (h *Handler) ListResumes(c *gin.Context) {
	if h == nil || h.list == nil {
		common.WriteError(c, errors.New("list resumes handler dependencies are nil"))
		return
	}
	h.list.Handle(c)
}

func (h *Handler) CreateResume(c *gin.Context) {
	if h == nil || h.create == nil {
		common.WriteError(c, errors.New("create resume handler dependencies are nil"))
		return
	}
	h.create.Handle(c)
}

func (h *Handler) UpdateResume(c *gin.Context) {
	if h == nil || h.update == nil {
		common.WriteError(c, errors.New("update resume handler dependencies are nil"))
		return
	}
	h.update.Handle(c)
}

func (h *Handler) PatchResume(c *gin.Context) {
	if h == nil || h.patch == nil {
		common.WriteError(c, errors.New("patch resume handler dependencies are nil"))
		return
	}
	h.patch.Handle(c)
}
