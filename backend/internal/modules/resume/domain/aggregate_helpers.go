package domain

import (
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func cloneYear(value *valueobjects.Year) *valueobjects.Year {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneExternalURL(value *valueobjects.ExternalURL) *valueobjects.ExternalURL {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneAssetPaths(values []valueobjects.AssetPath) []valueobjects.AssetPath {
	if values == nil {
		return nil
	}
	return append([]valueobjects.AssetPath(nil), values...)
}
