package responses

type PublicContactViewResponse struct {
	Contact  *string                  `json:"contact,omitempty"`
	Value    *string                  `json:"value,omitempty"`
	Contacts *[]PublicContactResponse `json:"contacts,omitempty"`
}

type PublicContactResponse struct {
	Contact string `json:"contact"`
	Value   string `json:"value"`
}
