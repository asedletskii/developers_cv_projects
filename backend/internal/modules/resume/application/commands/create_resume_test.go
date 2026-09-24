package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type fakeResumeRepository struct {
	saved *domain.Resume
	err   error
	calls int
}

func (f *fakeResumeRepository) GetByID(context.Context, useridentity.ID, valueobjects.ID) (*domain.Resume, error) {
	return nil, nil
}
func (f *fakeResumeRepository) Save(_ context.Context, resume *domain.Resume) error {
	f.calls++
	f.saved = resume
	return f.err
}

func (f *fakeResumeRepository) UpdateTitle(context.Context, useridentity.ID, valueobjects.ID, string) error {
	return f.err
}

type fakeUnitOfWork struct {
	err   error
	calls int
	ctx   context.Context
}

func (f *fakeUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	f.calls++
	f.ctx = ctx
	if f.err != nil {
		return f.err
	}
	return fn(ctx)
}

func validCreateCommand(t *testing.T) CreateResumeCommand {
	date, err := valueobjects.ParseDate("1998-04-12")
	require.NoError(t, err)
	photo, err := valueobjects.NewAssetPath("/assets/photo.svg")
	require.NoError(t, err)
	first, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Привет",
			EN: "Hello",
		},
	)
	require.NoError(t, err)
	second, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Я",
			EN: "I am",
		},
	)
	require.NoError(t, err)
	highlight, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Тест",
			EN: "Test",
		},
	)
	require.NoError(t, err)
	third, err := valueobjects.NewLocalizedText(
		valueobjects.LocalizedTextParams{
			RU: "Разработчик",
			EN: "Developer",
		},
	)
	require.NoError(t, err)
	greeting := valueobjects.NewGreeting(
		valueobjects.GreetingParams{
			First:           first,
			Second:          second,
			SecondHighlight: highlight,
			Third:           third,
		},
	)
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: "+79999999999"})
	require.NoError(t, err)
	userID := useridentity.NewID()
	return CreateResumeCommand{
		UserID:      userID,
		Name:        "Test",
		Surname:     "User",
		Position:    "Developer",
		Greeting:    greeting,
		PhotoPath:   photo,
		DateOfBirth: date,
		Location:    "Remote",
		Contacts:    contacts,
	}
}

func TestCreateResumeHandlerValidCommand(t *testing.T) {
	repo := &fakeResumeRepository{}
	uow := &fakeUnitOfWork{}
	h := NewCreateResumeHandler(repo, uow)

	command := validCreateCommand(t)
	id, err := h.Handle(context.Background(), command)

	require.NoError(t, err)
	require.Equal(t, 1, uow.calls)
	require.Equal(t, 1, repo.calls)
	require.NotNil(t, repo.saved)
	require.False(t, repo.saved.ID().IsZero())
	require.Equal(t, repo.saved.ID(), id)
	require.Equal(t, "Test", repo.saved.Name())
	require.Equal(t, "Developer", repo.saved.Title())
	require.Equal(t, command.UserID.UUID(), repo.saved.UserID().UUID())
}

func TestCreateResumeHandlerDomainValidationDoesNotSave(t *testing.T) {
	repo := &fakeResumeRepository{}
	uow := &fakeUnitOfWork{}
	h := NewCreateResumeHandler(repo, uow)
	command := validCreateCommand(t)
	command.Name = ""

	_, err := h.Handle(context.Background(), command)

	require.Error(t, err)
	require.Equal(t, 0, repo.calls)
}

func TestCreateResumeHandlerRepositoryError(t *testing.T) {
	repoErr := errors.New("repository failed")
	repo := &fakeResumeRepository{err: repoErr}
	uow := &fakeUnitOfWork{}
	h := NewCreateResumeHandler(repo, uow)

	_, err := h.Handle(context.Background(), validCreateCommand(t))

	require.ErrorIs(t, err, repoErr)
	require.Equal(t, 1, repo.calls)
}

func TestCreateResumeHandlerUnitOfWorkError(t *testing.T) {
	uowErr := errors.New("uow failed")
	repo := &fakeResumeRepository{}
	uow := &fakeUnitOfWork{err: uowErr}
	h := NewCreateResumeHandler(repo, uow)

	_, err := h.Handle(context.Background(), validCreateCommand(t))

	require.ErrorIs(t, err, uowErr)
	require.Equal(t, 0, repo.calls)
}

func TestCreateResumeHandlerPassesAggregateToRepository(t *testing.T) {
	repo := &fakeResumeRepository{}
	uow := &fakeUnitOfWork{}
	h := NewCreateResumeHandler(repo, uow)
	command := validCreateCommand(t)

	_, err := h.Handle(context.Background(), command)
	require.NoError(t, err)
	require.NotNil(t, repo.saved)
	require.Equal(t, "Test", repo.saved.Name())
	require.Equal(t, command.UserID.UUID(), repo.saved.UserID().UUID())
}

var _ ResumeRepository = (*fakeResumeRepository)(nil)
var _ UnitOfWork = (*fakeUnitOfWork)(nil)
