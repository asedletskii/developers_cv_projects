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
	getUser *queries.GetUserHandler
}

func NewHandler(getUser *queries.GetUserHandler) *Handler {
	return &Handler{getUser: getUser}
}

func (h *Handler) GetProfile(c *gin.Context) {
	if h == nil || h.getUser == nil {
		writeError(c, apperrors.Internal(errors.New("get user handler dependencies are nil")))
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

	user, err := h.getUser.Handle(c.Request.Context(), queries.GetUserQuery{UserID: identity})
	if err != nil {
		if errors.Is(err, queries.ErrUserNotFound) {
			writeError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Profile not found"))
			return
		}
		writeError(c, apperrors.Internal(err))
		return
	}

	c.JSON(http.StatusOK, user)
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
