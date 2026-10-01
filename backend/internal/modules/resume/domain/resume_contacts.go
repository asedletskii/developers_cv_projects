package domain

import (
	"fmt"
	"net/mail"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactsParams struct {
	ID             sharedvalueobjects.ID
	ResumeID       sharedvalueobjects.ID
	Contacts       ContactMethods
	HiddenContacts []ContactKey
	Text           *ContactsText
}

type Contacts struct {
	id             sharedvalueobjects.ID
	resumeID       sharedvalueobjects.ID
	contacts       ContactMethods
	hiddenContacts []ContactKey
	text           *ContactsText
}

func NewContacts(params ContactsParams) (*Contacts, error) {
	if params.ID.IsZero() || params.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: invalid contacts identity", ErrInvalidResume)
	}
	if !params.Contacts.HasAny() {
		return nil, fmt.Errorf("%w: contacts are empty", ErrInvalidResume)
	}
	seen := make(map[ContactKey]struct{}, len(params.HiddenContacts))
	for _, key := range params.HiddenContacts {
		if !key.IsValid() {
			return nil, fmt.Errorf("%w: invalid hidden contact", ErrInvalidResume)
		}
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("%w: duplicate hidden contact", ErrInvalidResume)
		}
		if !params.Contacts.Has(key) {
			return nil, fmt.Errorf("%w: hidden contact is not filled", ErrResumeInvariant)
		}
		seen[key] = struct{}{}
	}
	return &Contacts{id: params.ID, resumeID: params.ResumeID, contacts: params.Contacts, hiddenContacts: append([]ContactKey(nil), params.HiddenContacts...), text: cloneContactsText(params.Text)}, nil
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

func (c *Contacts) ID() sharedvalueobjects.ID       { return c.id }
func (c *Contacts) ResumeID() sharedvalueobjects.ID { return c.resumeID }
func (c *Contacts) Contacts() ContactMethods        { return c.contacts }
func (c *Contacts) HiddenContacts() []ContactKey {
	return append([]ContactKey(nil), c.hiddenContacts...)
}
func (c *Contacts) Text() *ContactsText { return cloneContactsText(c.text) }

func cloneContactsText(value *ContactsText) *ContactsText {
	if value == nil {
		return nil
	}
	copy := value.Clone()
	return &copy
}

func (r *Resume) ReplaceContacts(params ContactsParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.contacts != nil {
			params.ID = r.contacts.ID()
		} else {
			params.ID = sharedvalueobjects.NewID()
		}
	} else if r.contacts == nil || params.ID != r.contacts.ID() {
		return fmt.Errorf("%w: contacts do not belong to resume", ErrResumeInvariant)
	}
	params.ResumeID = r.id
	contacts, err := NewContacts(params)
	if err != nil {
		return err
	}
	r.contacts = contacts
	return nil
}

func (r *Resume) Contacts() *Contacts {
	if r == nil || r.contacts == nil {
		return nil
	}
	copy := *r.contacts
	return &copy
}

type ContactKey string

const (
	ContactKeyTelegram ContactKey = "telegram"
	ContactKeyEmail    ContactKey = "email"
	ContactKeyPhone    ContactKey = "phone"
)

func (v ContactKey) IsValid() bool {
	switch v {
	case ContactKeyTelegram, ContactKeyEmail, ContactKeyPhone:
		return true
	default:
		return false
	}
}

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	parsed, err := mail.ParseAddress(value)
	if err != nil || parsed.Address != value || !strings.Contains(value, "@") {
		return Email{}, fmt.Errorf("%w: invalid email", ErrInvalidValue)
	}
	return Email{value: value}, nil
}

func (e Email) String() string { return e.value }

type Telegram struct {
	value string
}

func NewTelegram(value string) (Telegram, error) {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "@")
	if value == "" {
		return Telegram{}, fmt.Errorf("%w: telegram is empty", ErrInvalidValue)
	}
	for _, r := range value {
		if !isValidTelegramRune(r) {
			return Telegram{}, fmt.Errorf("%w: invalid telegram username", ErrInvalidValue)
		}
	}
	return Telegram{value: value}, nil
}

func isValidTelegramRune(r rune) bool {
	return r == '_' ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

func (t Telegram) String() string { return t.value }

type ContactMethodsParams struct {
	Phone    string
	Telegram *Telegram
	Email    *Email
}

type ContactMethods struct {
	phone    *string
	telegram *Telegram
	email    *Email
}

func NewContactMethods(params ContactMethodsParams) (ContactMethods, error) {
	phone := strings.TrimSpace(params.Phone)
	if phone == "" && params.Telegram == nil && params.Email == nil {
		return ContactMethods{}, fmt.Errorf("%w: at least one contact is required", ErrInvalidValue)
	}
	var phoneValue *string
	if phone != "" {
		phoneCopy := phone
		phoneValue = &phoneCopy
	}
	return ContactMethods{
		phone:    phoneValue,
		telegram: cloneTelegram(params.Telegram),
		email:    cloneEmail(params.Email),
	}, nil
}

func (c ContactMethods) HasAny() bool {
	return c.phone != nil || c.telegram != nil || c.email != nil
}

func (c ContactMethods) Has(key ContactKey) bool {
	switch key {
	case ContactKeyPhone:
		return c.phone != nil
	case ContactKeyTelegram:
		return c.telegram != nil
	case ContactKeyEmail:
		return c.email != nil
	default:
		return false
	}
}

func (c ContactMethods) Phone() *string      { return cloneString(c.phone) }
func (c ContactMethods) Telegram() *Telegram { return cloneTelegram(c.telegram) }
func (c ContactMethods) Email() *Email       { return cloneEmail(c.email) }

func cloneTelegram(value *Telegram) *Telegram {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneEmail(value *Email) *Email {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

type ContactsTextParams struct {
	Title    *string
	Subtitle *string
}

type ContactsText struct {
	title    *string
	subtitle *string
}

func NewContactsText(params ContactsTextParams) (ContactsText, error) {
	title := cloneString(params.Title)
	subtitle := cloneString(params.Subtitle)
	if title == nil && subtitle == nil {
		return ContactsText{}, fmt.Errorf("%w: contacts text is empty", ErrInvalidValue)
	}
	return ContactsText{title: title, subtitle: subtitle}, nil
}

func (c ContactsText) Clone() ContactsText {
	return ContactsText{title: cloneString(c.title), subtitle: cloneString(c.subtitle)}
}

func (c ContactsText) Title() *string    { return cloneString(c.title) }
func (c ContactsText) Subtitle() *string { return cloneString(c.subtitle) }
