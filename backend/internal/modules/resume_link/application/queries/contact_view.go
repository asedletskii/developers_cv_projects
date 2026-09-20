package queries

import (
	"context"
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ContactViewThroughLinkQuery struct {
	LinkID    sharedvalueobjects.ID
	Initiator ContactViewInitiator
	Contact   contactkey.Key
	ViewID    sharedvalueobjects.ID
}

type ContactViewInitiator string

const (
	ContactViewInitiatorContact  ContactViewInitiator = "contact"
	ContactViewInitiatorDownload ContactViewInitiator = "download"
)

type ContactViewThroughLinkResult struct {
	Contact  *contactkey.Key
	Value    *string
	Contacts *[]ContactViewThroughLinkContact
}

type ContactViewThroughLinkContact struct {
	Contact contactkey.Key
	Value   string
}

type ContactViewAnalyticsErrorReporter func(error)

type ContactViewThroughLinkQueryHandler struct {
	linkResolver           PublicResumeLinkResolver
	resumeContactView      resumecontracts.ContactViewQueryHandler
	resumeContactsDownload resumecontracts.DownloadContactsQueryHandler
	analytics              analyticscontracts.PublicInteractionEventRecorder
	reportAnalyticsError   ContactViewAnalyticsErrorReporter
}

func NewContactViewThroughLinkQueryHandler(linkResolver PublicResumeLinkResolver, resume resumecontracts.ContactViewQueryHandler) *ContactViewThroughLinkQueryHandler {
	return &ContactViewThroughLinkQueryHandler{linkResolver: linkResolver, resumeContactView: resume}
}

func NewContactViewThroughLinkQueryHandlerWithAnalytics(
	linkResolver PublicResumeLinkResolver,
	resumeContactView resumecontracts.ContactViewQueryHandler,
	resumeContactsDownload resumecontracts.DownloadContactsQueryHandler,
	analytics analyticscontracts.PublicInteractionEventRecorder,
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
		if _, err := contactkey.Parse(string(query.Contact)); err != nil {
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
		view, err := h.resumeContactView.Handle(ctx, resumecontracts.ContactViewQuery{
			ResumeID: resumeID,
			Contact:  query.Contact,
		})
		if err != nil {
			return ContactViewThroughLinkResult{}, application.NormalizeProviderError(err)
		}
		if _, err := contactkey.Parse(string(view.Contact)); err != nil {
			return ContactViewThroughLinkResult{}, fmt.Errorf("map resume contact key: %w", err)
		}
		result.Contact = &view.Contact
		result.Value = &view.Value
	case ContactViewInitiatorDownload:
		contacts, err := h.resumeContactsDownload.Handle(ctx, resumecontracts.DownloadContactsQuery{ResumeID: resumeID})
		if err != nil {
			return ContactViewThroughLinkResult{}, application.NormalizeProviderError(err)
		}
		downloadedContacts := make([]ContactViewThroughLinkContact, 0, 3)
		for _, contact := range []struct {
			key   contactkey.Key
			value *string
		}{
			{key: contactkey.Telegram, value: contacts.Telegram},
			{key: contactkey.Email, value: contacts.Email},
			{key: contactkey.Phone, value: contacts.Phone},
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
