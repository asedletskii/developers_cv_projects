package update

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
	update *commands.UpdateResumeHandler
	get    *queries.GetResumeHandler
}

func NewHandler(updateHandler *commands.UpdateResumeHandler, getHandler *queries.GetResumeHandler) *Handler {
	return &Handler{update: updateHandler, get: getHandler}
}

func (h *Handler) Handle(c *gin.Context) {
	if h == nil || h.update == nil || h.get == nil {
		common.WriteError(c, errors.New("update resume handler dependencies are nil"))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	resumeID, ok := common.ParseResumeID(c)
	if !ok {
		return
	}

	var body request.ResumeRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteInvalidBody(c)
		return
	}
	command, err := body.ToUpdateCommand(userID.UUID(), resumeID)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if _, err := h.update.Handle(c.Request.Context(), command); err != nil {
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
	c.JSON(http.StatusOK, resume)
}
