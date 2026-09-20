package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/transport/http/common"
	"developers-cv/backend/internal/modules/user/transport/http/requests"
)

type profileReplacer interface {
	Handle(context.Context, commands.ReplaceProfileCommand) error
}

type ProfileHandler struct {
	query   *queries.GetProfileHandler
	replace profileReplacer
}

func NewProfileHandler(query *queries.GetProfileHandler, replace profileReplacer) *ProfileHandler {
	return &ProfileHandler{query: query, replace: replace}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	profile, err := h.query.Handle(c.Request.Context(), queries.GetProfileQuery{UserID: userID})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) ReplaceProfile(c *gin.Context) {
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}

	var body requests.ReplaceUserProfileRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}

	err := h.replace.Handle(c.Request.Context(), commands.ReplaceProfileCommand{
		UserID: userID,
		RU: commands.LocalizedProfileInput{
			Name: body.RU.Name, Surname: body.RU.Surname, Patronymic: body.RU.Patronymic,
			Position: body.RU.Position, Location: body.RU.Location,
		},
		EN: commands.LocalizedProfileInput{
			Name: body.EN.Name, Surname: body.EN.Surname, Patronymic: body.EN.Patronymic,
			Position: body.EN.Position, Location: body.EN.Location,
		},
		PhotoURL: body.PhotoURL,
	})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
