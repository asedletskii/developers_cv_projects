package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type stubListResumesReader struct {
	items  []ResumeListItem
	err    error
	userID sharedvalueobjects.ID
}

func (s *stubListResumesReader) ListResumes(_ context.Context, query ListResumesQuery) ([]ResumeListItem, error) {
	s.userID = query.UserID
	return s.items, s.err
}

func TestListResumesUsesUserID(t *testing.T) {
	reader := &stubListResumesReader{items: []ResumeListItem{{ID: "resume-a", Title: "Senior Go Developer", Status: "draft", Version: 3}}}
	h := NewListResumesQueryHandler(reader)
	userID := sharedvalueobjects.NewID()

	items, err := h.Handle(context.Background(), ListResumesQuery{UserID: userID})

	require.NoError(t, err)
	require.Equal(t, userID.UUID(), reader.userID.UUID())
	require.Equal(t, reader.items, items)
}

func TestListResumesRejectsEmptyUserID(t *testing.T) {
	reader := &stubListResumesReader{}
	h := NewListResumesQueryHandler(reader)

	_, err := h.Handle(context.Background(), ListResumesQuery{})

	require.Error(t, err)
	require.True(t, reader.userID.IsZero())
}
