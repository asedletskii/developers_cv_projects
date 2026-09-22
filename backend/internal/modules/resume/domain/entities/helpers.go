package entities

import "developers-cv/backend/internal/modules/resume/domain/valueobjects"

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneYear(value *valueobjects.Year) *valueobjects.Year {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneYearMonth(value *valueobjects.YearMonth) *valueobjects.YearMonth {
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

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}

func cloneTechnologies(values []Technology) []Technology {
	if values == nil {
		return nil
	}
	return append([]Technology{}, values...)
}

func cloneAssetPaths(values []valueobjects.AssetPath) []valueobjects.AssetPath {
	if values == nil {
		return nil
	}
	return append([]valueobjects.AssetPath{}, values...)
}

func cloneContributors(values []Contributor) []Contributor {
	if values == nil {
		return nil
	}
	return append([]Contributor{}, values...)
}

func cloneTelegram(value *valueobjects.Telegram) *valueobjects.Telegram {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneEmail(value *valueobjects.Email) *valueobjects.Email {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
