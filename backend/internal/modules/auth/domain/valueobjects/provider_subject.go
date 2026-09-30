package valueobjects

import (
	"fmt"
	"strings"
)

const MaxProviderSubjectLength = 255

type ProviderSubject struct {
	value string
}

func NewProviderSubject(value string) (ProviderSubject, error) {
	if strings.TrimSpace(value) == "" {
		return ProviderSubject{}, fmt.Errorf("provider subject is empty")
	}
	if len([]rune(value)) > MaxProviderSubjectLength {
		return ProviderSubject{}, fmt.Errorf("provider subject exceeds %d characters", MaxProviderSubjectLength)
	}
	return ProviderSubject{value: value}, nil
}

func (s ProviderSubject) String() string { return s.value }
func (s ProviderSubject) IsZero() bool   { return s.value == "" }
