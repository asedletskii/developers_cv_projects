//go:build !devauth

package app

import (
	"errors"
	"net/http"

	userqueries "developers-cv/backend/internal/modules/user/application/queries"
)

func newMockAuthenticator(_ userqueries.UserReader) (runtimeAuthenticator, error) {
	return nil, errors.New("mock authentication is not compiled into this binary; rebuild with -tags=devauth for development and tests")
}

func isMockUnauthenticated(error) bool {
	return false
}

func hasMockCredentials(*http.Request) bool {
	return false
}
