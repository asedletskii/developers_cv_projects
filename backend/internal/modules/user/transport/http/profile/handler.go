package profile

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
)

type profileReplacer interface {
	Handle(context.Context, commands.ReplaceProfileCommand) error
}

type Handler struct {
	query   *queries.GetProfileHandler
	replace profileReplacer
}

func NewHandler(query *queries.GetProfileHandler, replace profileReplacer) *Handler {
	return &Handler{query: query, replace: replace}
}

func (h *Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	identity, err := valueobjects.IDFromUUID(userID)
	if err != nil {
		apperrors.WriteError(c, err)
		return
	}
	profile, err := h.query.Handle(c.Request.Context(), queries.GetProfileQuery{UserID: identity})
	if err != nil {
		if errors.Is(err, queries.ErrProfileNotFound) {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Profile not found"))
			return
		}
		apperrors.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

type localizedProfileRequest struct {
	Name       *string `json:"name"`
	Surname    *string `json:"surname"`
	Patronymic *string `json:"patronymic"`
	Position   *string `json:"position"`
	Location   *string `json:"location"`
}

type replaceProfileRequest struct {
	RU        localizedProfileRequest `json:"ru"`
	EN        localizedProfileRequest `json:"en"`
	PhotoPath *string                 `json:"photoPath"`
}

func (h *Handler) Replace(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.Unauthorized(c)
		return
	}
	identity, err := valueobjects.IDFromUUID(userID)
	if err != nil {
		apperrors.WriteError(c, err)
		return
	}

	var body replaceProfileRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}

	err = h.replace.Handle(c.Request.Context(), commands.ReplaceProfileCommand{
		UserID:    identity,
		RU:        commands.LocalizedProfileInput{Name: body.RU.Name, Surname: body.RU.Surname, Patronymic: body.RU.Patronymic, Position: body.RU.Position, Location: body.RU.Location},
		EN:        commands.LocalizedProfileInput{Name: body.EN.Name, Surname: body.EN.Surname, Patronymic: body.EN.Patronymic, Position: body.EN.Position, Location: body.EN.Location},
		PhotoPath: body.PhotoPath,
	})
	if err != nil {
		apperrors.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
