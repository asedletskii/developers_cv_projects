package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/application/queries"
	contactkey "developers-cv/backend/internal/shared/domain/contact"
)

func TestNewPublicContactViewResponseMapsReveal(t *testing.T) {
	contact := contactkey.Email
	value := "person@example.com"

	payload, err := json.Marshal(NewPublicContactViewResponse(queries.ContactViewThroughLinkResult{
		Contact: &contact,
		Value:   &value,
	}))
	require.NoError(t, err)
	require.JSONEq(t, `{"contact":"email","value":"person@example.com"}`, string(payload))
}

func TestNewPublicContactViewResponseMapsDownload(t *testing.T) {
	contacts := []queries.ContactViewThroughLinkContact{
		{Contact: contactkey.Email, Value: "person@example.com"},
		{Contact: contactkey.Phone, Value: "+33123456789"},
	}

	payload, err := json.Marshal(NewPublicContactViewResponse(queries.ContactViewThroughLinkResult{
		Contacts: &contacts,
	}))
	require.NoError(t, err)
	require.JSONEq(t, `{"contacts":[{"contact":"email","value":"person@example.com"},{"contact":"phone","value":"+33123456789"}]}`, string(payload))
}

func TestNewPublicContactViewResponseKeepsEmptyDownloadArray(t *testing.T) {
	contacts := []queries.ContactViewThroughLinkContact{}

	payload, err := json.Marshal(NewPublicContactViewResponse(queries.ContactViewThroughLinkResult{
		Contacts: &contacts,
	}))
	require.NoError(t, err)
	require.JSONEq(t, `{"contacts":[]}`, string(payload))
}

func TestNewPublicContactViewResponseOmitsNilFields(t *testing.T) {
	payload, err := json.Marshal(NewPublicContactViewResponse(queries.ContactViewThroughLinkResult{}))
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(payload))
}
