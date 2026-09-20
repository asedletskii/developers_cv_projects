package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	resumelinkcontracts "developers-cv/backend/internal/modules/resume_link/contracts"
	usercontracts "developers-cv/backend/internal/modules/user/contracts"
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
	users      usercontracts.UserExistenceChecker
	uow        ports.UnitOfWork
}

func NewCreateResumeHandler(
	repository ports.ResumeRepository,
	users usercontracts.UserExistenceChecker,
	uow ports.UnitOfWork,
) *CreateResumeHandler {
	return &CreateResumeHandler{repository: repository, users: users, uow: uow}
}

func (h *CreateResumeHandler) Handle(ctx context.Context, c CreateResumeCommand) (sharedvalueobjects.ID, error) {
	if h == nil || h.repository == nil || h.users == nil || h.uow == nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("create resume handler dependencies are nil")
	}
	resume, err := domain.NewResume(domain.ResumeParams{UserID: c.UserID, Title: c.Title, Language: c.Language})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}
	if err := resume.ReplaceOverview(domain.ResumeOverviewParams{
		Greeting: domain.Greeting{},
		Sections: []domain.SectionSettings{},
	}); err != nil {
		return sharedvalueobjects.ID{}, err
	}
	err = h.uow.Do(ctx, func(tx context.Context) error {
		if err := h.users.EnsureUserExists(tx, c.UserID); err != nil {
			return err
		}
		return h.repository.Create(tx, resume)
	})
	if err != nil {
		return sharedvalueobjects.ID{}, err
	}
	return resume.ID(), nil
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
	if h == nil {
		return fmt.Errorf("change resume status handler dependencies are nil")
	}
	return executeResumeMutation(
		ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version,
		"change resume status", ports.ResumeMutationStatus,
		func(resume *domain.Resume) error {
			return resume.ReplaceStatus(c.Status)
		},
	)
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
		ctx, h.repository, h.uow, c.ResumeID, c.UserID, c.Version,
		"replace resume title", ports.ResumeMutationTitle,
		func(resume *domain.Resume) error {
			return resume.ReplaceTitle(c.Title)
		},
	)
}

type DeleteResumeHandler struct {
	repository ports.ResumeRepository
	links      resumelinkcontracts.ResumeLinkCleanup
	uow        ports.UnitOfWork
}

func NewDeleteResumeHandler(
	repository ports.ResumeRepository,
	links resumelinkcontracts.ResumeLinkCleanup,
	uow ports.UnitOfWork,
) *DeleteResumeHandler {
	return &DeleteResumeHandler{repository: repository, links: links, uow: uow}
}

func (h *DeleteResumeHandler) Handle(ctx context.Context, c DeleteResumeCommand) error {
	if h == nil || h.repository == nil || h.links == nil || h.uow == nil {
		return fmt.Errorf("delete resume handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() || c.Version < 1 {
		return fmt.Errorf("invalid delete resume command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		params := ports.DeleteResumeParams{ResumeID: c.ResumeID, UserID: c.UserID, Version: c.Version}
		if err := h.repository.LockForDelete(tx, params); err != nil {
			return err
		}
		if err := h.links.DeleteByResumeIDs(tx, []sharedvalueobjects.ID{c.ResumeID}); err != nil {
			return fmt.Errorf("delete resume: clean links: %w", err)
		}
		return h.repository.Delete(tx, params)
	})
}
