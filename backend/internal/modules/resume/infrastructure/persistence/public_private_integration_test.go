//go:build integration

package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestPostgresOwnedAndPublicResumeReadContracts(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeRepo := repository.NewResumeRepository(fixture.gorm)

	email, err := valueobjects.NewEmail("private@example.com")
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: "+33123456789", Email: &email})
	require.NoError(t, err)
	require.NoError(t, fixture.resume.ReplaceContacts(domain.ContactsParams{
		Contacts:       contacts,
		HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail},
	}))

	gitLink, err := valueobjects.NewExternalURL("https://github.com/example/project")
	require.NoError(t, err)
	screenshot, err := valueobjects.NewExternalURL("https://example.com/project.png")
	require.NoError(t, err)
	require.NoError(t, fixture.resume.ReplaceProjects([]domain.ProjectParams{{
		Name:        "Public Project",
		Role:        stringPtr("Lead Developer"),
		GitLink:     &gitLink,
		Screenshots: []valueobjects.ExternalURL{screenshot},
	}}))

	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, fixture.resume)
	}))

	reader := readers.NewResumeReader(fixture.gorm)
	owned, err := reader.GetOwned(ctx, queries.GetResumeQuery{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.NotNil(t, owned.Version)
	require.Equal(t, "private@example.com", *owned.Contacts.Email)
	require.Equal(t, "Lead Developer", *owned.Projects[0].Role)
	require.Equal(t, "https://github.com/example/project", *owned.Projects[0].GitLink)
	require.Len(t, owned.Projects[0].Screenshots, 1)

	publicHandler := queries.NewPublicResumeQueryHandler(reader)
	public, err := publicHandler.Handle(ctx, queries.PublicResumeQuery{ID: fixture.resume.ID()})
	require.NoError(t, err)
	require.Nil(t, public.Version)
	require.NotNil(t, public.Contacts)
	require.Nil(t, public.Contacts.Email)
	require.Equal(t, "+33123456789", *public.Contacts.Phone)
	require.Equal(t, "Public Project", public.Projects[0].Name)
	require.Nil(t, public.Projects[0].Role)
	require.Nil(t, public.Projects[0].GitLink)
	require.Empty(t, public.Projects[0].Screenshots)
	require.Empty(t, public.Projects[0].Contributors)

	foreign := useridentity.NewID()
	_, err = reader.GetOwned(ctx, queries.GetResumeQuery{ResumeID: fixture.resume.ID(), UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
}

func TestPostgresContactViewReturnsOnlyHiddenContact(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeRepo := repository.NewResumeRepository(fixture.gorm)

	email, err := valueobjects.NewEmail("private@example.com")
	require.NoError(t, err)
	contacts, err := entities.NewContacts(entities.ContactsParams{Phone: "+33123456789", Email: &email})
	require.NoError(t, err)
	require.NoError(t, fixture.resume.ReplaceContacts(domain.ContactsParams{
		Contacts:       contacts,
		HiddenContacts: []enums.ContactKey{enums.ContactKeyEmail},
	}))
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, fixture.resume)
	}))

	reader := readers.NewResumeReader(fixture.gorm)
	view, err := reader.GetContactForView(ctx, queries.ContactViewQuery{
		ResumeID: fixture.resume.ID(),
		Contact:  enums.ContactKeyEmail,
	})
	require.NoError(t, err)
	require.Equal(t, enums.ContactKeyEmail, view.Contact)
	require.Equal(t, "private@example.com", view.Value)

	_, err = reader.GetContactForView(ctx, queries.ContactViewQuery{
		ResumeID: fixture.resume.ID(),
		Contact:  enums.ContactKeyPhone,
	})
	require.ErrorIs(t, err, queries.ErrContactNotHidden)
}
