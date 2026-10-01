package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ActivityInput struct {
	ID          identity.ID
	Type        domain.ActivityType
	Title       string
	Description *string
	Date        *domain.YearMonth
	Link        *domain.ExternalURL
}

type ActivitiesInput struct{ Items []ActivityInput }

func (i ActivitiesInput) toDomain() []domain.ActivityParams {
	out := make([]domain.ActivityParams, 0, len(i.Items))
	for _, x := range i.Items {
		out = append(out, domain.ActivityParams{ID: x.ID, Type: x.Type, Title: x.Title, Description: x.Description, Date: x.Date, Link: x.Link})
	}
	return out
}

type ReplaceActivitiesCommand struct {
	ResumeID identity.ID
	UserID   identity.ID
	Version  int64
	Input    ActivitiesInput
}

type ReplaceActivitiesHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewReplaceActivitiesHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceActivitiesHandler {
	return &ReplaceActivitiesHandler{repository: repository, uow: uow}
}

func (h *ReplaceActivitiesHandler) Handle(ctx context.Context, c ReplaceActivitiesCommand) error {
	if h == nil {
		return fmt.Errorf("replace activities handler dependencies are nil")
	}

	return executeResumeMutation(ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version, "replace activities", func(resume *domain.Resume) error {
		return resume.ReplaceActivities(c.Input.toDomain())
	})
}
