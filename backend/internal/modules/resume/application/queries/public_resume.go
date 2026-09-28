package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type PublicResumeQuery struct {
	ID valueobjects.ID
}

type PublicResumeReader interface {
	GetPublic(context.Context, PublicResumeQuery) (ResumeView, error)
}

type PublicResumeQueryHandler struct {
	reader PublicResumeReader
}

func NewPublicResumeQueryHandler(reader PublicResumeReader) *PublicResumeQueryHandler {
	return &PublicResumeQueryHandler{reader: reader}
}

func (h *PublicResumeQueryHandler) Handle(ctx context.Context, query PublicResumeQuery) (ResumeView, error) {
	if h == nil || h.reader == nil {
		return ResumeView{}, fmt.Errorf("public resume query handler dependencies are nil")
	}
	if query.ID.IsZero() {
		return ResumeView{}, fmt.Errorf("invalid resume id")
	}

	view, err := h.reader.GetPublic(ctx, query)
	if err != nil {
		return ResumeView{}, err
	}

	return toPublic(view), nil
}

func toPublic(view ResumeView) ResumeView {
	result := view
	result.Version = nil
	if view.Contacts != nil {
		contacts := *view.Contacts
		contacts.HiddenContacts = append([]string(nil), view.Contacts.HiddenContacts...)
		contacts.ContactsText = clonePublicContactsText(view.Contacts.ContactsText)
		contacts.Phone = cloneStringPointer(view.Contacts.Phone)
		contacts.Telegram = cloneStringPointer(view.Contacts.Telegram)
		contacts.Email = cloneStringPointer(view.Contacts.Email)
		redactHiddenContacts(&contacts)
		result.Contacts = &contacts
	}
	result.Projects = make([]ProjectView, len(view.Projects))

	for index, project := range view.Projects {
		result.Projects[index] = project
		result.Projects[index].Role = nil
		result.Projects[index].GitLink = nil
		result.Projects[index].Screenshots = make([]string, 0)
		result.Projects[index].Contributors = make([]ProjectContributorView, 0)
	}

	return result
}

func redactHiddenContacts(contacts *ResumeContactsView) {
	if contacts == nil {
		return
	}

	for _, key := range contacts.HiddenContacts {
		switch enums.ContactKey(key) {
		case enums.ContactKeyPhone:
			contacts.Phone = nil
		case enums.ContactKeyTelegram:
			contacts.Telegram = nil
		case enums.ContactKeyEmail:
			contacts.Email = nil
		}
	}
}

func clonePublicContactsText(value *ResumeContactsTextView) *ResumeContactsTextView {
	if value == nil {
		return nil
	}
	result := *value
	result.Title = cloneStringPointer(value.Title)
	result.Subtitle = cloneStringPointer(value.Subtitle)
	return &result
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}
