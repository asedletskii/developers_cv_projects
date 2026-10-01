//go:build integration

package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/repository"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestPostgresOwnedAndPublicResumeReadContracts(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeRepo := repository.NewResumeRepository(fixture.gorm)

	email, err := domain.NewEmail("private@example.com")
	require.NoError(t, err)
	contacts, err := domain.NewContactMethods(domain.ContactMethodsParams{Phone: "+33123456789", Email: &email})
	require.NoError(t, err)
	require.NoError(t, fixture.resume.ReplaceContacts(domain.ContactsParams{
		Contacts:       contacts,
		HiddenContacts: []domain.ContactKey{domain.ContactKeyEmail},
	}))

	gitLink, err := domain.NewExternalURL("https://github.com/example/project")
	require.NoError(t, err)
	screenshot, err := domain.NewExternalURL("https://example.com/project.png")
	require.NoError(t, err)
	require.NoError(t, fixture.resume.ReplaceProjects([]domain.ProjectParams{{
		Name:        "Public Project",
		Role:        stringPtr("Lead Developer"),
		GitLink:     &gitLink,
		Screenshots: []domain.ExternalURL{screenshot},
	}}))

	require.NoError(t, fixture.resume.Publish())
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, fixture.resume)
	}))

	reader := readers.NewResumeReader(fixture.gorm)
	owned, err := reader.GetOwned(ctx, queries.GetResumeQuery{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.Equal(t, "opened", owned.Status)
	require.NotNil(t, owned.Version)
	require.Equal(t, "private@example.com", *owned.Contacts.Email)
	require.Equal(t, "Lead Developer", *owned.Projects[0].Role)
	require.Equal(t, "https://github.com/example/project", *owned.Projects[0].GitLink)
	require.Len(t, owned.Projects[0].Screenshots, 1)

	publicReaderView, err := reader.GetPublic(ctx, queries.PublicResumeQuery{ID: fixture.resume.ID()})
	require.NoError(t, err)
	require.Equal(t, "opened", publicReaderView.Status)
	require.NotNil(t, publicReaderView.Contacts)
	require.Equal(t, "private@example.com", *publicReaderView.Contacts.Email)

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

	foreign := sharedvalueobjects.NewID()
	_, err = reader.GetOwned(ctx, queries.GetResumeQuery{ResumeID: fixture.resume.ID(), UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
}

func TestPostgresContactViewReturnsOnlyHiddenContact(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	ctx := fixture.ctx
	resumeRepo := repository.NewResumeRepository(fixture.gorm)

	email, err := domain.NewEmail("private@example.com")
	require.NoError(t, err)
	contacts, err := domain.NewContactMethods(domain.ContactMethodsParams{Phone: "+33123456789", Email: &email})
	require.NoError(t, err)
	require.NoError(t, fixture.resume.ReplaceContacts(domain.ContactsParams{
		Contacts:       contacts,
		HiddenContacts: []domain.ContactKey{domain.ContactKeyEmail},
	}))
	require.NoError(t, fixture.resume.Publish())
	require.NoError(t, fixture.uow.Do(ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, fixture.resume)
	}))

	reader := readers.NewResumeReader(fixture.gorm)
	view, err := reader.GetContactForView(ctx, queries.ContactViewQuery{
		ResumeID: fixture.resume.ID(),
		Contact:  domain.ContactKeyEmail,
	})
	require.NoError(t, err)
	require.Equal(t, domain.ContactKeyEmail, view.Contact)
	require.Equal(t, "private@example.com", view.Value)

	_, err = reader.GetContactForView(ctx, queries.ContactViewQuery{
		ResumeID: fixture.resume.ID(),
		Contact:  domain.ContactKeyPhone,
	})
	require.ErrorIs(t, err, queries.ErrContactNotHidden)
}

func TestPostgresPublicResumeRequiresOpenedStatus(t *testing.T) {
	fixture := newPhase6PostgresFixture(t)
	reader := readers.NewResumeReader(fixture.gorm)

	_, err := reader.GetPublic(fixture.ctx, queries.PublicResumeQuery{ID: fixture.resume.ID()})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)

	resumeRepo := repository.NewResumeRepository(fixture.gorm)
	require.NoError(t, fixture.resume.Publish())
	require.NoError(t, fixture.uow.Do(fixture.ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, fixture.resume)
	}))

	_, err = reader.GetPublic(fixture.ctx, queries.PublicResumeQuery{ID: fixture.resume.ID()})
	require.NoError(t, err)

	opened, err := resumeRepo.Get(fixture.ctx, ports.GetResumeParams{ResumeID: fixture.resume.ID(), UserID: fixture.owner})
	require.NoError(t, err)
	require.NoError(t, opened.Close())
	require.NoError(t, fixture.uow.Do(fixture.ctx, func(tx context.Context) error {
		return resumeRepo.Save(tx, opened)
	}))

	_, err = reader.GetPublic(fixture.ctx, queries.PublicResumeQuery{ID: fixture.resume.ID()})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
}
