package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ContactsToDBModels(c *domain.Contacts) (
	gormmodel.ResumeContacts,
	[]gormmodel.ResumeContactHidden,
	*gormmodel.ResumeContactsText,
	error,
) {
	if c == nil {
		return gormmodel.ResumeContacts{}, nil, nil, fmt.Errorf("contacts is nil")
	}
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
	return m, h, text, nil
}

func ContactsToDomain(
	m gormmodel.ResumeContacts,
	h []gormmodel.ResumeContactHidden,
	t *gormmodel.ResumeContactsText,
) (*domain.Contacts, error) {
	var tg *domain.Telegram
	if m.Telegram != nil {
		x, err := domain.NewTelegram(*m.Telegram)
		if err != nil {
			return nil, err
		}
		tg = &x
	}
	var em *domain.Email
	if m.Email != nil {
		x, err := domain.NewEmail(*m.Email)
		if err != nil {
			return nil, err
		}
		em = &x
	}
	contactValues := domain.ContactMethodsParams{
		Phone:    PointerStringValue(m.Phone),
		Telegram: tg,
		Email:    em,
	}
	contactMethods, err := domain.NewContactMethods(contactValues)
	if err != nil {
		return nil, err
	}
	hidden := make([]domain.ContactKey, len(h))
	for i, v := range h {
		hidden[i] = domain.ContactKey(v.Value)
	}
	var text *domain.ContactsText
	if t != nil {
		x, err := domain.NewContactsText(domain.ContactsTextParams{
			Title:    t.Title,
			Subtitle: t.Subtitle,
		})
		if err != nil {
			return nil, err
		}
		text = &x
	}
	id, err := identity.IDFromUUID(m.ID)
	if err != nil {
		return nil, err
	}
	rid, err := identity.IDFromUUID(m.ResumeID)
	if err != nil {
		return nil, err
	}
	return domain.NewContacts(domain.ContactsParams{
		ID:             id,
		ResumeID:       rid,
		Contacts:       contactMethods,
		HiddenContacts: hidden,
		Text:           text,
	})
}

func PointerStringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
