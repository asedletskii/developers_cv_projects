package mock

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

func TestPersonReaderReturnsLocalizedPerson(t *testing.T) {
	reader := NewPersonReader()

	person, err := reader.GetPerson(context.Background(), queries.GetPersonQuery{
		ID:   uuid.MustParse(PersonOneID),
		Lang: queries.LanguageEN,
	})

	require.NoError(t, err)
	require.Equal(t, PersonOneID, person.ID)
	require.Equal(t, "Hello!", person.Greeting.First)
	require.Equal(t, "Anton Sedletskii", person.Greeting.SecondHighlight)
	require.NotNil(t, person.ContactsText)
	require.Equal(t, "Tell me about your project", *person.ContactsText.Title)
	require.NotNil(t, person.UpdatedAt)
	require.NotNil(t, person.Availability)
	require.Equal(t, "Looking for a Senior Frontend role, remote or hybrid", *person.Availability.Note)
}

func TestPersonReaderReturnsNotFound(t *testing.T) {
	reader := NewPersonReader()

	_, err := reader.GetPerson(context.Background(), queries.GetPersonQuery{
		ID:   uuid.New(),
		Lang: queries.LanguageRU,
	})

	require.ErrorIs(t, err, queries.ErrPersonNotFound)
}

func TestPersonReaderKeepsEmptyArrays(t *testing.T) {
	reader := NewPersonReader()

	person, err := reader.GetPerson(context.Background(), queries.GetPersonQuery{
		ID:   uuid.MustParse(PersonTwoID),
		Lang: queries.LanguageRU,
	})

	require.NoError(t, err)
	require.NotNil(t, person.HiddenContacts)
	require.Empty(t, person.HiddenContacts)
	require.NotNil(t, person.Courses)
	require.Empty(t, person.Courses)
	require.NotNil(t, person.Activities)
	require.Empty(t, person.Activities)
}
