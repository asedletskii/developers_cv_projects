package requests

type PublicContactViewRequest struct {
	Initiator      string  `json:"initiator"`
	Contact        *string `json:"contact,omitempty"`
	TurnstileToken string  `json:"turnstileToken"`
	ViewID         string  `json:"viewId,omitempty"`
}
