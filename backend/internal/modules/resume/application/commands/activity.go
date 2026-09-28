package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ActivityInput struct {
	ID          valueobjects.ID
	Type        enums.ActivityType
	Title       string
	Description *string
	Date        *valueobjects.YearMonth
	Link        *valueobjects.ExternalURL
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
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Input    ActivitiesInput
}

type ReplaceActivitiesHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceActivitiesHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceActivitiesHandler {
	return &ReplaceActivitiesHandler{repository: repository, uow: uow}
}

func (h *ReplaceActivitiesHandler) Handle(ctx context.Context, c ReplaceActivitiesCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace activities handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace activities command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		r, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if r.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := r.ReplaceActivities(c.Input.toDomain()); err != nil {
			return err
		}
		return h.repository.Save(tx, r)
	})
}
