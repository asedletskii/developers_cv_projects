package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/application/queries"
	profilehandler "developers-cv/backend/internal/modules/user/transport/http/profile"
)

type Handler struct {
	profile *profilehandler.Handler
}

type HandlerDependencies struct {
	GetProfile     *queries.GetProfileHandler
	ReplaceProfile *commands.ReplaceProfileHandler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{profile: profilehandler.NewHandler(deps.GetProfile, deps.ReplaceProfile)}
}

func (h *Handler) GetProfile(c *gin.Context)     { h.profile.Get(c) }
func (h *Handler) ReplaceProfile(c *gin.Context) { h.profile.Replace(c) }
