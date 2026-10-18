package mappers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/analytics/domain"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ViewToDB(view *domain.View) (gormmodel.AnalyticsView, error) {
	if view == nil {
		return gormmodel.AnalyticsView{}, fmt.Errorf("analytics view is nil")
	}
	layout, err := json.Marshal(view.Layout())
	if err != nil {
		return gormmodel.AnalyticsView{}, fmt.Errorf("marshal layout: %w", err)
	}
	state := view.State()
	return gormmodel.AnalyticsView{
		ID:            view.ID().UUID(),
		LinkID:        view.LinkID().UUID(),
		VisitorID:     view.VisitorID().UUID(),
		IPHash:        view.IPHash(),
		Device:        string(view.Device()),
		Browser:       string(view.Browser()),
		OS:            string(view.OS()),
		ReferrerHost:  view.ReferrerHost(),
		Traffic:       string(view.Traffic()),
		TrafficReason: view.TrafficReason(),
		Layout:        layout,
		StartedAt:     view.StartedAt(),
		ConfirmedAt:   nullableTime(state.ConfirmedAt()),
		LastSeenAt:    nullableTime(state.LastSeenAt()),
		ActiveMS:      state.ActiveMS(),
		LastStateSeq:  state.LastSeq(),
		MaxSection:    string(state.MaxSection()),
		ExitSection:   string(state.ExitSection()),
		EndedAt:       nullableTime(state.EndedAt()),
		CreatedAt:     view.CreatedAt(),
	}, nil
}

func ViewToDomain(model gormmodel.AnalyticsView) (*domain.View, error) {
	var raw struct {
		Order        []domain.AnalyticsSectionID          `json:"order"`
		Hashes       map[domain.AnalyticsSectionID]string `json:"hashes"`
		Technologies *[]string                            `json:"technologies"`
	}
	if err := json.Unmarshal(model.Layout, &raw); err != nil {
		return nil, fmt.Errorf("unmarshal layout: %w", err)
	}
	var layout domain.LayoutSnapshot
	var err error
	if raw.Technologies != nil {
		layout, err = domain.RestoreLayoutSnapshot(raw.Order, raw.Hashes, *raw.Technologies)
	} else {
		layout, err = domain.RestoreLayoutSnapshot(raw.Order, raw.Hashes)
	}
	if err != nil {
		return nil, err
	}
	visitorID, err := valueObjectID(model.VisitorID)
	if err != nil {
		return nil, fmt.Errorf("map visitor id: %w", err)
	}
	id, err := valueObjectID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map view id: %w", err)
	}
	linkID, err := valueObjectID(model.LinkID)
	if err != nil {
		return nil, fmt.Errorf("map link id: %w", err)
	}
	return domain.RestoreView(domain.ViewSnapshot{
		ID: id, LinkID: linkID, VisitorID: visitorID, IPHash: model.IPHash,
		Device: domain.DeviceType(model.Device), Browser: domain.BrowserType(model.Browser), OS: domain.OSType(model.OS),
		ReferrerHost: model.ReferrerHost, Traffic: domain.TrafficType(model.Traffic),
		TrafficReason: model.TrafficReason, Layout: layout, StartedAt: model.StartedAt,
		CreatedAt: model.CreatedAt, State: domain.ViewStateSnapshot{
			LastSeq: model.LastStateSeq, ConfirmedAt: timeValue(model.ConfirmedAt),
			LastSeenAt: timeValue(model.LastSeenAt), ActiveMS: model.ActiveMS,
			MaxSection:  domain.AnalyticsSectionID(model.MaxSection),
			ExitSection: domain.AnalyticsSectionID(model.ExitSection), EndedAt: timeValue(model.EndedAt),
		},
	})
}

func timeValue(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.UTC()
}

func nullableTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	value = value.UTC()
	return &value
}

func valueObjectID(value uuid.UUID) (sharedvalueobjects.ID, error) {
	return sharedvalueobjects.IDFromUUID(value)
}
