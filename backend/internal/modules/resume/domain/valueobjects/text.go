package valueobjects

import (
	"fmt"
	"strings"
)

type LocalizedTextParams struct {
	RU string
	EN string
}

type LocalizedText struct {
	ru string
	en string
}

func NewLocalizedText(params LocalizedTextParams) (LocalizedText, error) {
	ru := strings.TrimSpace(params.RU)
	en := strings.TrimSpace(params.EN)
	if ru == "" && en == "" {
		return LocalizedText{}, fmt.Errorf("%w: both localized values are empty", ErrInvalidValue)
	}
	return LocalizedText{ru: ru, en: en}, nil
}

func (t LocalizedText) RU() string { return t.ru }
func (t LocalizedText) EN() string { return t.en }

func (t LocalizedText) Value(language string) string {
	if language == "en" && t.en != "" {
		return t.en
	}
	return t.ru
}
