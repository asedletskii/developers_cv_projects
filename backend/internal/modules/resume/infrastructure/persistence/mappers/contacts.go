package mappers

import (
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ContactsToDBModels(c *domain.Contacts) (
	gormmodel.ResumeContacts,
	[]gormmodel.ResumeContactHidden,
	*gormmodel.ResumeContactsText,
) {
	d := c.Contacts()
	m := gormmodel.ResumeContacts{ID: c.ID().UUID(), ResumeID: c.ResumeID().UUID()}
	if v := d.Phone(); v != nil {
		m.Phone = v
	}
	if v := d.Telegram(); v != nil {
		x := v.String()
		m.Telegram = &x
	}
	if v := d.Email(); v != nil {
		x := v.String()
		m.Email = &x
	}
	hiddenValues := c.HiddenContacts()
	h := make([]gormmodel.ResumeContactHidden, 0, len(hiddenValues))
	for i, v := range hiddenValues {
		h = append(h, gormmodel.ResumeContactHidden{
			ContactsID: m.ID,
			Value:      string(v),
			Position:   i,
		})
	}
	var text *gormmodel.ResumeContactsText
	if x := c.Text(); x != nil {
		text = &gormmodel.ResumeContactsText{ContactsID: m.ID}
		if v := x.Title(); v != nil {
			value := *v
			text.Title = &value
		}
		if v := x.Subtitle(); v != nil {
			value := *v
			text.Subtitle = &value
		}
	}
	return m, h, text
}

func ContactsToDomain(
	m gormmodel.ResumeContacts,
	h []gormmodel.ResumeContactHidden,
	t *gormmodel.ResumeContactsText,
) (*domain.Contacts, error) {
	var telegram *domain.Telegram
	if m.Telegram != nil {
		value := domain.RestoreTelegram(*m.Telegram)
		telegram = &value
	}
	var email *domain.Email
	if m.Email != nil {
		value := domain.RestoreEmail(*m.Email)
		email = &value
	}
	contactMethods := domain.RestoreContactMethods(domain.ContactMethodsParams{
		Phone: PointerStringValue(m.Phone), Telegram: telegram, Email: email,
	})
	hidden := make([]contactkey.Key, len(h))
	for i, value := range h {
		hidden[i] = contactkey.Key(value.Value)
	}
	var text *domain.ContactsText
	if t != nil {
		value := domain.RestoreContactsText(domain.ContactsTextParams{Title: t.Title, Subtitle: t.Subtitle})
		text = &value
	}
	id, err := sharedvalueobjects.IDFromUUID(m.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(m.ResumeID)
	if err != nil {
		return nil, err
	}
	return domain.RestoreContacts(domain.ContactsSnapshot{
		ID: id, ResumeID: resumeID, Contacts: contactMethods, HiddenContacts: hidden, Text: text,
	}), nil
}

func PointerStringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
