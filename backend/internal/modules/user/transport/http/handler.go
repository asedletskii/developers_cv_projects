package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/queries"
	profilehandler "developers-cv/backend/internal/modules/user/transport/http/profile"
)

type Handler struct {
	profile *profilehandler.Handler
}

func NewHandler(getProfile *queries.GetProfileHandler) *Handler {
	return &Handler{profile: profilehandler.NewHandler(getProfile)}
}

func (h *Handler) GetProfile(c *gin.Context) {
	h.profile.Get(c)
}
