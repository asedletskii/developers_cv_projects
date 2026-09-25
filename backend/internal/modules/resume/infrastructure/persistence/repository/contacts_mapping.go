package repository

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

func contactsModels(c *domain.Contacts) (gormmodel.ResumeContacts, []gormmodel.ResumeContactHidden, *gormmodel.ResumeContactsText, error) {
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
	h := make([]gormmodel.ResumeContactHidden, 0, len(c.HiddenContacts()))
	for i, v := range c.HiddenContacts() {
		h = append(h, gormmodel.ResumeContactHidden{ContactsID: m.ID, Value: string(v), Position: i})
	}
	var t *gormmodel.ResumeContactsText
	if x := c.Text(); x != nil {
		t = &gormmodel.ResumeContactsText{ContactsID: m.ID}
		if v := x.Title(); v != nil {
			ru, en := v.RU(), v.EN()
			t.TitleRU = &ru
			t.TitleEN = &en
		}
		if v := x.Subtitle(); v != nil {
			ru, en := v.RU(), v.EN()
			t.SubtitleRU = &ru
			t.SubtitleEN = &en
		}
	}
	return m, h, t, nil
}
func mapContacts(m gormmodel.ResumeContacts, h []gormmodel.ResumeContactHidden, t *gormmodel.ResumeContactsText) (*domain.Contacts, error) {
	var tg *valueobjects.Telegram
	if m.Telegram != nil {
		x, e := valueobjects.NewTelegram(*m.Telegram)
		if e != nil {
			return nil, e
		}
		tg = &x
	}
	var em *valueobjects.Email
	if m.Email != nil {
		x, e := valueobjects.NewEmail(*m.Email)
		if e != nil {
			return nil, e
		}
		em = &x
	}
	contacts, e := entities.NewContacts(entities.ContactsParams{Phone: pointerStringValue(m.Phone), Telegram: tg, Email: em})
	if e != nil {
		return nil, e
	}
	hidden := make([]enums.ContactKey, len(h))
	for i, v := range h {
		hidden[i] = enums.ContactKey(v.Value)
	}
	var text *valueobjects.ContactsText
	if t != nil {
		var title, subtitle *valueobjects.LocalizedText
		if t.TitleRU != nil || t.TitleEN != nil {
			ru, en := "", ""
			if t.TitleRU != nil {
				ru = *t.TitleRU
			}
			if t.TitleEN != nil {
				en = *t.TitleEN
			}
			x, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: ru, EN: en})
			if e != nil {
				return nil, e
			}
			title = &x
		}
		if t.SubtitleRU != nil || t.SubtitleEN != nil {
			ru, en := "", ""
			if t.SubtitleRU != nil {
				ru = *t.SubtitleRU
			}
			if t.SubtitleEN != nil {
				en = *t.SubtitleEN
			}
			x, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: ru, EN: en})
			if e != nil {
				return nil, e
			}
			subtitle = &x
		}
		x, e := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: title, Subtitle: subtitle})
		if e != nil {
			return nil, e
		}
		text = &x
	}
	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return nil, e
	}
	rid, e := valueobjects.IDFromUUID(m.ResumeID)
	if e != nil {
		return nil, e
	}
	return domain.NewContacts(domain.ContactsParams{ID: id, ResumeID: rid, Contacts: contacts, HiddenContacts: hidden, Text: text})
}

func pointerStringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
