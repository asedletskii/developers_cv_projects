package queries

import (
	"context"
	"fmt"

	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactViewThroughLinkQuery struct {
	LinkID    sharedvalueobjects.ID
	Initiator ContactViewInitiator
	Contact   resumedomain.ContactKey
	ViewID    sharedvalueobjects.ID
}

type ContactViewInitiator string

const (
	ContactViewInitiatorContact  ContactViewInitiator = "contact"
	ContactViewInitiatorDownload ContactViewInitiator = "download"
)

type ContactViewThroughLinkResult struct {
	Contact  *resumedomain.ContactKey         `json:"contact,omitempty"`
	Value    *string                          `json:"value,omitempty"`
	Contacts *[]ContactViewThroughLinkContact `json:"contacts,omitempty"`
}

type ContactViewThroughLinkContact struct {
	Contact resumedomain.ContactKey `json:"contact"`
	Value   string                  `json:"value"`
}

type ResumeContactViewQuery interface {
	Handle(context.Context, resumequeries.ContactViewQuery) (resumequeries.ContactView, error)
}

type ResumeContactsDownloadQuery interface {
	Handle(context.Context, resumequeries.DownloadContactsQuery) (resumequeries.DownloadContactsView, error)
}

type ContactViewAnalyticsErrorReporter func(error)

type ContactViewThroughLinkQueryHandler struct {
	linkResolver           PublicResumeLinkResolver
	resumeContactView      ResumeContactViewQuery
	resumeContactsDownload ResumeContactsDownloadQuery
	analytics              analyticsports.PublicInteractionEventRecorder
	reportAnalyticsError   ContactViewAnalyticsErrorReporter
}

func NewContactViewThroughLinkQueryHandler(linkResolver PublicResumeLinkResolver, resume ResumeContactViewQuery) *ContactViewThroughLinkQueryHandler {
	return &ContactViewThroughLinkQueryHandler{linkResolver: linkResolver, resumeContactView: resume}
}

func NewContactViewThroughLinkQueryHandlerWithAnalytics(
	linkResolver PublicResumeLinkResolver,
	resumeContactView ResumeContactViewQuery,
	resumeContactsDownload ResumeContactsDownloadQuery,
	analytics analyticsports.PublicInteractionEventRecorder,
	reportAnalyticsError ContactViewAnalyticsErrorReporter,
) *ContactViewThroughLinkQueryHandler {
	return &ContactViewThroughLinkQueryHandler{
		linkResolver:           linkResolver,
		resumeContactView:      resumeContactView,
		resumeContactsDownload: resumeContactsDownload,
		analytics:              analytics,
		reportAnalyticsError:   reportAnalyticsError,
	}
}

func (h *ContactViewThroughLinkQueryHandler) Handle(
	ctx context.Context,
	query ContactViewThroughLinkQuery,
) (ContactViewThroughLinkResult, error) {
	if h == nil || h.linkResolver == nil {
		return ContactViewThroughLinkResult{}, fmt.Errorf("contact view through link query handler dependencies are nil")
	}
	initiator := query.Initiator
	if initiator == "" && query.Contact != "" {
		initiator = ContactViewInitiatorContact
	}
	if initiator != ContactViewInitiatorContact && initiator != ContactViewInitiatorDownload {
		return ContactViewThroughLinkResult{}, fmt.Errorf("invalid contact view through link query")
	}
	if initiator == ContactViewInitiatorContact && h.resumeContactView == nil {
		return ContactViewThroughLinkResult{}, fmt.Errorf("contact view through link query handler dependencies are nil")
	}
	if initiator == ContactViewInitiatorDownload && h.resumeContactsDownload == nil {
		return ContactViewThroughLinkResult{}, fmt.Errorf("contact download through link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return ContactViewThroughLinkResult{}, fmt.Errorf("invalid contact view through link query")
	}
	if initiator == ContactViewInitiatorContact {
		if _, err := resumedomain.ParseContactKey(string(query.Contact)); err != nil {
			return ContactViewThroughLinkResult{}, fmt.Errorf("invalid contact view through link query")
		}
	} else if query.Contact != "" {
		return ContactViewThroughLinkResult{}, fmt.Errorf("invalid contact view through link query")
	}

	resumeID, err := h.linkResolver.Resolve(ctx, query.LinkID)
	if err != nil {
		return ContactViewThroughLinkResult{}, err
	}

	var result ContactViewThroughLinkResult
	switch initiator {
	case ContactViewInitiatorContact:
		view, err := h.resumeContactView.Handle(ctx, resumequeries.ContactViewQuery{
			ResumeID: resumeID,
			Contact:  query.Contact,
		})
		if err != nil {
			return ContactViewThroughLinkResult{}, err
		}
		result.Contact = &view.Contact
		result.Value = &view.Value
	case ContactViewInitiatorDownload:
		contacts, err := h.resumeContactsDownload.Handle(ctx, resumequeries.DownloadContactsQuery{ResumeID: resumeID})
		if err != nil {
			return ContactViewThroughLinkResult{}, err
		}
		downloadedContacts := make([]ContactViewThroughLinkContact, 0, 3)
		for _, contact := range []struct {
			key   resumedomain.ContactKey
			value *string
		}{
			{key: resumedomain.ContactKeyTelegram, value: contacts.Telegram},
			{key: resumedomain.ContactKeyEmail, value: contacts.Email},
			{key: resumedomain.ContactKeyPhone, value: contacts.Phone},
		} {
			if contact.value != nil && *contact.value != "" {
				downloadedContacts = append(downloadedContacts, ContactViewThroughLinkContact{
					Contact: contact.key,
					Value:   *contact.value,
				})
			}
		}
		result.Contacts = &downloadedContacts
	}
	if initiator == ContactViewInitiatorContact && !query.ViewID.IsZero() && h.analytics != nil {
		analyticsErr := h.analytics.RecordContactView(ctx, query.ViewID, query.LinkID, string(query.Contact))
		if analyticsErr != nil {
			if h.reportAnalyticsError != nil {
				h.reportAnalyticsError(analyticsErr)
			}
		}
	}
	return result, nil
}
