package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	"developers-cv/backend/internal/modules/analytics/domain"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/mappers"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	resumelinkdomain "developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func loadAnalyticsLifecycle(
	ctx context.Context,
	tx *gorm.DB,
	linkID sharedvalueobjects.ID,
) (ports.LinkAnalyticsContext, error) {
	var model mappers.AnalyticsLifecycle
	err := tx.WithContext(ctx).
		Table("resume_links AS rl").
		Select("rl.resume_id AS resume_id, rl.status AS status, r.status AS resume_status").
		Joins("JOIN resumes AS r ON r.id = rl.resume_id").
		Where("rl.id = ?", linkID.UUID()).
		Clauses(clause.Locking{Strength: "SHARE"}).
		Take(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ports.LinkAnalyticsContext{}, domain.ErrAnalyticsSubjectNotFound
		}
		return ports.LinkAnalyticsContext{}, fmt.Errorf("load link lifecycle: %w", mapDBError(err))
	}

	linkStatus, resumeStatus, err := mappers.AnalyticsLifecycleToDomain(model)
	if err != nil {
		return ports.LinkAnalyticsContext{}, fmt.Errorf("map link lifecycle: %w", err)
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return ports.LinkAnalyticsContext{}, fmt.Errorf("map link lifecycle: resume id: %w", err)
	}
	return ports.LinkAnalyticsContext{
		LinkID:          linkID,
		ResumeID:        resumeID,
		LinkActive:      linkStatus == resumelinkdomain.ResumeLinkStatusActive,
		ResumePublished: resumeStatus == resumedomain.ResumeStatusOpened,
	}, nil
}
