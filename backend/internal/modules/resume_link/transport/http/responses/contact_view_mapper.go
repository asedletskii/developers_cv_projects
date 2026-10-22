package responses

import "developers-cv/backend/internal/modules/resume_link/application/queries"

func NewPublicContactViewResponse(result queries.ContactViewThroughLinkResult) PublicContactViewResponse {
	response := PublicContactViewResponse{}
	if result.Contact != nil {
		contact := string(*result.Contact)
		response.Contact = &contact
	}
	if result.Value != nil {
		value := *result.Value
		response.Value = &value
	}
	if result.Contacts != nil {
		contacts := make([]PublicContactResponse, len(*result.Contacts))
		for i, contact := range *result.Contacts {
			contacts[i] = PublicContactResponse{Contact: string(contact.Contact), Value: contact.Value}
		}
		response.Contacts = &contacts
	}
	return response
}
