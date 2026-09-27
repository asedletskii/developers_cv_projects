package overview

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	get    *queries.GetOverviewQueryHandler
	update *commands.UpdateOverviewHandler
}

func NewHandler(g *queries.GetOverviewQueryHandler, u *commands.UpdateOverviewHandler) *Handler {
	return &Handler{get: g, update: u}
}
func (h *Handler) Get(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	id, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	v, e := h.get.Handle(c.Request.Context(), queries.GetOverviewQuery{ResumeID: id, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, v)
}
func (h *Handler) Update(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	id, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var body overviewRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	cmd, e := body.ToCommand(uid, id)
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	if e = h.update.Handle(c.Request.Context(), cmd); e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
