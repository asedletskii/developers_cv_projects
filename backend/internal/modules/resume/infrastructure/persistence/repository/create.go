package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/platform/database"
)

func (r *ResumeRepository) Create(ctx context.Context, command commands.CreateResumeRepositoryCommand) error {
	resume := command.Resume
	if resume == nil {
		return fmt.Errorf("create resume: %w", domain.ErrInvalidResume)
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume: %w", err)
	}
	if err := createAggregate(ctx, tx, resume); err != nil {
		return fmt.Errorf("create resume: %w", err)
	}
	return nil
}

func createAggregate(ctx context.Context, tx *gorm.DB, resume *domain.Resume) error {
	m := toResumeModel(resume)
	if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	if err := createProfile(ctx, tx, resume); err != nil {
		return err
	}
	if err := createSettings(ctx, tx, resume); err != nil {
		return err
	}
	if err := createContent(ctx, tx, resume); err != nil {
		return err
	}
	return createProjects(ctx, tx, resume)
}
