package commands

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func executeResumeMutation(
	ctx context.Context,
	repository ports.ResumeRepository,
	uow ports.UnitOfWork,
	resumeID sharedvalueobjects.ID,
	userID sharedvalueobjects.ID,
	version int64,
	operation string,
	mutate func(*domain.Resume) error,
) error {
	if repository == nil || uow == nil {
		return fmt.Errorf("%s handler dependencies are nil", operation)
	}
	if resumeID.IsZero() || userID.IsZero() || version < 1 {
		return fmt.Errorf("invalid %s command", operation)
	}
	return uow.Do(ctx, func(tx context.Context) error {
		resume, err := repository.Get(tx, ports.GetResumeParams{ResumeID: resumeID, UserID: userID})
		if err != nil {
			return err
		}
		if resume.Version() != version {
			return domain.ErrConcurrencyConflict
		}
		if err := mutate(resume); err != nil {
			return err
		}
		return repository.Save(tx, resume)
	})
}
