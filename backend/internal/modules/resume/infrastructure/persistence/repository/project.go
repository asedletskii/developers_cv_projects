package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

var _ ports.CreateProjectRepository = (*ProjectRepository)(nil)
var _ ports.GetProjectRepository = (*ProjectRepository)(nil)
var _ ports.UpdateProjectRepository = (*ProjectRepository)(nil)
var _ ports.DeleteProjectRepository = (*ProjectRepository)(nil)

func (r *ProjectRepository) CreateProject(ctx context.Context, c ports.CreateProjectRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	if c.Project == nil {
		return fmt.Errorf("project is nil")
	}
	item, err := projectDomainModel(c.Project)
	if err != nil {
		return err
	}
	item.SortOrder, err = nextSortOrder(tx, item.TableName(), c.ResumeID.UUID())
	if err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
		return err
	}
	return writeProjectChildrenFromDomain(ctx, tx, &item, c.Project)
}

func (r *ProjectRepository) GetProject(ctx context.Context, q ports.GetProjectRepositoryCommand) (*domain.Project, error) {
	tx, err := r.owned(ctx, q.ResumeID, q.UserID)
	if err != nil {
		return nil, err
	}
	var model gormmodel.ResumeProject
	if err := tx.WithContext(ctx).Where("id = ? AND resume_id = ?", q.ProjectID.UUID(), q.ResumeID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}
	return r.mapProject(ctx, tx, model)
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, c ports.UpdateProjectRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	if c.Project == nil {
		return fmt.Errorf("project is nil")
	}
	model, err := projectDomainModel(c.Project)
	if err != nil {
		return err
	}
	var existing gormmodel.ResumeProject
	if err := tx.WithContext(ctx).Where("id = ? AND resume_id = ?", c.ProjectID.UUID(), c.ResumeID.UUID()).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return err
	}
	if err := tx.WithContext(ctx).Model(&gormmodel.ResumeProject{}).Where("id = ? AND resume_id = ?", c.ProjectID.UUID(), c.ResumeID.UUID()).Updates(map[string]any{
		"name": model.Name, "description": model.Description, "role": model.Role, "link": model.Link, "git_link": model.GitLink,
	}).Error; err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Where("project_id = ?", c.ProjectID.UUID()).Delete(&gormmodel.ResumeProjectScreenshot{}).Error; err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Where("project_id = ?", c.ProjectID.UUID()).Delete(&gormmodel.ResumeProjectTechnology{}).Error; err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Where("project_id = ?", c.ProjectID.UUID()).Delete(&gormmodel.ResumeProjectContributor{}).Error; err != nil {
		return err
	}
	model.ID = existing.ID
	return writeProjectChildrenFromDomain(ctx, tx, &model, c.Project)
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, c ports.DeleteProjectRepositoryCommand) error {
	tx, err := r.owned(ctx, c.ResumeID, c.UserID)
	if err != nil {
		return err
	}
	result := tx.WithContext(ctx).Where("id = ? AND resume_id = ?", c.ProjectID.UUID(), c.ResumeID.UUID()).Delete(&gormmodel.ResumeProject{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func projectDomainModel(item *domain.Project) (gormmodel.ResumeProject, error) {
	return gormmodel.ResumeProject{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Name: item.Name(), Description: item.Description(), Role: item.Role(), Link: externalURLString(item.Link()), GitLink: externalURLString(item.GitLink())}, nil
}

func projectGradePtr(value *enums.TechnologyGrade) *string {
	if value == nil {
		return nil
	}
	text := string(*value)
	return &text
}
