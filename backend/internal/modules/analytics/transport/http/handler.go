package http

import (
	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/analytics/application/commands"
	"developers-cv/backend/internal/modules/analytics/application/queries"
	"developers-cv/backend/internal/platform/http/middleware"
)

type Handler struct {
	events    *EventsHandler
	read      *ReadHandler
	hasClicks bool
}

type HandlerDependencies struct {
	Ingest           *commands.IngestAnalyticsBatchHandler
	Statistics       *queries.LinkStatisticsHandler
	Daily            *queries.DailyStatisticsHandler
	Visits           *queries.VisitsHandler
	Sections         *queries.SectionsHandler
	Technologies     *queries.TechnologiesHandler
	ResumeStatistics *queries.ResumeStatisticsHandler
	Clicks           *queries.AnalyticsClickQueryHandler
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		events: NewEventsHandler(deps.Ingest),
		read: NewReadHandler(ReadHandlerDependencies{
			Statistics:       deps.Statistics,
			Daily:            deps.Daily,
			Visits:           deps.Visits,
			Sections:         deps.Sections,
			Technologies:     deps.Technologies,
			ResumeStatistics: deps.ResumeStatistics,
			Clicks:           deps.Clicks,
		}),
		hasClicks: deps.Clicks != nil,
	}
}

func RegisterRoutes(handler *Handler, anonymous, authenticated *gin.RouterGroup, eventsRateLimiter gin.HandlerFunc) {
	if handler == nil {
		return
	}
	if handler.events != nil && handler.events.ingest != nil {
		routeHandlers := []gin.HandlerFunc{middleware.BodyLimit(AnalyticsBatchMaxBodyBytes)}
		if eventsRateLimiter != nil {
			routeHandlers = append(routeHandlers, eventsRateLimiter)
		}
		routeHandlers = append(routeHandlers, handler.events.Handle)
		anonymous.POST("/cv/:linkId/events", routeHandlers...)
	}
	if handler.read == nil {
		return
	}
	authenticated.GET("/link/:linkId/statistics", handler.read.Statistics)
	authenticated.GET("/link/:linkId/statistics/daily", handler.read.DailyStatistics)
	authenticated.GET("/link/:linkId/views", handler.read.Visits)
	authenticated.GET("/link/:linkId/statistics/sections", handler.read.SectionStatistics)
	authenticated.GET("/link/:linkId/statistics/technologies", handler.read.TechnologyStatistics)
	if handler.hasClicks {
		authenticated.GET("/link/:linkId/clicks", handler.read.Clicks)
	}
	authenticated.GET("/resume/:resumeId/statistics", handler.read.ResumeStatistics)
}
