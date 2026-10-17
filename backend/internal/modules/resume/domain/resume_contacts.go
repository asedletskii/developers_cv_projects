package domain

import (
	"fmt"
	"net/mail"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactsParams struct {
	ID             sharedvalueobjects.ID
	Contacts       ContactMethods
	HiddenContacts []ContactKey
	Text           *ContactsText
}

type ContactsSnapshot struct {
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

func newContacts(params ContactsParams, resumeID sharedvalueobjects.ID) (*Contacts, error) {
	contacts := &Contacts{
		id: params.ID, resumeID: resumeID, contacts: params.Contacts,
		hiddenContacts: append([]ContactKey(nil), params.HiddenContacts...), text: cloneContactsText(params.Text),
	}
	if err := contacts.validate(); err != nil {
		return nil, err
	}
	return contacts, nil
}

func RestoreContacts(snapshot ContactsSnapshot) *Contacts {
	return &Contacts{
		id: snapshot.ID, resumeID: snapshot.ResumeID, contacts: snapshot.Contacts,
		hiddenContacts: append([]ContactKey(nil), snapshot.HiddenContacts...), text: cloneContactsText(snapshot.Text),
	}
}

func (c *Contacts) validate() error {
	if c == nil || c.id.IsZero() || c.resumeID.IsZero() {
		return fmt.Errorf("%w: invalid contacts identity", ErrInvalidResume)
	}
	if !c.contacts.hasAny() {
		return fmt.Errorf("%w: contacts are empty", ErrInvalidResume)
	}
	if err := c.contacts.validate(); err != nil {
		return err
	}
	if c.text != nil {
		if err := c.text.validate(); err != nil {
			return fmt.Errorf("%w: invalid contacts text", ErrInvalidResume)
		}
	}
	if len(c.hiddenContacts) > 3 {
		return fmt.Errorf("%w: too many hidden contacts", ErrInvalidResume)
	}
	seen := make(map[ContactKey]struct{}, len(c.hiddenContacts))
	for _, key := range c.hiddenContacts {
		if !key.isValid() {
			return fmt.Errorf("%w: invalid hidden contact", ErrInvalidResume)
		}
		if _, exists := seen[key]; exists {
			return fmt.Errorf("%w: duplicate hidden contact", ErrInvalidResume)
		}
		if !c.contacts.has(key) {
			return fmt.Errorf("%w: hidden contact is not filled", ErrResumeInvariant)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func (c *Contacts) ID() sharedvalueobjects.ID       { return c.id }
func (c *Contacts) ResumeID() sharedvalueobjects.ID { return c.resumeID }
func (c *Contacts) Contacts() ContactMethods        { return cloneContactMethods(c.contacts) }
func (c *Contacts) HiddenContacts() []ContactKey {
	return append([]ContactKey(nil), c.hiddenContacts...)
}
func (c *Contacts) Text() *ContactsText { return cloneContactsText(c.text) }

func cloneContactsText(value *ContactsText) *ContactsText {
	if value == nil {
		return nil
	}
	copy := value.clone()
	return &copy
}

func (r *Resume) ReplaceContacts(params ContactsParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	var currentID sharedvalueobjects.ID
	if r.contacts != nil {
		currentID = r.contacts.ID()
	}
	resolvedID, err := resolveSingletonID(params.ID, currentID, "contacts")
	if err != nil {
		return err
	}
	params.ID = resolvedID
	contacts, err := newContacts(params, r.id)
	if err != nil {
		return err
	}
	r.contacts = contacts
	return nil
}

func (r *Resume) Contacts() *Contacts {
	if r == nil {
		return nil
	}
	return cloneContacts(r.contacts)
}

type ContactKey string

const (
	ContactKeyTelegram ContactKey = "telegram"
	ContactKeyEmail    ContactKey = "email"
	ContactKeyPhone    ContactKey = "phone"
)

func (v ContactKey) isValid() bool {
	switch v {
	case ContactKeyTelegram, ContactKeyEmail, ContactKeyPhone:
		return true
	default:
		return false
	}
}

func ParseContactKey(value string) (ContactKey, error) {
	key := ContactKey(value)
	if !key.isValid() {
		return "", fmt.Errorf("%w: invalid contact key", ErrInvalidValue)
	}
	return key, nil
}

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	email := Email{value: strings.TrimSpace(value)}
	if err := email.validate(); err != nil {
		return Email{}, err
	}
	return email, nil
}

func RestoreEmail(value string) Email { return Email{value: value} }

func (e Email) validate() error {
	parsed, err := mail.ParseAddress(e.value)
	if err != nil || parsed.Address != e.value || !strings.Contains(e.value, "@") || !textLengthWithin(e.value, MaxContactEmailLength) {
		return fmt.Errorf("%w: invalid email", ErrInvalidValue)
	}
	return nil
}

func (e Email) String() string { return e.value }

func (e Email) IsZero() bool { return e.value == "" }

type Telegram struct {
	value string
}

func NewTelegram(value string) (Telegram, error) {
	telegram := Telegram{value: strings.TrimPrefix(strings.TrimSpace(value), "@")}
	if err := telegram.validate(); err != nil {
		return Telegram{}, err
	}
	return telegram, nil
}

func RestoreTelegram(value string) Telegram { return Telegram{value: value} }

func (t Telegram) validate() error {
	if t.value == "" || !textLengthWithin(t.value, MaxContactTelegramLength) {
		return fmt.Errorf("%w: telegram is empty or too long", ErrInvalidValue)
	}
	for _, r := range t.value {
		if !isValidTelegramRune(r) {
			return fmt.Errorf("%w: invalid telegram username", ErrInvalidValue)
		}
	}
	return nil
}

func isValidTelegramRune(r rune) bool {
	return r == '_' ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

func (t Telegram) String() string { return t.value }

func (t Telegram) IsZero() bool { return t.value == "" }

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
	methods := ContactMethods{telegram: cloneTelegram(params.Telegram), email: cloneEmail(params.Email)}
	phone := strings.TrimSpace(params.Phone)
	if phone != "" {
		methods.phone = &phone
	}
	if err := methods.validate(); err != nil {
		return ContactMethods{}, err
	}
	return methods, nil
}

func RestoreContactMethods(params ContactMethodsParams) ContactMethods {
	var phone *string
	if params.Phone != "" {
		value := params.Phone
		phone = &value
	}
	return ContactMethods{phone: phone, telegram: cloneTelegram(params.Telegram), email: cloneEmail(params.Email)}
}

func (c ContactMethods) validate() error {
	if !optionalTextLengthWithin(c.phone, MaxContactPhoneLength) {
		return fmt.Errorf("%w: phone exceeds %d characters", ErrInvalidValue, MaxContactPhoneLength)
	}
	if !c.hasAny() {
		return fmt.Errorf("%w: at least one contact is required", ErrInvalidValue)
	}
	if c.telegram != nil {
		if err := c.telegram.validate(); err != nil {
			return err
		}
	}
	if c.email != nil {
		if err := c.email.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c ContactMethods) hasAny() bool {
	return c.phone != nil || c.telegram != nil || c.email != nil
}

func (c ContactMethods) has(key ContactKey) bool {
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

func (c ContactMethods) Phone() *string      { return copyString(c.phone) }
func (c ContactMethods) Telegram() *Telegram { return cloneTelegram(c.telegram) }
func (c ContactMethods) Email() *Email       { return cloneEmail(c.email) }

func cloneContactMethods(value ContactMethods) ContactMethods {
	return ContactMethods{phone: copyString(value.phone), telegram: cloneTelegram(value.telegram), email: cloneEmail(value.email)}
}

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
	text := ContactsText{title: normalizeOptionalString(params.Title), subtitle: normalizeOptionalString(params.Subtitle)}
	if err := text.validate(); err != nil {
		return ContactsText{}, err
	}
	return text, nil
}

func RestoreContactsText(params ContactsTextParams) ContactsText {
	return ContactsText{title: copyString(params.Title), subtitle: copyString(params.Subtitle)}
}

func (c ContactsText) validate() error {
	if c.title == nil && c.subtitle == nil {
		return fmt.Errorf("%w: contacts text is empty", ErrInvalidValue)
	}
	if !optionalTextLengthWithin(c.title, MaxShortTextLength) || !optionalTextLengthWithin(c.subtitle, MaxShortTextLength) {
		return fmt.Errorf("%w: contacts text exceeds configured limits", ErrInvalidValue)
	}
	return nil
}

func (c ContactsText) clone() ContactsText {
	return ContactsText{title: copyString(c.title), subtitle: copyString(c.subtitle)}
}

func (c ContactsText) Title() *string    { return copyString(c.title) }
func (c ContactsText) Subtitle() *string { return copyString(c.subtitle) }

func cloneContacts(value *Contacts) *Contacts {
	if value == nil {
		return nil
	}
	copy := *value
	copy.contacts = cloneContactMethods(value.contacts)
	copy.hiddenContacts = append([]ContactKey(nil), value.hiddenContacts...)
	copy.text = cloneContactsText(value.text)
	return &copy
}
