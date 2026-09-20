package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/transport/http/handlers"
)

type Handler struct {
	profile *handlers.ProfileHandler
}

type HandlerDependencies struct {
	GetProfile     *queries.GetProfileHandler
	ReplaceProfile *commands.ReplaceProfileHandler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{profile: handlers.NewProfileHandler(deps.GetProfile, deps.ReplaceProfile)}
}

func (h *Handler) GetProfile(c *gin.Context)     { h.profile.GetProfile(c) }
func (h *Handler) ReplaceProfile(c *gin.Context) { h.profile.ReplaceProfile(c) }
