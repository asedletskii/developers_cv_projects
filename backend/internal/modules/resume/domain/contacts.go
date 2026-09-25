package domain

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ContactsParams struct {
	ID             valueobjects.ID
	ResumeID       valueobjects.ID
	Contacts       entities.Contacts
	HiddenContacts []enums.ContactKey
	Text           *valueobjects.ContactsText
}

type Contacts struct {
	id             valueobjects.ID
	resumeID       valueobjects.ID
	contacts       entities.Contacts
	hiddenContacts []enums.ContactKey
	text           *valueobjects.ContactsText
}

func NewContacts(params ContactsParams) (*Contacts, error) {
	if params.ID.IsZero() || params.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: invalid contacts identity", ErrInvalidResume)
	}
	seen := make(map[enums.ContactKey]struct{}, len(params.HiddenContacts))
	for _, key := range params.HiddenContacts {
		if !key.IsValid() {
			return nil, fmt.Errorf("%w: invalid hidden contact", ErrInvalidResume)
		}
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("%w: duplicate hidden contact", ErrInvalidResume)
		}
		seen[key] = struct{}{}
	}
	return &Contacts{id: params.ID, resumeID: params.ResumeID, contacts: params.Contacts, hiddenContacts: append([]enums.ContactKey(nil), params.HiddenContacts...), text: cloneContactsText(params.Text)}, nil
}

func (c *Contacts) Replace(params ContactsParams) error {
	params.ID, params.ResumeID = c.id, c.resumeID
	updated, err := NewContacts(params)
	if err != nil {
		return err
	}
	*c = *updated
	return nil
}

func (c *Contacts) ID() valueobjects.ID         { return c.id }
func (c *Contacts) ResumeID() valueobjects.ID   { return c.resumeID }
func (c *Contacts) Contacts() entities.Contacts { return c.contacts }
func (c *Contacts) HiddenContacts() []enums.ContactKey {
	return append([]enums.ContactKey(nil), c.hiddenContacts...)
}
func (c *Contacts) Text() *valueobjects.ContactsText { return cloneContactsText(c.text) }

func cloneContactsText(value *valueobjects.ContactsText) *valueobjects.ContactsText {
	if value == nil {
		return nil
	}
	copy, err := valueobjects.NewContactsText(valueobjects.ContactsTextParams{
		Title:    value.Title(),
		Subtitle: value.Subtitle(),
	})
	if err != nil {
		return nil
	}
	return &copy
}
