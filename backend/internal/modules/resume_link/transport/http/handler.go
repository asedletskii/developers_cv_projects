package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/transport/http/handlers"
	"developers-cv/backend/internal/platform/http/middleware"
)

type Handler struct {
	management         *handlers.ResumeLinkManagementHandler
	public             *handlers.PublicResumeLinkHandler
	contactView        *handlers.PublicContactViewHandler
	analyticsEvents    *handlers.ResumeLinkAnalyticsEventsHandler
	analyticsRead      *handlers.ResumeLinkAnalyticsReadHandler
	hasAnalyticsClicks bool
}

type HandlerDependencies struct {
	Create                    *commands.CreateResumeLinkHandler
	List                      *queries.ListResumeLinksQueryHandler
	Delete                    *commands.DeleteResumeLinkHandler
	Public                    *handlers.PublicResumeLinkHandler
	ContactView               *handlers.PublicContactViewHandler
	AnalyticsEvents           *handlers.ResumeLinkAnalyticsEventsHandler
	AnalyticsStatistics       *queries.LinkStatisticsHandler
	AnalyticsDaily            *queries.DailyStatisticsHandler
	AnalyticsVisits           *queries.VisitsHandler
	AnalyticsSections         *queries.SectionsHandler
	AnalyticsTechnologies     *queries.TechnologiesHandler
	AnalyticsResumeStatistics *queries.ResumeStatisticsHandler
	AnalyticsClicks           *queries.AnalyticsClickQueryHandler
	Lifecycle                 *commands.ReplaceStatusResumeLinkHandler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		management:         handlers.NewResumeLinkManagementHandler(deps.Create, deps.List, deps.Delete, deps.Lifecycle),
		public:             deps.Public,
		contactView:        deps.ContactView,
		analyticsEvents:    deps.AnalyticsEvents,
		hasAnalyticsClicks: deps.AnalyticsClicks != nil,
		analyticsRead: handlers.NewResumeLinkAnalyticsReadHandler(handlers.ResumeLinkAnalyticsReadDependencies{
			Statistics:       deps.AnalyticsStatistics,
			Daily:            deps.AnalyticsDaily,
			Visits:           deps.AnalyticsVisits,
			Sections:         deps.AnalyticsSections,
			Technologies:     deps.AnalyticsTechnologies,
			ResumeStatistics: deps.AnalyticsResumeStatistics,
			Clicks:           deps.AnalyticsClicks,
		}),
	}
}

func RegisterRoutes(
	handler *Handler,
	publicRoutes *gin.RouterGroup,
	anonymous *gin.RouterGroup,
	lk *gin.RouterGroup,
	contactViewRateLimiter,
	analyticsEventsRateLimiter gin.HandlerFunc,
) {
	registerPublicRoutes(handler, publicRoutes, anonymous, contactViewRateLimiter)
	registerAnalyticsEventRoutes(handler, anonymous, analyticsEventsRateLimiter)
	registerManagementRoutes(handler, lk)
	registerAnalyticsReadRoutes(handler, lk)
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

func registerAnalyticsEventRoutes(handler *Handler, anonymous *gin.RouterGroup, rateLimiter gin.HandlerFunc) {
	if handler.analyticsEvents == nil {
		return
	}
	routeHandlers := []gin.HandlerFunc{middleware.BodyLimit(handlers.AnalyticsBatchMaxBodyBytes)}
	if rateLimiter != nil {
		routeHandlers = append(routeHandlers, rateLimiter)
	}
	routeHandlers = append(routeHandlers, handler.analyticsEvents.Handle)
	anonymous.POST("/cv/:linkId/events", routeHandlers...)
}

func registerManagementRoutes(handler *Handler, lk *gin.RouterGroup) {
	management := handler.management
	lk.GET("/links", management.List)
	lk.POST("/link", management.Create)
	lk.DELETE("/link/:linkId", management.Delete)
	lk.PUT("/link/:linkId/status", management.ReplaceStatus)
}

func registerAnalyticsReadRoutes(handler *Handler, lk *gin.RouterGroup) {
	analytics := handler.analyticsRead
	lk.GET("/link/:linkId/statistics", analytics.Statistics)
	lk.GET("/link/:linkId/statistics/daily", analytics.DailyStatistics)
	lk.GET("/link/:linkId/views", analytics.Visits)
	lk.GET("/link/:linkId/statistics/sections", analytics.SectionStatistics)
	lk.GET("/link/:linkId/statistics/technologies", analytics.TechnologyStatistics)
	if handler.hasAnalyticsClicks {
		lk.GET("/link/:linkId/clicks", analytics.Clicks)
	}
	lk.GET("/resume/:resumeId/statistics", analytics.ResumeStatistics)
}
