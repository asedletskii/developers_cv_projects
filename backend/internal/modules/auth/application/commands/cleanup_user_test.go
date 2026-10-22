package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type userDataCleanupRepositoryStub struct {
	ids []sharedvalueobjects.ID
	err error
}

func (s *userDataCleanupRepositoryStub) DeleteByUserID(_ context.Context, userID sharedvalueobjects.ID) error {
	s.ids = append(s.ids, userID)
	return s.err
}

func TestDeleteUserDataHandlerDeletesByUserID(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	repository := &userDataCleanupRepositoryStub{}
	handler := NewDeleteUserDataHandler(repository)

	require.NoError(t, handler.DeleteByUserID(context.Background(), userID))
	require.Equal(t, []sharedvalueobjects.ID{userID}, repository.ids)
}

func TestDeleteUserDataHandlerRejectsEmptyUserID(t *testing.T) {
	repository := &userDataCleanupRepositoryStub{}
	handler := NewDeleteUserDataHandler(repository)

	require.Error(t, handler.DeleteByUserID(context.Background(), sharedvalueobjects.ID{}))
	require.Empty(t, repository.ids)
}

func TestDeleteUserDataHandlerPropagatesRepositoryError(t *testing.T) {
	repoErr := errors.New("delete failed")
	handler := NewDeleteUserDataHandler(&userDataCleanupRepositoryStub{err: repoErr})

	err := handler.DeleteByUserID(context.Background(), sharedvalueobjects.NewID())

	require.ErrorIs(t, err, repoErr)
}
