package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrPersonNotFound = errors.New("person not found")

type GetPersonQuery struct {
	ID   uuid.UUID
	Lang Language
}

type PersonReader interface {
	GetPerson(ctx context.Context, query GetPersonQuery) (Person, error)
}

type GetPersonHandler struct {
	reader PersonReader
}

func NewGetPersonHandler(reader PersonReader) *GetPersonHandler {
	return &GetPersonHandler{reader: reader}
}

func (h *GetPersonHandler) Handle(ctx context.Context, query GetPersonQuery) (Person, error) {
	if query.ID == uuid.Nil {
		return Person{}, fmt.Errorf("invalid person id")
	}
	if query.Lang != LanguageRU && query.Lang != LanguageEN {
		return Person{}, fmt.Errorf("unsupported language %q", query.Lang)
	}
	return h.reader.GetPerson(ctx, query)
}
