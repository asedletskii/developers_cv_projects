package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/application/queries"
	profilehandler "developers-cv/backend/internal/modules/user/transport/http/profile"
)

type Handler struct{ profile *profilehandler.Handler }

func NewHandler(getProfile *queries.GetProfileHandler, replaceProfile *commands.ReplaceProfileHandler) *Handler {
	return &Handler{profile: profilehandler.NewHandler(getProfile, replaceProfile)}
}

func (h *Handler) GetProfile(c *gin.Context)     { h.profile.Get(c) }
func (h *Handler) ReplaceProfile(c *gin.Context) { h.profile.Replace(c) }
