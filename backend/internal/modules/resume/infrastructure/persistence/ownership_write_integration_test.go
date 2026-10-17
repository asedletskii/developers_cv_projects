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
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
	"developers-cv/backend/internal/testutil/pgtest"
)

func TestPostgresResumeRepositoryEnforcesOwnership(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := sharedvalueobjects.NewID()
	foreign := sharedvalueobjects.NewID()
	insertResumeTestUser(t, ctx, dsn, owner)
	insertResumeTestUser(t, ctx, dsn, foreign)

	resume := newTestResume(t, owner)
	require.NoError(t, resume.ReplaceOverview(domain.ResumeOverviewParams{
		Greeting: mustGreeting(t, domain.GreetingParams{}),
		Sections: []domain.SectionSettings{},
	}))
	require.NoError(t, resume.ReplaceProjects([]domain.ProjectParams{{Name: "Resume Platform"}}))

	resumeRepository := repository.NewResumeRepository(gormDB)
	uow := database.NewUnitOfWork(gormDB)
	require.NoError(t, uow.Do(ctx, func(tx context.Context) error {
		return resumeRepository.Create(tx, resume)
	}))

	_, err = resumeRepository.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: foreign})
	require.ErrorIs(t, err, domain.ErrEntityNotFound)

	_, err = resumeRepository.Get(ctx, ports.GetResumeParams{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)

	resumeReader := readers.NewResumeReader(gormDB)
	owned, err := resumeReader.OwnsResume(ctx, queries.GetResumeQuery{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.True(t, owned)
	overview, err := resumeReader.GetOwnedOverview(ctx, queries.GetResumeQuery{ResumeID: resume.ID(), UserID: owner})
	require.NoError(t, err)
	require.Empty(t, overview.Sections)
	_, err = resumeReader.GetOwnedOverview(ctx, queries.GetResumeQuery{ResumeID: resume.ID(), UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	owned, err = resumeReader.OwnsResume(ctx, queries.GetResumeQuery{ResumeID: resume.ID(), UserID: foreign})
	require.NoError(t, err)
	require.False(t, owned)
	owned, err = resumeReader.OwnsResume(ctx, queries.GetResumeQuery{ResumeID: sharedvalueobjects.NewID(), UserID: owner})
	require.NoError(t, err)
	require.False(t, owned)
}
