package profile

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
)

type Handler struct {
	query *queries.GetProfileHandler
}

func NewHandler(query *queries.GetProfileHandler) *Handler {
	return &Handler{query: query}
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
