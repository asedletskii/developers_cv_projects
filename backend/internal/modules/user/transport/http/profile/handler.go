package profile

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/modules/user/transport/http/common"
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
	if h == nil || h.query == nil {
		common.Internal(c, errors.New("get profile handler dependencies are nil"))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		common.Unauthorized(c)
		return
	}

	identity, err := valueobjects.IDFromUUID(userID)
	if err != nil {
		common.Internal(c, err)
		return
	}

	profile, err := h.query.Handle(c.Request.Context(), queries.GetProfileQuery{UserID: identity})
	if err != nil {
		if errors.Is(err, queries.ErrProfileNotFound) {
			common.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Profile not found"))
			return
		}
		common.Internal(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}
