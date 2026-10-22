package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestReplaceEducationsHandlerRejectsStaleVersion(t *testing.T) {
	r := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 2)
	repo := &commandRepositoryStub{resume: r}
	h := NewReplaceEducationsHandler(repo, &unitOfWorkStub{})
	err := h.Handle(context.Background(), ReplaceEducationsCommand{ResumeID: r.ID(), UserID: r.UserID(), Version: 1, Input: EducationsInput{}})
	require.ErrorIs(t, err, domain.ErrConcurrencyConflict)
	require.False(t, repo.saveCalled)
}
func TestReplaceEducationsHandlerUsesAggregateRepository(t *testing.T) {
	r := newCommandTestResume(t, sharedvalueobjects.NewID(), "Resume", 1)
	repo := &commandRepositoryStub{resume: r}
	h := NewReplaceEducationsHandler(repo, &unitOfWorkStub{})
	err := h.Handle(context.Background(), ReplaceEducationsCommand{ResumeID: r.ID(), UserID: r.UserID(), Version: 1, Input: EducationsInput{Items: []EducationInput{{Location: "Berlin", University: "Uni", Grade: domain.EducationLevelMaster, ProgramName: "CS", Faculty: "Engineering"}}}})
	require.NoError(t, err)
	require.True(t, repo.saveCalled)
	require.Len(t, r.Educations(), 1)
}
