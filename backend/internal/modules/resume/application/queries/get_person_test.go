package queries

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakePersonReader struct {
	query  GetPersonQuery
	person Person
	err    error
}

func (f *fakePersonReader) GetPerson(_ context.Context, query GetPersonQuery) (Person, error) {
	f.query = query
	return f.person, f.err
}

func TestGetPersonHandlerPassesQueryObjectToReader(t *testing.T) {
	reader := &fakePersonReader{person: Person{ID: "person-1"}}
	handler := NewGetPersonHandler(reader)
	id := uuid.New()
	query := GetPersonQuery{ID: id, Lang: LanguageEN}

	person, err := handler.Handle(context.Background(), query)

	require.NoError(t, err)
	require.Equal(t, reader.person, person)
	require.Equal(t, query, reader.query)
}

func TestGetPersonHandlerRejectsInvalidQuery(t *testing.T) {
	reader := &fakePersonReader{}
	handler := NewGetPersonHandler(reader)

	_, err := handler.Handle(context.Background(), GetPersonQuery{Lang: LanguageEN})
	require.Error(t, err)

	_, err = handler.Handle(context.Background(), GetPersonQuery{ID: uuid.New(), Lang: "de"})
	require.Error(t, err)
}
