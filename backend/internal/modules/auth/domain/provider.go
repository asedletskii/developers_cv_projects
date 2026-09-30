package domain

import (
	"fmt"
	"strings"
)

type Provider string

const ProviderTelegram Provider = "telegram"

func NewProvider(value string) (Provider, error) {
	value = strings.TrimSpace(value)
	provider := Provider(value)
	if provider == "" {
		return "", fmt.Errorf("%w: provider is empty", ErrInvalidProvider)
	}
	if provider != ProviderTelegram {
		return "", fmt.Errorf("%w: unsupported provider %q", ErrInvalidProvider, provider)
	}
	return provider, nil
}

func (p Provider) IsValid() bool {
	return p == ProviderTelegram
}

func (p Provider) String() string { return string(p) }
