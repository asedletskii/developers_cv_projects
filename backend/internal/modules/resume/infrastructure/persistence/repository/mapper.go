package repository

import (
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func toResumeModel(p *domain.Resume) gormmodel.Resume {
	now := time.Now().UTC()
	return gormmodel.Resume{
		ID:          p.ID().UUID(),
		UserID:      p.UserID().UUID(),
		Name:        p.Name(),
		Surname:     p.Surname(),
		Patronymic:  p.Patronymic(),
		Position:    p.Position(),
		Title:       p.Title(),
		ResumeLink:  assetPathPtr(p.ResumeLink()),
		PhotoPath:   p.PhotoPath().String(),
		DateOfBirth: p.DateOfBirth().Time(),
		Location:    p.Location(),
		Relocation:  p.Relocation(),
		UpdatedAt:   now,
	}
}

func localizedPart(v *valueobjects.LocalizedText, language string) *string {
	if v == nil {
		return nil
	}
	value := v.Value(language)
	if value == "" {
		return nil
	}
	return &value
}
func localizedPtr(v *valueobjects.LocalizedText, language string) *string {
	return localizedPart(v, language)
}
func colorPtr(v *valueobjects.Color) *string {
	if v == nil {
		return nil
	}
	value := v.String()
	return &value
}
func assetPathPtr(v *valueobjects.AssetPath) *string {
	if v == nil {
		return nil
	}
	value := v.String()
	return &value
}
func externalURLPtr(v *valueobjects.ExternalURL) *string {
	if v == nil {
		return nil
	}
	value := v.String()
	return &value
}
func enumPtr[T ~string](v *T) *string {
	if v == nil {
		return nil
	}
	value := string(*v)
	return &value
}
func newGreeting(m gormmodel.Greeting) (valueobjects.Greeting, error) {
	f, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.FirstRU, EN: m.FirstEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	s, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.SecondRU, EN: m.SecondEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	h, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.SecondHighlightRU, EN: m.SecondHighlightEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	t, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.ThirdRU, EN: m.ThirdEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	return valueobjects.NewGreeting(valueobjects.GreetingParams{First: f, Second: s, SecondHighlight: h, Third: t}), nil
}

func parseYearMonthDate(value string) (valueobjects.YearMonth, error) {
	if len(value) < 7 {
		return valueobjects.YearMonth{}, fmt.Errorf("invalid year-month database value %q", value)
	}
	return valueobjects.ParseYearMonth(value[:7])
}

func parseYearDate(value string) (valueobjects.Year, error) {
	if len(value) < 4 {
		return valueobjects.Year{}, fmt.Errorf("invalid year database value %q", value)
	}
	return valueobjects.ParseYear(value[:4])
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}
