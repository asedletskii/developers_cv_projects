package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/user/application"
	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

const MockUserIDHeader = "X-Mock-User-ID"

var ErrUnauthenticated = errors.New("unauthenticated")

type MockAuthenticator struct {
	repository application.UserRepository
}

func NewMockAuthenticator(repository application.UserRepository) *MockAuthenticator {
	return &MockAuthenticator{repository: repository}
}

func (a *MockAuthenticator) Authenticate(ctx context.Context, request *http.Request) (uuid.UUID, error) {
	if a == nil || a.repository == nil {
		return uuid.Nil, fmt.Errorf("mock authenticator dependencies are nil")
	}

	rawID := strings.TrimSpace(request.Header.Get(MockUserIDHeader))
	if rawID == "" {
		return uuid.Nil, ErrUnauthenticated
	}

	id, err := uuid.Parse(rawID)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, ErrUnauthenticated
	}

	userID, err := valueobjects.IDFromUUID(id)
	if err != nil {
		return uuid.Nil, ErrUnauthenticated
	}
	if _, err := a.repository.GetByID(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return uuid.Nil, ErrUnauthenticated
		}
		return uuid.Nil, fmt.Errorf("validate mock user: %w", err)
	}
	return id, nil
}
