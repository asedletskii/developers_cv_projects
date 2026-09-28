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

type CreateResumeCommand struct {
	UserID   useridentity.ID
	Title    string
	Language enums.ResumeLanguage
}

type ReplaceResumeTitleCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Version  int64
	Title    string
}

type DeleteResumeCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type CreateResumeHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewCreateResumeHandler(repository ports.ResumeRepository, uow UnitOfWork) *CreateResumeHandler {
	return &CreateResumeHandler{repository: repository, uow: uow}
}

func (h *CreateResumeHandler) Handle(ctx context.Context, c CreateResumeCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume handler dependencies are nil")
	}
	var resumeID valueobjects.ID
	err := h.uow.Do(ctx, func(tx context.Context) error {
		r, e := domain.NewResume(domain.ResumeParams{UserID: c.UserID, Title: c.Title, Language: c.Language})
		if e != nil {
			return e
		}
		resumeID = r.ID()
		overviewParams := domain.ResumeOverviewParams{
			Greeting: valueobjects.NewGreeting(valueobjects.GreetingParams{}),
			Sections: []valueobjects.SectionSettings{},
		}
		if e := r.ReplaceOverview(overviewParams); e != nil {
			return e
		}
		return h.repository.Create(tx, r)
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return resumeID, nil
}

type ChangeResumeStatusCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
	Status   enums.ResumeStatus
	Version  int64
}

type ChangeResumeStatusHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewChangeResumeStatusHandler(repository ports.ResumeRepository, uow UnitOfWork) *ChangeResumeStatusHandler {
	return &ChangeResumeStatusHandler{repository: repository, uow: uow}
}

func (h *ChangeResumeStatusHandler) Handle(ctx context.Context, c ChangeResumeStatusCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("change resume status handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid change resume status command")
	}
	if !c.Status.IsValid() {
		return fmt.Errorf("%w: invalid target status %q", domain.ErrInvalidResume, c.Status)
	}
	if c.Status == enums.ResumeStatusDraft {
		return fmt.Errorf("%w: target status %q is not allowed", domain.ErrInvalidResumeStatusTransition, c.Status)
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		switch c.Status {
		case enums.ResumeStatusOpened:
			err = resume.Publish()
		case enums.ResumeStatusClosed:
			err = resume.Close()
		}
		if err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}

type ReplaceResumeTitleHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewReplaceResumeTitleHandler(repository ports.ResumeRepository, uow UnitOfWork) *ReplaceResumeTitleHandler {
	return &ReplaceResumeTitleHandler{repository: repository, uow: uow}
}

func (h *ReplaceResumeTitleHandler) Handle(ctx context.Context, c ReplaceResumeTitleCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("replace resume title handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid replace resume title command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := h.repository.Get(tx, ports.GetResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if resume.Version() != c.Version {
			return domain.ErrConcurrencyConflict
		}
		if err := resume.ReplaceTitle(c.Title); err != nil {
			return err
		}
		return h.repository.Save(tx, resume)
	})
}

type DeleteResumeHandler struct {
	repository ports.ResumeRepository
	uow        UnitOfWork
}

func NewDeleteResumeHandler(repository ports.ResumeRepository, uow UnitOfWork) *DeleteResumeHandler {
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
		return h.repository.Delete(tx, ports.DeleteResumeParams{ResumeID: c.ResumeID, UserID: c.UserID})
	})
}
