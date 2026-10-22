//go:build devauth

package devauth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"net/http"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/user/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const MockUserIDHeader = "X-Mock-User-ID"

var ErrUnauthenticated = errors.New("unauthenticated")

type MockAuthenticator struct {
	reader queries.UserReader
}

func NewMockAuthenticator(reader queries.UserReader) *MockAuthenticator {
	return &MockAuthenticator{reader: reader}
}

func (a *MockAuthenticator) Authenticate(ctx context.Context, request *http.Request) (uuid.UUID, error) {
	if a == nil || a.reader == nil {
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

	userID, err := sharedvalueobjects.IDFromUUID(id)
	if err != nil {
		return uuid.Nil, ErrUnauthenticated
	}
	if _, err := a.reader.GetUser(ctx, queries.GetUserQuery{ID: userID}); err != nil {
		if errors.Is(err, queries.ErrUserNotFound) {
			return uuid.Nil, ErrUnauthenticated
		}
		return uuid.Nil, fmt.Errorf("validate mock user: %w", err)
	}
	return id, nil
}
