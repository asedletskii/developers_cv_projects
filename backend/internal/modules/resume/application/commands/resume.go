package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CreateResumeCommand struct {
	UserID   sharedvalueobjects.ID
	Title    string
	Language domain.ResumeLanguage
}

type ReplaceResumeTitleCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
	Title    string
}

type DeleteResumeCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Version  int64
}

type CreateResumeHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewCreateResumeHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *CreateResumeHandler {
	return &CreateResumeHandler{repository: repository, uow: uow}
}

func (h *CreateResumeHandler) Handle(ctx context.Context, c CreateResumeCommand) (sharedvalueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("create resume handler dependencies are nil")
	}
	var resumeID sharedvalueobjects.ID
	err := h.uow.Do(ctx, func(tx context.Context) error {
		r, e := domain.NewResume(domain.ResumeParams{UserID: c.UserID, Title: c.Title, Language: c.Language})
		if e != nil {
			return e
		}
		resumeID = r.ID()
		overviewParams := domain.ResumeOverviewParams{
			Greeting: domain.NewGreeting(domain.GreetingParams{}),
			Sections: []domain.SectionSettings{},
		}
		if e := r.ReplaceOverview(overviewParams); e != nil {
			return e
		}
		return h.repository.Create(tx, r)
	})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}
	return resumeID, nil
}

type ChangeResumeStatusCommand struct {
	ResumeID sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Status   domain.ResumeStatus
	Version  int64
}

type ChangeResumeStatusHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewChangeResumeStatusHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ChangeResumeStatusHandler {
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
	if c.Status == domain.ResumeStatusDraft {
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
		case domain.ResumeStatusOpened:
			err = resume.Publish()
		case domain.ResumeStatusClosed:
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
	uow        ports.UnitOfWork
}

func NewReplaceResumeTitleHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *ReplaceResumeTitleHandler {
	return &ReplaceResumeTitleHandler{repository: repository, uow: uow}
}

func (h *ReplaceResumeTitleHandler) Handle(ctx context.Context, c ReplaceResumeTitleCommand) error {
	if h == nil {
		return fmt.Errorf("replace resume title handler dependencies are nil")
	}
	return executeResumeMutation(
		ctx,
		h.repository,
		h.uow,
		c.ResumeID,
		c.UserID,
		c.Version,
		"replace resume title",
		func(resume *domain.Resume) error {
			return resume.ReplaceTitle(c.Title)
		},
	)
}

type DeleteResumeHandler struct {
	repository ports.ResumeRepository
	uow        ports.UnitOfWork
}

func NewDeleteResumeHandler(repository ports.ResumeRepository, uow ports.UnitOfWork) *DeleteResumeHandler {
	return &DeleteResumeHandler{repository: repository, uow: uow}
}

func (h *DeleteResumeHandler) Handle(ctx context.Context, c DeleteResumeCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("delete resume handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid delete resume command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.Delete(tx, ports.DeleteResumeParams{ResumeID: c.ResumeID, UserID: c.UserID, Version: c.Version})
	})
}
