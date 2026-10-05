package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

func parseAnalyticsPeriod(c *gin.Context) (queries.AnalyticsPeriod, bool) {
	period, err := queries.NewAnalyticsPeriod(c.Query("from"), c.Query("to"), c.Query("tz"), time.Now().UTC())
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid analytics period"))
		return queries.AnalyticsPeriod{}, false
	}
	return period, true
}

func writeAnalyticsReadError(c *gin.Context, err error) {
	if errors.Is(err, queries.ErrInvalidAnalyticsPeriod) || errors.Is(err, queries.ErrInvalidAnalyticsLimit) || errors.Is(err, queries.ErrInvalidAnalyticsCursor) || errors.Is(err, domain.ErrInvalidAnalyticsClick) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid analytics query"))
		return
	}
	if errors.Is(err, resumequeries.ErrResumeNotFound) || errors.Is(err, queries.ErrResumeLinkNotFound) || errors.Is(err, domain.ErrEntityNotFound) {
		common.WriteResumeLinkNotFound(c)
		return
	}
	apperrors.WriteError(c, err)
}

func (h *Handler) Statistics(c *gin.Context) {
	if h == nil || h.analyticsRead == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics read handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	period, ok := parseAnalyticsPeriod(c)
	if !ok {
		return
	}
	result, err := h.analyticsRead.Statistics(c.Request.Context(), queries.LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	if err != nil {
		writeAnalyticsReadError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) DailyStatistics(c *gin.Context) {
	if h == nil || h.analyticsRead == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics read handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	period, ok := parseAnalyticsPeriod(c)
	if !ok {
		return
	}
	result, err := h.analyticsRead.Daily(c.Request.Context(), queries.LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	if err != nil {
		writeAnalyticsReadError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Visits(c *gin.Context) {
	if h == nil || h.analyticsRead == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics read handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	limit := 0
	if raw := c.Query("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid limit"))
			return
		}
		limit = v
	}
	includeFiltered := false
	if raw, exists := c.GetQuery("includeFiltered"); exists {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid includeFiltered"))
			return
		}
		includeFiltered = parsed
	}
	result, err := h.analyticsRead.Visits(c.Request.Context(), queries.LinkVisitsQuery{UserID: userID, LinkID: linkID, Cursor: c.Query("cursor"), Limit: limit, IncludeFiltered: includeFiltered})
	if err != nil {
		writeAnalyticsReadError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) SectionStatistics(c *gin.Context) {
	if h == nil || h.analyticsRead == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics read handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	period, ok := parseAnalyticsPeriod(c)
	if !ok {
		return
	}
	result, err := h.analyticsRead.Sections(c.Request.Context(), queries.LinkSectionsQuery{UserID: userID, LinkID: linkID, Period: period})
	if err != nil {
		writeAnalyticsReadError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) TechnologyStatistics(c *gin.Context) {
	if h == nil || h.analyticsRead == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics read handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	period, ok := parseAnalyticsPeriod(c)
	if !ok {
		return
	}
	result, err := h.analyticsRead.Technologies(c.Request.Context(), queries.LinkAnalyticsQuery{UserID: userID, LinkID: linkID, Period: period})
	if err != nil {
		writeAnalyticsReadError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Clicks(c *gin.Context) {
	if h == nil || h.analyticsClicks == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics click handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}
	section := strings.TrimSpace(c.Query("section"))
	breakpoint := strings.TrimSpace(c.Query("breakpoint"))
	if section == "" || breakpoint == "" {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid click-map query"))
		return
	}
	includeStale := false
	if raw, exists := c.GetQuery("includeStale"); exists {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid includeStale"))
			return
		}
		includeStale = parsed
	}
	result, err := h.analyticsClicks.Handle(c.Request.Context(), queries.LinkClicksQuery{UserID: userID, LinkID: linkID, Section: section, Breakpoint: breakpoint, IncludeStale: includeStale})
	if err != nil {
		writeAnalyticsReadError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ResumeStatistics(c *gin.Context) {
	if h == nil || h.analyticsRead == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics read handler is not configured")))
		return
	}
	userID, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	resumeID, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	period, ok := parseAnalyticsPeriod(c)
	if !ok {
		return
	}
	result, err := h.analyticsRead.ResumeStatistics(c.Request.Context(), queries.ResumeAnalyticsQuery{UserID: userID, ResumeID: resumeID, Period: period})
	if err != nil {
		writeAnalyticsReadError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
