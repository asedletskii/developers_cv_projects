package queries

type ResumeContactsView struct {
	ID             string                  `json:"id"`
	ResumeID       string                  `json:"resumeId"`
	Phone          *string                 `json:"phone,omitempty"`
	Telegram       *string                 `json:"telegram,omitempty"`
	Email          *string                 `json:"email,omitempty"`
	HiddenContacts []string                `json:"hiddenContacts"`
	ContactsText   *ResumeContactsTextView `json:"contactsText,omitempty"`
}

type ResumeContactsTextView struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}
