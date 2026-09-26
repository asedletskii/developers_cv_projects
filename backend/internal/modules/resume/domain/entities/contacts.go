package entities

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ContactsParams struct {
	Phone    string
	Telegram *valueobjects.Telegram
	Email    *valueobjects.Email
}

type Contacts struct {
	phone    *string
	telegram *valueobjects.Telegram
	email    *valueobjects.Email
}

func NewContacts(params ContactsParams) (Contacts, error) {
	phone := strings.TrimSpace(params.Phone)
	if phone == "" && params.Telegram == nil && params.Email == nil {
		return Contacts{}, fmt.Errorf("%w: at least one contact is required", valueobjects.ErrInvalidValue)
	}
	var phoneValue *string
	if phone != "" {
		phoneCopy := phone
		phoneValue = &phoneCopy
	}
	return Contacts{
		phone:    phoneValue,
		telegram: cloneTelegram(params.Telegram),
		email:    cloneEmail(params.Email),
	}, nil
}

func EmptyContacts() Contacts { return Contacts{} }

func (c Contacts) Phone() *string                   { return cloneString(c.phone) }
func (c Contacts) Telegram() *valueobjects.Telegram { return cloneTelegram(c.telegram) }
func (c Contacts) Email() *valueobjects.Email       { return cloneEmail(c.email) }

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
