package valueobjects

import (
	"fmt"
	"time"
)

type Date struct {
	value time.Time
}

func NewDate(value time.Time) (Date, error) {
	if value.IsZero() {
		return Date{}, ErrInvalidValue
	}
	return Date{value: time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)}, nil
}

func ParseDate(value string) (Date, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return Date{}, fmt.Errorf("%w: invalid date %q", ErrInvalidValue, value)
	}
	return NewDate(parsed)
}

func (d Date) Time() time.Time        { return d.value }
func (d Date) IsZero() bool           { return d.value.IsZero() }
func (d Date) String() string         { return d.value.Format("2006-01-02") }
func (d Date) Before(other Date) bool { return d.value.Before(other.value) }
func (d Date) After(other Date) bool  { return d.value.After(other.value) }
func (d Date) Equal(other Date) bool  { return d.value.Equal(other.value) }

type YearMonth struct {
	value time.Time
}

func NewYearMonth(value time.Time) (YearMonth, error) {
	if value.IsZero() {
		return YearMonth{}, ErrInvalidValue
	}
	return YearMonth{value: time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)}, nil
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
func (ym YearMonth) Before(other YearMonth) bool { return ym.value.Before(other.value) }
func (ym YearMonth) After(other YearMonth) bool  { return ym.value.After(other.value) }
func (ym YearMonth) Equal(other YearMonth) bool  { return ym.value.Equal(other.value) }

type Year struct {
	value int
}

func NewYear(value int) (Year, error) {
	if value < 1 || value > 9999 {
		return Year{}, ErrInvalidValue
	}
	return Year{value: value}, nil
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
