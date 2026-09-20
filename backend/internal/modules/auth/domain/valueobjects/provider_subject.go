package valueobjects

import (
	"errors"
	"fmt"
	"strings"
)

const MaxProviderSubjectLength = 255

var ErrInvalidProviderSubject = errors.New("invalid provider subject")

type ProviderSubject struct {
	value string
}

func NewProviderSubject(value string) (ProviderSubject, error) {
	if strings.TrimSpace(value) == "" {
		return ProviderSubject{}, fmt.Errorf("%w: value is empty", ErrInvalidProviderSubject)
	}
	if len([]rune(value)) > MaxProviderSubjectLength {
		return ProviderSubject{}, fmt.Errorf("%w: value exceeds %d characters", ErrInvalidProviderSubject, MaxProviderSubjectLength)
	}
	return ProviderSubject{value: value}, nil
}

func (s ProviderSubject) String() string { return s.value }
func (s ProviderSubject) IsZero() bool   { return s.value == "" }
