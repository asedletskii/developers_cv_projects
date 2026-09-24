package create

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"developers-cv/backend/internal/modules/resume/transport/http/request"
)

type Handler struct {
	create *commands.CreateResumeHandler
	get    *queries.GetResumeHandler
}

func NewHandler(createHandler *commands.CreateResumeHandler, getHandler *queries.GetResumeHandler) *Handler {
	return &Handler{create: createHandler, get: getHandler}
}

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.create == nil || h.get == nil {
		common.WriteError(c, errors.New("create resume handler dependencies are nil"))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}

	var body request.ResumeRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteInvalidBody(c)
		return
	}
	command, err := body.ToCreateCommand(userID.UUID())
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	resumeID, err := h.create.Handle(c.Request.Context(), command)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}

	resume, err := h.get.Handle(c.Request.Context(), queries.GetResumeQuery{
		UserID: &userID,
		ID:     resumeID.UUID(),
		Lang:   queries.LanguageRU,
	})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resume)
}
