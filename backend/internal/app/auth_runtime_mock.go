//go:build devauth

package app

import (
	"errors"
	"net/http"

	"developers-cv/backend/internal/devauth"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
)

func newMockAuthenticator(reader userqueries.UserReader) (runtimeAuthenticator, error) {
	return devauth.NewMockAuthenticator(reader), nil
}

func isMockUnauthenticated(err error) bool {
	return errors.Is(err, devauth.ErrUnauthenticated)
}

func hasMockCredentials(request *http.Request) bool {
	return request != nil && request.Header.Get(devauth.MockUserIDHeader) != ""
}
