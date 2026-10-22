package application

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume_link/domain"
)

func TestNormalizeProviderError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
	}{
		{name: "resume not found", err: resumecontracts.ErrResumeNotFound, target: ErrResumeNotFound},
		{name: "contact not found", err: resumecontracts.ErrContactNotFound, target: ErrContactNotFound},
		{name: "contact not hidden", err: resumecontracts.ErrContactNotHidden, target: ErrContactNotHidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := fmt.Errorf("provider call: %w", tt.err)

			got := NormalizeProviderError(wrapped)

			require.ErrorIs(t, got, tt.target)
			require.ErrorIs(t, got, tt.err)
		})
	}
}

func TestNormalizeProviderErrorLeavesUnknownErrorUntouched(t *testing.T) {
	want := errors.New("database unavailable")

	got := NormalizeProviderError(want)

	require.ErrorIs(t, got, want)
}

func TestNormalizeProviderErrorHandlesNil(t *testing.T) {
	require.NoError(t, NormalizeProviderError(nil))
}

func TestNormalizeLinkNotFoundError(t *testing.T) {
	wrapped := fmt.Errorf("get link: %w", domain.ErrEntityNotFound)

	got := NormalizeLinkNotFoundError(wrapped)

	require.ErrorIs(t, got, ErrResumeLinkNotFound)
	require.ErrorIs(t, got, domain.ErrEntityNotFound)
}

func TestNormalizeLinkNotFoundErrorLeavesOtherErrorsUntouched(t *testing.T) {
	want := errors.New("database unavailable")

	got := NormalizeLinkNotFoundError(want)

	require.ErrorIs(t, got, want)
}
