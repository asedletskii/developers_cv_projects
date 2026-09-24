package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type stubListResumesReader struct {
	items  []ResumeListItem
	err    error
	userID useridentity.ID
}

func (s *stubListResumesReader) ListResumes(_ context.Context, query GetListResumesQuery) ([]ResumeListItem, error) {
	s.userID = query.UserID
	return s.items, s.err
}

func TestListResumesUsesUserID(t *testing.T) {
	reader := &stubListResumesReader{items: []ResumeListItem{{ID: "resume-a", Title: "Senior Go Developer"}}}
	h := NewListResumesHandler(reader)
	userID := useridentity.NewID()

	items, err := h.Handle(context.Background(), GetListResumesQuery{UserID: userID})

	require.NoError(t, err)
	require.Equal(t, userID.UUID(), reader.userID.UUID())
	require.Equal(t, reader.items, items)
}

func TestListResumesRejectsEmptyUserID(t *testing.T) {
	reader := &stubListResumesReader{}
	h := NewListResumesHandler(reader)

	_, err := h.Handle(context.Background(), GetListResumesQuery{})

	require.Error(t, err)
	require.True(t, reader.userID.IsZero())
}
