package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const analyticsBatchMaxBodyBytes int64 = 32 * 1024

type AnalyticsEventsHandler struct {
	ingest *commands.IngestAnalyticsBatchHandler
}

func NewAnalyticsEventsHandler(ingest *commands.IngestAnalyticsBatchHandler) *AnalyticsEventsHandler {
	return &AnalyticsEventsHandler{ingest: ingest}
}

type analyticsBatchRequest struct {
	ViewID string             `json:"viewId"`
	Seq    *int64             `json:"seq"`
	Final  *bool              `json:"final"`
	State  *analyticsStateDTO `json:"state"`
	Events *[]json.RawMessage `json:"events"`
	Clicks *[]json.RawMessage `json:"clicks"`
}

type analyticsStateDTO struct {
	ActiveMS    int64  `json:"activeMs"`
	MaxSection  string `json:"maxSection"`
	ExitSection string `json:"exitSection"`
}

func (h *AnalyticsEventsHandler) Handle(c *gin.Context) {
	if h == nil || h.ingest == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics events handler is not configured")))
		return
	}

	linkID, ok := common.ParseLinkID(c)
	if !ok {
		return
	}

	var body analyticsBatchRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}

	viewUUID, err := uuid.Parse(strings.TrimSpace(body.ViewID))
	if err != nil || viewUUID == uuid.Nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid view id"))
		return
	}
	viewID, err := sharedvalueobjects.IDFromUUID(viewUUID)
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid view id"))
		return
	}
	if body.Seq == nil || body.Final == nil || body.State == nil || body.Events == nil || body.Clicks == nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid analytics batch"))
		return
	}
	if len(*body.Events) > commands.MaxAnalyticsEventsPerBatch || len(*body.Clicks) > commands.MaxAnalyticsClicksPerBatch {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Analytics batch limits exceeded"))
		return
	}

	events := make([]commands.AnalyticsEventInput, 0, len(*body.Events))
	for _, raw := range *body.Events {
		event, ok := decodeAnalyticsEvent(raw)
		if ok {
			events = append(events, event)
		}
	}
	clicks := make([]commands.AnalyticsClickInput, 0, len(*body.Clicks))
	for _, raw := range *body.Clicks {
		click, ok := decodeAnalyticsClick(raw)
		if ok {
			clicks = append(clicks, click)
		}
	}

	err = h.ingest.Handle(c.Request.Context(), commands.AnalyticsBatchInput{
		ViewID: viewID, LinkID: linkID, Seq: *body.Seq,
		ActiveMS: body.State.ActiveMS, MaxSection: body.State.MaxSection, ExitSection: body.State.ExitSection,
		Final: *body.Final, Events: events, Clicks: clicks,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResumeLinkInactive):
			apperrors.WriteError(c, apperrors.New(apperrors.CodeGone, http.StatusGone, "Resume link is no longer active"))
		case errors.Is(err, domain.ErrEntityNotFound), errors.Is(err, domain.ErrAnalyticsViewExpired):
			apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Analytics view not found"))
		case errors.Is(err, domain.ErrInvalidAnalyticsEvent), errors.Is(err, domain.ErrInvalidViewState):
			apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid analytics batch"))
		default:
			apperrors.WriteError(c, err)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func decodeAnalyticsEvent(raw json.RawMessage) (commands.AnalyticsEventInput, bool) {
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return commands.AnalyticsEventInput{}, false
	}

	var idString, typeString string
	var at int64
	if json.Unmarshal(fields["id"], &idString) != nil || json.Unmarshal(fields["type"], &typeString) != nil || json.Unmarshal(fields["at"], &at) != nil || at < 0 {
		return commands.AnalyticsEventInput{}, false
	}
	clientID, err := uuid.Parse(strings.TrimSpace(idString))
	if err != nil || clientID == uuid.Nil {
		return commands.AnalyticsEventInput{}, false
	}

	payload := make(map[string]string)
	for _, key := range []string{"section", "technology", "action", "kind", "target", "contact"} {
		if value, ok := fields[key]; ok {
			var stringValue string
			if json.Unmarshal(value, &stringValue) == nil {
				payload[key] = stringValue
			}
		}
	}
	return commands.AnalyticsEventInput{
		ClientEventID: clientID,
		Type:          domain.AnalyticsEventType(typeString),
		Payload:       payload,
		At:            at,
	}, true
}

func decodeAnalyticsClick(raw json.RawMessage) (commands.AnalyticsClickInput, bool) {
	var value struct {
		Section string `json:"section"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		VW      int    `json:"vw"`
		Target  string `json:"target"`
		At      int64  `json:"at"`
	}
	if json.Unmarshal(raw, &value) != nil || value.At < 0 {
		return commands.AnalyticsClickInput{}, false
	}
	return commands.AnalyticsClickInput{Section: value.Section, XPct: value.X, YPct: value.Y, VW: value.VW, Target: value.Target, At: value.At}, true
}
