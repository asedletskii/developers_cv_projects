package valueobjects

import (
	"fmt"
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{3}(?:[0-9a-fA-F]{3})?$`)

type Color struct {
	value string
}

func NewColor(value string) (Color, error) {
	value = strings.TrimSpace(value)
	if !hexColorPattern.MatchString(value) {
		return Color{}, fmt.Errorf("%w: invalid color", ErrInvalidValue)
	}
	return Color{value: value}, nil
}

func (c Color) String() string { return c.value }

type ThemeSettingsParams struct {
	Primary   *Color
	Secondary *Color
}

type ThemeSettings struct {
	Primary   *Color
	Secondary *Color
}

func NewThemeSettings(params ThemeSettingsParams) (ThemeSettings, error) {
	if params.Primary == nil && params.Secondary == nil {
		return ThemeSettings{}, fmt.Errorf("%w: theme is empty", ErrInvalidValue)
	}
	return ThemeSettings{
		Primary:   cloneColor(params.Primary),
		Secondary: cloneColor(params.Secondary),
	}, nil
}

func cloneColor(value *Color) *Color {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
