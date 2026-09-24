package http

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
	getProfile *queries.GetProfileHandler
}

func NewHandler(getProfile *queries.GetProfileHandler) *Handler {
	return &Handler{getProfile: getProfile}
}

func (h *Handler) GetProfile(c *gin.Context) {
	if h == nil || h.getProfile == nil {
		writeError(c, apperrors.Internal(errors.New("get profile handler dependencies are nil")))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		writeError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return
	}

	identity, err := valueobjects.IDFromUUID(userID)
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}

	profile, err := h.getProfile.Handle(c.Request.Context(), queries.GetProfileQuery{UserID: identity})
	if err != nil {
		if errors.Is(err, queries.ErrProfileNotFound) {
			writeError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Profile not found"))
			return
		}
		writeError(c, apperrors.Internal(err))
		return
	}

	c.JSON(http.StatusOK, profile)
}

func writeError(c *gin.Context, err error) {
	appErr := apperrors.From(err)
	requestID := middleware.GetRequestID(c)
	c.AbortWithStatusJSON(appErr.Status, gin.H{
		"error": gin.H{
			"code":      appErr.Code,
			"message":   appErr.Message,
			"requestId": requestID,
		},
	})
}
