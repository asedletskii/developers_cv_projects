package entities

import "developers-cv/backend/internal/modules/resume/domain/valueobjects"

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
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
