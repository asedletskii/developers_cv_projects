package http

import (
	"context"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/user/application/commands"
	"developers-cv/backend/internal/modules/user/application/queries"
	profilehandler "developers-cv/backend/internal/modules/user/transport/http/profile"
)

type Handler struct{ profile *profilehandler.Handler }

type profileReplacer interface {
	Handle(context.Context, commands.ReplaceProfileCommand) error
}

func NewHandler(getProfile *queries.GetProfileHandler, replaceProfile profileReplacer) *Handler {
	return &Handler{profile: profilehandler.NewHandler(getProfile, replaceProfile)}
}

func (h *Handler) GetProfile(c *gin.Context)     { h.profile.Get(c) }
func (h *Handler) ReplaceProfile(c *gin.Context) { h.profile.Replace(c) }
