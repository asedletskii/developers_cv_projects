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

type AnalyticsClickRequest struct {
	Section string `json:"section"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	VW      int    `json:"vw"`
	Target  string `json:"target"`
	At      int64  `json:"at"`
}

type AnalyticsEventRequest map[string]json.RawMessage
