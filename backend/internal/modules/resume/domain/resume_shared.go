package domain

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func resolveSingletonID(requestedID, currentID sharedvalueobjects.ID, entityName string) (sharedvalueobjects.ID, error) {
	if requestedID.IsZero() {
		if !currentID.IsZero() {
			return currentID, nil
		}
		return sharedvalueobjects.NewID(), nil
	}
	if currentID.IsZero() || requestedID != currentID {
		return sharedvalueobjects.ID{}, fmt.Errorf("%w: %s does not belong to resume", ErrResumeInvariant, entityName)
	}
	return currentID, nil
}

func resolveCollectionItemID(requestedID sharedvalueobjects.ID, existingIDs map[sharedvalueobjects.ID]struct{}, entityName string) (sharedvalueobjects.ID, error) {
	if requestedID.IsZero() {
		return sharedvalueobjects.NewID(), nil
	}
	if _, exists := existingIDs[requestedID]; !exists {
		return sharedvalueobjects.ID{}, fmt.Errorf("%w: %w: %s does not belong to resume", ErrResumeInvariant, ErrEntityNotFound, entityName)
	}
	return requestedID, nil
}

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func copyString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneYearMonth(value *YearMonth) *YearMonth {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneYear(value *Year) *Year {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneExternalURL(value *ExternalURL) *ExternalURL {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneExternalURLs(values []ExternalURL) []ExternalURL {
	if values == nil {
		return nil
	}
	return append([]ExternalURL(nil), values...)
}

func cloneSectionSettings(values []SectionSettings) []SectionSettings {
	if values == nil {
		return nil
	}
	return append([]SectionSettings(nil), values...)
}

var ErrInvalidValue = errors.New("invalid value")

type Date struct {
	value time.Time
}

func NewDate(value time.Time) (Date, error) {
	date := RestoreDate(value)
	if err := date.validate(); err != nil {
		return Date{}, err
	}
	return date, nil
}

func RestoreDate(value time.Time) Date {
	if value.IsZero() {
		return Date{}
	}
	return Date{value: time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)}
}

func (d Date) validate() error {
	if d.value.IsZero() {
		return ErrInvalidValue
	}
	return nil
}

func ParseDate(value string) (Date, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return Date{}, fmt.Errorf("%w: invalid date %q", ErrInvalidValue, value)
	}
	return NewDate(parsed)
}

func (d Date) Time() time.Time { return d.value }
func (d Date) IsZero() bool    { return d.value.IsZero() }
func (d Date) String() string  { return d.value.Format("2006-01-02") }

type YearMonth struct {
	value time.Time
}

func NewYearMonth(value time.Time) (YearMonth, error) {
	yearMonth := RestoreYearMonth(value)
	if err := yearMonth.validate(); err != nil {
		return YearMonth{}, err
	}
	return yearMonth, nil
}

func RestoreYearMonth(value time.Time) YearMonth {
	if value.IsZero() {
		return YearMonth{}
	}
	return YearMonth{value: time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)}
}

func (ym YearMonth) validate() error {
	if ym.value.IsZero() {
		return ErrInvalidValue
	}
	return nil
}

func ParseYearMonth(value string) (YearMonth, error) {
	parsed, err := time.Parse("2006-01", value)
	if err != nil {
		return YearMonth{}, fmt.Errorf("%w: invalid year-month %q", ErrInvalidValue, value)
	}
	return NewYearMonth(parsed)
}

func (ym YearMonth) Time() time.Time             { return ym.value }
func (ym YearMonth) IsZero() bool                { return ym.value.IsZero() }
func (ym YearMonth) String() string              { return ym.value.Format("2006-01") }
func (ym YearMonth) before(other YearMonth) bool { return ym.value.Before(other.value) }

type Year struct {
	value int
}

func NewYear(value int) (Year, error) {
	year := RestoreYear(value)
	if err := year.validate(); err != nil {
		return Year{}, err
	}
	return year, nil
}

func RestoreYear(value int) Year { return Year{value: value} }

func (y Year) validate() error {
	if y.value < 1 || y.value > 9999 {
		return ErrInvalidValue
	}
	return nil
}

func ParseYear(value string) (Year, error) {
	parsed, err := time.Parse("2006", value)
	if err != nil {
		return Year{}, fmt.Errorf("%w: invalid year %q", ErrInvalidValue, value)
	}
	return NewYear(parsed.Year())
}

func (y Year) Int() int       { return y.value }
func (y Year) String() string { return fmt.Sprintf("%04d", y.value) }
func (y Year) IsZero() bool   { return y.value == 0 }

type ExternalURL struct {
	value string
}

const MaxExternalURLLength = 2048

func NewExternalURL(value string) (ExternalURL, error) {
	urlValue := ExternalURL{value: strings.TrimSpace(value)}
	if err := urlValue.validate(); err != nil {
		return ExternalURL{}, err
	}
	return urlValue, nil
}

func RestoreExternalURL(value string) ExternalURL { return ExternalURL{value: value} }

func (u ExternalURL) validate() error {
	if u.value == "" {
		return fmt.Errorf("%w: url is empty", ErrInvalidValue)
	}
	if len(u.value) > MaxExternalURLLength {
		return fmt.Errorf("%w: url exceeds %d bytes", ErrInvalidValue, MaxExternalURLLength)
	}
	parsed, err := url.ParseRequestURI(u.value)
	if err != nil || parsed.Scheme == "" || parsed.Hostname() == "" {
		return fmt.Errorf("%w: invalid external url", ErrInvalidValue)
	}
	if parsed.Scheme != "https" {
		return fmt.Errorf("%w: external url must use https", ErrInvalidValue)
	}
	if parsed.User != nil {
		return fmt.Errorf("%w: external url must not contain user information", ErrInvalidValue)
	}
	return nil
}

func (u ExternalURL) String() string { return u.value }
func (u ExternalURL) IsZero() bool   { return u.value == "" }

func cloneDate(value *Date) *Date {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
