package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func ValidateResumeTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w: resume title is empty", ErrInvalidResume)
	}
	return nil
}

func validateSections(sections []valueobjects.SectionSettings) error {
	seen := make(map[string]struct{}, len(sections))
	for _, section := range sections {
		key := string(section.ID())
		if _, ok := seen[key]; ok {
			return fmt.Errorf("%w: duplicate resume section", ErrInvalidResume)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func cloneTheme(value *valueobjects.ThemeSettings) *valueobjects.ThemeSettings {
	if value == nil {
		return nil
	}
	copy := *value
	if value.Primary != nil {
		primary := *value.Primary
		copy.Primary = &primary
	}
	if value.Secondary != nil {
		secondary := *value.Secondary
		copy.Secondary = &secondary
	}
	return &copy
}

func cloneAvailability(value *valueobjects.Availability) *valueobjects.Availability {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneAssetPath(value *valueobjects.AssetPath) *valueobjects.AssetPath {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneSections(values []valueobjects.SectionSettings) []valueobjects.SectionSettings {
	if values == nil {
		return nil
	}
	return append([]valueobjects.SectionSettings(nil), values...)
}
