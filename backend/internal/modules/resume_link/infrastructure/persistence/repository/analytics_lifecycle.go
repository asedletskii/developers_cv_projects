package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/mappers"
)

func loadAnalyticsLifecycle(
	ctx context.Context,
	tx *gorm.DB,
	linkID uuid.UUID,
) (domain.ResumeLinkStatus, resumedomain.ResumeStatus, error) {
	var model mappers.AnalyticsLifecycle
	err := tx.WithContext(ctx).
		Table("resume_links AS rl").
		Select("rl.status AS status, r.status AS resume_status").
		Joins("JOIN resumes AS r ON r.id = rl.resume_id").
		Where("rl.id = ?", linkID).
		Take(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", domain.ErrEntityNotFound
		}
		return "", "", fmt.Errorf("load link lifecycle: %w", mapDBError(err))
	}

	linkStatus, resumeStatus, err := mappers.AnalyticsLifecycleToDomain(model)
	if err != nil {
		return "", "", fmt.Errorf("map link lifecycle: %w", err)
	}
	return linkStatus, resumeStatus, nil
}
