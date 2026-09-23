package queries

import (
	"context"
	"errors"
	"fmt"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"github.com/google/uuid"
)

var ErrResumeNotFound = errors.New("resume not found")

type GetResumeQuery struct {
	UserID *useridentity.ID
	ID     uuid.UUID
	Lang   Language
}

type GetResumeReader interface {
	GetResume(ctx context.Context, query GetResumeQuery) (Resume, error)
}

type GetResumeHandler struct {
	reader GetResumeReader
}

func NewGetResumeHandler(reader GetResumeReader) *GetResumeHandler {
	return &GetResumeHandler{reader: reader}
}

func (h *GetResumeHandler) Handle(ctx context.Context, query GetResumeQuery) (Resume, error) {
	if query.UserID != nil && query.UserID.IsZero() {
		return Resume{}, fmt.Errorf("invalid user id")
	}
	if query.ID == uuid.Nil {
		return Resume{}, fmt.Errorf("invalid resume id")
	}
	if query.Lang != LanguageRU && query.Lang != LanguageEN {
		return Resume{}, fmt.Errorf("unsupported language %q", query.Lang)
	}
	return h.reader.GetResume(ctx, query)
}
