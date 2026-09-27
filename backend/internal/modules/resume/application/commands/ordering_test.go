package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeReorderLanguagesRepository struct {
	called bool
	cmd    ports.ReorderLanguagesRepositoryCommand
}

func (r *fakeReorderLanguagesRepository) ReorderLanguages(_ context.Context, cmd ports.ReorderLanguagesRepositoryCommand) error {
	r.called = true
	r.cmd = cmd
	return nil
}

type fakeReorderTechnologiesRepository struct {
	called bool
	cmd    ports.ReorderTechnologiesRepositoryCommand
}

func (r *fakeReorderTechnologiesRepository) ReorderTechnologies(_ context.Context, cmd ports.ReorderTechnologiesRepositoryCommand) error {
	r.called = true
	r.cmd = cmd
	return nil
}

func TestReorderLanguagesHandlerPassesFullCollectionOrder(t *testing.T) {
	repository := &fakeReorderLanguagesRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReorderLanguagesHandler(repository, uow)
	resumeID := valueobjects.NewID()
	userID := useridentity.NewID()
	ids := []valueobjects.ID{valueobjects.NewID(), valueobjects.NewID()}

	err := handler.Handle(context.Background(), ReorderLanguagesCommand{ResumeID: resumeID, UserID: userID, LanguageIDs: ids})

	require.NoError(t, err)
	require.True(t, repository.called)
	require.Equal(t, resumeID, repository.cmd.ResumeID)
	require.Equal(t, userID, repository.cmd.UserID)
	require.Equal(t, ids, repository.cmd.LanguageIDs)
}

func TestReorderTechnologiesHandlerPassesFullCollectionOrder(t *testing.T) {
	repository := &fakeReorderTechnologiesRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReorderTechnologiesHandler(repository, uow)
	resumeID := valueobjects.NewID()
	userID := useridentity.NewID()
	ids := []valueobjects.ID{valueobjects.NewID(), valueobjects.NewID()}

	err := handler.Handle(context.Background(), ReorderTechnologiesCommand{ResumeID: resumeID, UserID: userID, TechnologyIDs: ids})

	require.NoError(t, err)
	require.True(t, repository.called)
	require.Equal(t, resumeID, repository.cmd.ResumeID)
	require.Equal(t, userID, repository.cmd.UserID)
	require.Equal(t, ids, repository.cmd.TechnologyIDs)
}

func TestReorderLanguagesHandlerRejectsEmptyOrder(t *testing.T) {
	repository := &fakeReorderLanguagesRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReorderLanguagesHandler(repository, uow)

	err := handler.Handle(context.Background(), ReorderLanguagesCommand{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID()})

	require.Error(t, err)
	require.False(t, repository.called)
}

func TestReorderTechnologiesHandlerRejectsEmptyOrder(t *testing.T) {
	repository := &fakeReorderTechnologiesRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReorderTechnologiesHandler(repository, uow)

	err := handler.Handle(context.Background(), ReorderTechnologiesCommand{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID()})

	require.Error(t, err)
	require.False(t, repository.called)
}

func TestReorderLanguagesHandlerRejectsDuplicateIDs(t *testing.T) {
	repository := &fakeReorderLanguagesRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReorderLanguagesHandler(repository, uow)
	id := valueobjects.NewID()

	err := handler.Handle(context.Background(), ReorderLanguagesCommand{
		ResumeID:    valueobjects.NewID(),
		UserID:      useridentity.NewID(),
		LanguageIDs: []valueobjects.ID{id, id},
	})

	require.Error(t, err)
	require.False(t, repository.called)
}

func TestReorderTechnologiesHandlerRejectsDuplicateIDs(t *testing.T) {
	repository := &fakeReorderTechnologiesRepository{}
	uow := &fakeResumeUnitOfWork{}
	handler := NewReorderTechnologiesHandler(repository, uow)
	id := valueobjects.NewID()

	err := handler.Handle(context.Background(), ReorderTechnologiesCommand{
		ResumeID:      valueobjects.NewID(),
		UserID:        useridentity.NewID(),
		TechnologyIDs: []valueobjects.ID{id, id},
	})

	require.Error(t, err)
	require.False(t, repository.called)
}
