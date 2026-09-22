package valueobjects

import (
	"fmt"
	"net/url"
	"strings"
)

type ExternalURL struct {
	value string
}

func NewExternalURL(value string) (ExternalURL, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return ExternalURL{}, fmt.Errorf("%w: url is empty", ErrInvalidValue)
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ExternalURL{}, fmt.Errorf("%w: invalid external url", ErrInvalidValue)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ExternalURL{}, fmt.Errorf("%w: unsupported url scheme", ErrInvalidValue)
	}
	return ExternalURL{value: value}, nil
}

func (u ExternalURL) String() string { return u.value }
func (u ExternalURL) IsZero() bool   { return u.value == "" }

type AssetPath struct {
	value string
}

func NewAssetPath(value string) (AssetPath, error) {
	value = strings.TrimSpace(value)
	if value == "" || !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") {
		return AssetPath{}, fmt.Errorf("%w: invalid asset path", ErrInvalidValue)
	}
	return AssetPath{value: value}, nil
}

func (p AssetPath) String() string { return p.value }
func (p AssetPath) IsZero() bool   { return p.value == "" }
