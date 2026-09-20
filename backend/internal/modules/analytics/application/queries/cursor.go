package queries

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

var (
	ErrInvalidAnalyticsCursor = errors.New("invalid analytics cursor")
	ErrInvalidAnalyticsLimit  = errors.New("invalid analytics limit")
)

func EncodeVisitCursor(t time.Time, id string) (string, error) {
	b, err := json.Marshal(struct {
		T  time.Time `json:"t"`
		ID string    `json:"id"`
	}{t, id})
	if err != nil {
		return "", fmt.Errorf("encode visit cursor: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func DecodeVisitCursor(value string) (time.Time, string, error) {
	if value == "" {
		return time.Time{}, "", nil
	}
	b, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return time.Time{}, "", ErrInvalidAnalyticsCursor
	}
	var v struct {
		T  time.Time `json:"t"`
		ID string    `json:"id"`
	}
	if json.Unmarshal(b, &v) != nil || v.T.IsZero() || v.ID == "" {
		return time.Time{}, "", ErrInvalidAnalyticsCursor
	}
	return v.T, v.ID, nil
}

func NormalizeVisitLimit(v int) (int, error) {
	if v == 0 {
		return 20, nil
	}
	if v < 1 || v > 50 {
		return 0, ErrInvalidAnalyticsLimit
	}
	return v, nil
}

func NormalizeSectionOrder(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}
