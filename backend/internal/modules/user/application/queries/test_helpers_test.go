package queries

import "context"

type profileReaderStub struct {
	profile ProfileView
	err     error
	query   GetProfileQuery
}

func (r *profileReaderStub) GetProfile(_ context.Context, query GetProfileQuery) (ProfileView, error) {
	r.query = query
	return r.profile, r.err
}

func stringPtr(value string) *string { return &value }
