package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"fmt"
	"strings"
)

type CreateResumeCommand struct {
	UserID   useridentity.ID
	Title    string
	Language enums.ResumeLanguage
}
type ReplaceResumeTitleCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Title    string
}
type DeleteResumeCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type CreateResumeHandler struct {
	repository ports.CreateResumeRepository
	overview   ports.CreateOverviewRepository
	uow        UnitOfWork
}

func NewCreateResumeHandler(repository ports.CreateResumeRepository, overview ports.CreateOverviewRepository, uow UnitOfWork) *CreateResumeHandler {
	return &CreateResumeHandler{repository: repository, overview: overview, uow: uow}
}
func (h *CreateResumeHandler) Handle(ctx context.Context, c CreateResumeCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.overview == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume handler dependencies are nil")
	}
	resumeID := valueobjects.NewID()
	err := h.uow.Do(ctx, func(tx context.Context) error {
		r, e := domain.NewResume(domain.ResumeParams{ID: resumeID, UserID: c.UserID, Title: c.Title, Language: c.Language})
		if e != nil {
			return e
		}
		if e := h.repository.CreateResume(tx, ports.CreateResumeRepositoryCommand{Resume: r}); e != nil {
			return e
		}
		o, e := domain.NewResumeOverview(domain.ResumeOverviewParams{Greeting: valueobjects.NewGreeting(valueobjects.GreetingParams{}), Sections: []valueobjects.SectionSettings{}})
		if e != nil {
			return e
		}
		return h.overview.CreateOverview(tx, ports.CreateOverviewRepositoryCommand{ResumeID: resumeID, UserID: c.UserID, Overview: o})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return resumeID, nil
}

type ReplaceResumeTitleHandler struct {
	replacer ports.ReplaceResumeTitleRepository
	uow      UnitOfWork
}

func NewReplaceResumeTitleHandler(replacer ports.ReplaceResumeTitleRepository, uow UnitOfWork) *ReplaceResumeTitleHandler {
	return &ReplaceResumeTitleHandler{replacer: replacer, uow: uow}
}

func (h *ReplaceResumeTitleHandler) Handle(ctx context.Context, c ReplaceResumeTitleCommand) error {
	if h == nil || h.replacer == nil || h.uow == nil {
		return fmt.Errorf("replace resume title handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid replace resume title command")
	}
	title := strings.TrimSpace(c.Title)
	if err := domain.ValidateResumeTitle(title); err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.replacer.ReplaceResumeTitle(tx, ports.ReplaceResumeTitleRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID, Title: title})
	})
}

type DeleteResumeHandler struct {
	repository ports.DeleteResumeRepository
	uow        UnitOfWork
}

func NewDeleteResumeHandler(repository ports.DeleteResumeRepository, uow UnitOfWork) *DeleteResumeHandler {
	return &DeleteResumeHandler{repository: repository, uow: uow}
}
func (h *DeleteResumeHandler) Handle(ctx context.Context, c DeleteResumeCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("delete resume handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid delete resume command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.DeleteResume(tx, ports.DeleteResumeRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID})
	})
}
