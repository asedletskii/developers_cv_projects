package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type ResumeLinkReader struct {
	db *gorm.DB
}

func NewResumeLinkReader(db *gorm.DB) *ResumeLinkReader {
	return &ResumeLinkReader{db: db}
}

func (r *ResumeLinkReader) GetByID(
	ctx context.Context,
	query queries.PublicResumeLinkQuery,
) (queries.ResumeLinkPublicView, error) {
	if r == nil || r.db == nil {
		return queries.ResumeLinkPublicView{}, fmt.Errorf("get resume link: reader dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.ResumeLink
	if err := db.Where("id = ?", query.LinkID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ResumeLinkPublicView{}, domain.ErrEntityNotFound
		}
		return queries.ResumeLinkPublicView{}, fmt.Errorf("get resume link: %w", err)
	}
	return mapView(model), nil
}

func (r *ResumeLinkReader) LockByID(
	ctx context.Context,
	query queries.PublicResumeLinkQuery,
) (queries.ResumeLinkPublicView, error) {
	if r == nil || r.db == nil {
		return queries.ResumeLinkPublicView{}, fmt.Errorf("lock resume link: reader dependencies are nil")
	}
	db, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return queries.ResumeLinkPublicView{}, fmt.Errorf("lock resume link: %w", err)
	}
	var model gormmodel.ResumeLink
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "SHARE"}).Where("id = ?", query.LinkID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ResumeLinkPublicView{}, domain.ErrEntityNotFound
		}
		return queries.ResumeLinkPublicView{}, fmt.Errorf("lock resume link: %w", err)
	}
	return mapView(model), nil
}

func mapView(model gormmodel.ResumeLink) queries.ResumeLinkPublicView {
	return queries.ResumeLinkPublicView{
		ID: model.ID.String(), ResumeID: model.ResumeID.String(), Status: model.Status,
	}
}

func (r *ResumeLinkReader) ListByResumeID(
	ctx context.Context,
	query queries.ListResumeLinksQuery,
) ([]queries.ResumeLinkListItem, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("list resume links: reader dependencies are nil")
	}
	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var models []gormmodel.ResumeLink
	if err := db.Where("resume_id = ?", query.ResumeID.UUID()).Order("created_at, id").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list resume links: %w", err)
	}
	result := make([]queries.ResumeLinkListItem, 0, len(models))
	for _, model := range models {
		result = append(result, queries.ResumeLinkListItem{
			ID:        model.ID.String(),
			Name:      model.Name,
			Status:    model.Status,
			Version:   model.Version,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		})
	}
	return result, nil
}
