package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/transport/http/handlers"
)

type Handler struct {
	management  *handlers.ResumeLinkManagementHandler
	public      *handlers.PublicResumeLinkHandler
	contactView *handlers.PublicContactViewHandler
}

type HandlerDependencies struct {
	Create      *commands.CreateResumeLinkHandler
	List        *queries.ListResumeLinksQueryHandler
	Delete      *commands.DeleteResumeLinkHandler
	Public      *handlers.PublicResumeLinkHandler
	ContactView *handlers.PublicContactViewHandler
	Lifecycle   *commands.ReplaceStatusResumeLinkHandler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		management:  handlers.NewResumeLinkManagementHandler(deps.Create, deps.List, deps.Delete, deps.Lifecycle),
		public:      deps.Public,
		contactView: deps.ContactView,
	}
}

func RegisterRoutes(
	handler *Handler,
	publicRoutes *gin.RouterGroup,
	anonymous *gin.RouterGroup,
	lk *gin.RouterGroup,
	contactViewRateLimiter gin.HandlerFunc,
) {
	registerPublicRoutes(handler, publicRoutes, anonymous, contactViewRateLimiter)
	registerManagementRoutes(handler, lk)
}

func registerPublicRoutes(handler *Handler, publicRoutes, anonymous *gin.RouterGroup, contactViewRateLimiter gin.HandlerFunc) {
	publicRoutes.GET("/cv/:linkId", handler.public.Handle)
	contactHandlers := make([]gin.HandlerFunc, 0, 2)
	if contactViewRateLimiter != nil {
		contactHandlers = append(contactHandlers, contactViewRateLimiter)
	}
	contactHandlers = append(contactHandlers, handler.contactView.Handle)
	anonymous.POST("/cv/:linkId/contact-views", contactHandlers...)
}

func registerManagementRoutes(handler *Handler, lk *gin.RouterGroup) {
	management := handler.management
	lk.GET("/links", management.List)
	lk.POST("/link", management.Create)
	lk.DELETE("/link/:linkId", management.Delete)
	lk.PUT("/link/:linkId/status", management.ReplaceStatus)
}
