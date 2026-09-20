package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistProjects(tx *gorm.DB, resume *domain.Resume) error {
	models, screenshots, technologies, contributors := mappers.ProjectsToDBModels(resume.Projects())
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeProject{}).Error; err != nil {
		return fmt.Errorf("replace resume projects: %w", mapDBError(err))
	}
	if len(models) > 0 {
		if err := tx.Create(&models).Error; err != nil {
			return fmt.Errorf("create resume projects: %w", mapDBError(err))
		}
	}
	if len(screenshots) > 0 {
		if err := tx.Create(&screenshots).Error; err != nil {
			return fmt.Errorf("create resume project screenshots: %w", mapDBError(err))
		}
	}
	if len(technologies) > 0 {
		if err := tx.Create(&technologies).Error; err != nil {
			return fmt.Errorf("create resume project technologies: %w", mapDBError(err))
		}
	}
	if len(contributors) > 0 {
		if err := tx.Create(&contributors).Error; err != nil {
			return fmt.Errorf("create resume project contributors: %w", mapDBError(err))
		}
	}
	return nil
}

func loadProjects(tx *gorm.DB, resumeID sharedvalueobjects.ID) ([]*domain.Project, error) {
	var rows []gormmodel.ResumeProject
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume projects: %w", mapDBError(err))
	}
	if len(rows) == 0 {
		return make([]*domain.Project, 0), nil
	}

	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	screenshots, err := loadProjectScreenshots(tx, ids)
	if err != nil {
		return nil, err
	}
	technologies, err := loadProjectTechnologies(tx, ids)
	if err != nil {
		return nil, err
	}
	contributors, err := loadProjectContributors(tx, ids)
	if err != nil {
		return nil, err
	}

	return mappers.ProjectsToDomain(rows, screenshots, technologies, contributors)
}

func loadProjectScreenshots(tx *gorm.DB, projectIDs []uuid.UUID) (map[uuid.UUID][]gormmodel.ResumeProjectScreenshot, error) {
	var rows []gormmodel.ResumeProjectScreenshot
	if err := tx.Where("project_id IN ?", projectIDs).Order("project_id, position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load resume project screenshots: %w", mapDBError(err))
	}
	grouped := make(map[uuid.UUID][]gormmodel.ResumeProjectScreenshot, len(projectIDs))
	for _, row := range rows {
		grouped[row.ProjectID] = append(grouped[row.ProjectID], row)
	}
	return grouped, nil
}

func loadProjectTechnologies(tx *gorm.DB, projectIDs []uuid.UUID) (map[uuid.UUID][]gormmodel.ResumeProjectTechnology, error) {
	var rows []gormmodel.ResumeProjectTechnology
	if err := tx.Where("project_id IN ?", projectIDs).Order("project_id, position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load resume project technologies: %w", mapDBError(err))
	}
	grouped := make(map[uuid.UUID][]gormmodel.ResumeProjectTechnology, len(projectIDs))
	for _, row := range rows {
		grouped[row.ProjectID] = append(grouped[row.ProjectID], row)
	}
	return grouped, nil
}

func loadProjectContributors(tx *gorm.DB, projectIDs []uuid.UUID) (map[uuid.UUID][]gormmodel.ResumeProjectContributor, error) {
	var rows []gormmodel.ResumeProjectContributor
	if err := tx.Where("project_id IN ?", projectIDs).Order("project_id, position").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load resume project contributors: %w", mapDBError(err))
	}
	grouped := make(map[uuid.UUID][]gormmodel.ResumeProjectContributor, len(projectIDs))
	for _, row := range rows {
		grouped[row.ProjectID] = append(grouped[row.ProjectID], row)
	}
	return grouped, nil
}
