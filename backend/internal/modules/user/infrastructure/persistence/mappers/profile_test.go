package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/user/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestProfileToDBModels(t *testing.T) {
	userID := sharedvalueobjects.NewID()
	name, enName := "Дмитрий", "Dmitry"
	ru, err := domain.NewLocalizedProfile(&name, nil, nil, nil, nil)
	require.NoError(t, err)
	en, err := domain.NewLocalizedProfile(&enName, nil, nil, nil, nil)
	require.NoError(t, err)
	photoURL := "https://example.com/profile.jpg"
	profile, err := domain.NewProfile(userID, ru, en, &photoURL)
	require.NoError(t, err)
	now := time.Date(2026, 9, 30, 8, 0, 0, 0, time.UTC)

	model, rows, err := ProfileToDBModels(profile, now)
	require.NoError(t, err)
	require.Equal(t, userID.UUID(), model.UserID)
	require.Equal(t, photoURL, *model.PhotoURL)
	require.Equal(t, now, model.CreatedAt)
	require.Equal(t, now, model.UpdatedAt)
	require.Len(t, rows, 2)
	require.Equal(t, "ru", rows[0].Language)
	require.Equal(t, "en", rows[1].Language)
}

func TestProfileToDBModelsRejectsNil(t *testing.T) {
	_, _, err := ProfileToDBModels(nil, time.Now().UTC())
	require.Error(t, err)
}
