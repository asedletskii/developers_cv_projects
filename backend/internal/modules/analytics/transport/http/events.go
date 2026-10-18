package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/analytics/application/commands"
	"developers-cv/backend/internal/modules/analytics/domain"
	"developers-cv/backend/internal/modules/analytics/transport/http/requests"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	jsonhttp "developers-cv/backend/internal/platform/http/json"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const AnalyticsBatchMaxBodyBytes int64 = 32 * 1024

type EventsHandler struct {
	ingest *commands.IngestAnalyticsBatchHandler
}

func NewEventsHandler(ingest *commands.IngestAnalyticsBatchHandler) *EventsHandler {
	return &EventsHandler{ingest: ingest}
}

func (h *EventsHandler) Handle(c *gin.Context) {
	if h == nil || h.ingest == nil {
		apperrors.WriteError(c, apperrors.Internal(errors.New("analytics events handler is not configured")))
		return
	}

	linkID, ok := parseAnalyticsID(c, "linkId", "Invalid link id")
	if !ok {
		return
	}

	var body requests.AnalyticsBatchRequest
	if err := jsonhttp.BindJSON(c, &body); err != nil {
		h.ingest.RecordRejectedTransportBatch("invalid_batch")
		writeAnalyticsBindError(c, err)
		return
	}

	viewUUID, err := uuid.Parse(strings.TrimSpace(body.ViewID))
	if err != nil || viewUUID == uuid.Nil {
		h.ingest.RecordRejectedTransportBatch("invalid_batch")
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid view id"))
		return
	}
	viewID, err := sharedvalueobjects.IDFromUUID(viewUUID)
	if err != nil {
		h.ingest.RecordRejectedTransportBatch("invalid_batch")
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid view id"))
		return
	}
	if body.Seq == nil || body.Final == nil || body.State == nil || body.Events == nil || body.Clicks == nil {
		h.ingest.RecordRejectedTransportBatch("invalid_batch")
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid analytics batch"))
		return
	}
	if len(*body.Events) > commands.MaxAnalyticsEventsPerBatch || len(*body.Clicks) > commands.MaxAnalyticsClicksPerBatch {
		h.ingest.RecordRejectedTransportBatch("batch_limit")
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Analytics batch limits exceeded"))
		return
	}

	events := make([]commands.AnalyticsEventInput, 0, len(*body.Events))
	for _, raw := range *body.Events {
		event, ok, reason := decodeAnalyticsEvent(raw)
		if !ok {
			h.ingest.RecordDroppedTransportInput("event", reason)
			continue
		}
		events = append(events, event)
	}
	clicks := make([]commands.AnalyticsClickInput, 0, len(*body.Clicks))
	for _, raw := range *body.Clicks {
		click, ok, reason := decodeAnalyticsClick(raw)
		if !ok {
			h.ingest.RecordDroppedTransportInput("click", reason)
			continue
		}
		clicks = append(clicks, click)
	}

	err = h.ingest.Handle(c.Request.Context(), commands.AnalyticsBatchInput{
		ViewID: viewID, LinkID: linkID, Seq: *body.Seq,
		ActiveMS: body.State.ActiveMS, MaxSection: body.State.MaxSection, ExitSection: body.State.ExitSection,
		Final: *body.Final, Events: events, Clicks: clicks,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAnalyticsLinkInactive):
			apperrors.WriteError(c, apperrors.New(apperrors.CodeGone, http.StatusGone, "Resume link is no longer active"))
		case errors.Is(err, domain.ErrAnalyticsSubjectNotFound), errors.Is(err, domain.ErrAnalyticsViewExpired):
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

func decodeAnalyticsEvent(raw json.RawMessage) (commands.AnalyticsEventInput, bool, string) {
	var fields requests.AnalyticsEventRequest
	if json.Unmarshal(raw, &fields) != nil || fields == nil {
		return commands.AnalyticsEventInput{}, false, "invalid_json_shape"
	}

	var idString, typeString string
	var at int64
	if json.Unmarshal(fields["id"], &idString) != nil || json.Unmarshal(fields["type"], &typeString) != nil || json.Unmarshal(fields["at"], &at) != nil {
		return commands.AnalyticsEventInput{}, false, "invalid_fields"
	}
	if at < 0 {
		return commands.AnalyticsEventInput{}, false, "negative_timestamp"
	}
	clientID, err := uuid.Parse(strings.TrimSpace(idString))
	if err != nil || clientID == uuid.Nil {
		return commands.AnalyticsEventInput{}, false, "invalid_client_id"
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
		ClientEventID:     clientID,
		Type:              domain.AnalyticsEventType(typeString),
		Payload:           payload,
		ClientTimestampMS: at,
	}, true, ""
}

func decodeAnalyticsClick(raw json.RawMessage) (commands.AnalyticsClickInput, bool, string) {
	var value requests.AnalyticsClickRequest
	if json.Unmarshal(raw, &value) != nil {
		return commands.AnalyticsClickInput{}, false, "invalid_fields"
	}
	if strings.TrimSpace(string(raw)) == "null" || strings.TrimSpace(string(raw)) == "" {
		return commands.AnalyticsClickInput{}, false, "invalid_json_shape"
	}
	if value.At < 0 {
		return commands.AnalyticsClickInput{}, false, "negative_timestamp"
	}
	return commands.AnalyticsClickInput{
		Section: value.Section, XPct: value.X, YPct: value.Y, ViewportW: value.VW,
		Target: value.Target, ClientTimestampMS: value.At,
	}, true, ""
}
