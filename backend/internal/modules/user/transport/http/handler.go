package http

import (
	"errors"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/transport/http/common"
	profilehandler "developers-cv/backend/internal/modules/user/transport/http/profile"
)

type Handler struct {
	profile *profilehandler.Handler
}

func NewHandler(getProfile *queries.GetProfileHandler) *Handler {
	return &Handler{profile: profilehandler.NewHandler(getProfile)}
}

func (h *Handler) GetProfile(c *gin.Context) {
	if h == nil || h.profile == nil {
		common.Internal(c, errors.New("get profile handler dependencies are nil"))
		return
	}
	h.profile.Get(c)
}
