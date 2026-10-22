package contracts

import (
	"context"

	contactkey "developers-cv/backend/internal/shared/domain/contact"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactKey = contactkey.Key

const (
	ContactKeyPhone    = contactkey.Phone
	ContactKeyTelegram = contactkey.Telegram
	ContactKeyEmail    = contactkey.Email
)

type ContactViewQuery struct {
	ResumeID sharedvalueobjects.ID
	Contact  ContactKey
}

type ContactView struct {
	Contact ContactKey
	Value   string
}

type ContactViewQueryHandler interface {
	Handle(context.Context, ContactViewQuery) (ContactView, error)
}

type DownloadContactsQuery struct {
	ResumeID sharedvalueobjects.ID
}

type DownloadContactsView struct {
	Phone    *string
	Telegram *string
	Email    *string
}

type DownloadContactsQueryHandler interface {
	Handle(context.Context, DownloadContactsQuery) (DownloadContactsView, error)
}
