package requests

import "encoding/json"

type AnalyticsBatchRequest struct {
	ViewID string                 `json:"viewId"`
	Seq    *int64                 `json:"seq"`
	Final  *bool                  `json:"final"`
	State  *AnalyticsStateRequest `json:"state"`
	Events *[]json.RawMessage     `json:"events"`
	Clicks *[]json.RawMessage     `json:"clicks"`
}

type AnalyticsStateRequest struct {
	ActiveMS    int64  `json:"activeMs"`
	MaxSection  string `json:"maxSection"`
	ExitSection string `json:"exitSection"`
}

type PublicContactViewRequest struct {
	Initiator      string  `json:"initiator"`
	Contact        *string `json:"contact,omitempty"`
	TurnstileToken string  `json:"turnstileToken"`
	ViewID         string  `json:"viewId,omitempty"`
}
