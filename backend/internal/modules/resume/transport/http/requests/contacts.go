package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactsRequest struct {
	Version        int64                `json:"version"`
	ID             string               `json:"id,omitempty"`
	Phone          *string              `json:"phone,omitempty"`
	Telegram       *string              `json:"telegram,omitempty"`
	Email          *string              `json:"email,omitempty"`
	HiddenContacts []string             `json:"hiddenContacts"`
	ContactsText   *ContactsTextRequest `json:"contactsText,omitempty"`
}

type ContactsTextRequest struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}

func (r ContactsRequest) ToInput() (commands.ContactsInput, error) {
	var telegram *domain.Telegram
	if r.Telegram != nil {
		value, err := domain.NewTelegram(*r.Telegram)
		if err != nil {
			return commands.ContactsInput{}, err
		}
		telegram = &value
	}
	var email *domain.Email
	if r.Email != nil {
		value, err := domain.NewEmail(*r.Email)
		if err != nil {
			return commands.ContactsInput{}, err
		}
		email = &value
	}
	hidden := make([]contactkey.Key, len(r.HiddenContacts))
	for i, value := range r.HiddenContacts {
		hidden[i] = contactkey.Key(value)
	}
	var text *commands.ContactsTextInput
	if r.ContactsText != nil {
		text = &commands.ContactsTextInput{Title: r.ContactsText.Title, Subtitle: r.ContactsText.Subtitle}
	}
	var id sharedvalueobjects.ID
	if r.ID != "" {
		parsed, err := sharedvalueobjects.ParseID(r.ID)
		if err != nil {
			return commands.ContactsInput{}, err
		}
		id = parsed
	}
	return commands.ContactsInput{
		ID: id, Phone: r.Phone, Telegram: telegram, Email: email,
		HiddenContacts: hidden, ContactsText: text,
	}, nil
}
