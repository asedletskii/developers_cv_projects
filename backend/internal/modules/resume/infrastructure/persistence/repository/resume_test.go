package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
)

func TestResumeRepositoryGetRejectsNilReceiver(t *testing.T) {
	var repository *ResumeRepository

	_, err := repository.Get(context.Background(), ports.GetResumeParams{})

	require.ErrorContains(t, err, "repository dependencies are nil")
}
