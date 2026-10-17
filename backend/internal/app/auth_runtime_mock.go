//go:build devauth

package app

import (
	"errors"
	"net/http"

	authmock "developers-cv/backend/internal/modules/auth/infrastructure/mock"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
)

func newMockAuthenticator(reader userqueries.UserReader) (runtimeAuthenticator, error) {
	return authmock.NewMockAuthenticator(reader), nil
}

func isMockUnauthenticated(err error) bool {
	return errors.Is(err, authmock.ErrUnauthenticated)
}

func hasMockCredentials(request *http.Request) bool {
	return request != nil && request.Header.Get(authmock.MockUserIDHeader) != ""
}
