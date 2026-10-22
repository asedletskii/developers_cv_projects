package mappers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestResumeLinkRoundTrip(t *testing.T) {
	id := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	link, err := domain.RestoreResumeLink(domain.ResumeLinkSnapshot{
		ID: id, ResumeID: resumeID, Name: "  HH  ", Version: 3, Status: domain.ResumeLinkStatusInactive,
	})
	require.NoError(t, err)

	createdAt := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Minute)
	model, err := ResumeLinkToDB(link, createdAt, updatedAt)
	require.NoError(t, err)
	require.Equal(t, id.UUID(), model.ID)
	require.Equal(t, resumeID.UUID(), model.ResumeID)
	require.Equal(t, "  HH  ", model.Name)
	require.Equal(t, string(domain.ResumeLinkStatusInactive), model.Status)
	require.Equal(t, int64(3), model.Version)
	require.Equal(t, createdAt, model.CreatedAt)
	require.Equal(t, updatedAt, model.UpdatedAt)

	mapped, err := ResumeLinkToDomain(model)
	require.NoError(t, err)
	require.Equal(t, id, mapped.ID())
	require.Equal(t, resumeID, mapped.ResumeID())
	require.Equal(t, "  HH  ", mapped.Name())
	require.Equal(t, domain.ResumeLinkStatusInactive, mapped.Status())
	require.Equal(t, int64(3), mapped.Version())
}

func TestResumeLinkToDomainRejectsInvalidStatus(t *testing.T) {
	model := gormmodel.ResumeLink{ID: uuid.New(), ResumeID: uuid.New(), Name: "HH", Status: "broken", Version: 1}
	_, err := ResumeLinkToDomain(model)
	require.ErrorIs(t, err, domain.ErrInvalidResumeLink)
}
