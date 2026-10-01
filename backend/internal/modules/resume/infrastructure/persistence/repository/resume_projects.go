package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func persistProjects(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeProject, 0, len(resume.Projects()))
	items := resume.Projects()
	for position, project := range items {
		if project == nil {
			return fmt.Errorf("%w: nil project", domain.ErrResumeInvariant)
		}
		if project.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: project does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := mappers.ProjectToDB(project)
		if err != nil {
			return err
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeProject{}).Error; err != nil {
		return fmt.Errorf("replace resume projects: %w", mapDBError(err))
	}
	for position, project := range items {
		model := models[position]
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("create resume project: %w", mapDBError(err))
		}
		if err := writeProjectChildrenFromDomain(tx, &model, project); err != nil {
			return err
		}
	}
	return nil
}

func loadProjects(tx *gorm.DB, resumeID identity.ID) ([]*domain.Project, error) {
	var rows []gormmodel.ResumeProject
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume projects: %w", mapDBError(err))
	}
	if len(rows) == 0 {
		return nil, nil
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

	out := make([]*domain.Project, 0, len(rows))
	for _, row := range rows {
		project, err := mappers.ProjectToDomain(row, screenshots[row.ID], technologies[row.ID], contributors[row.ID])
		if err != nil {
			return nil, err
		}
		out = append(out, project)
	}
	return out, nil
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

func writeProjectChildrenFromDomain(tx *gorm.DB, model *gormmodel.ResumeProject, project *domain.Project) error {
	projectID, err := identity.IDFromUUID(model.ID)
	if err != nil {
		return fmt.Errorf("parse project id: %w", err)
	}
	screenshots := mappers.ProjectScreenshotsToDB(projectID, project.Screenshots())
	if len(screenshots) > 0 {
		if err := tx.Create(&screenshots).Error; err != nil {
			return fmt.Errorf("create resume project screenshots: %w", mapDBError(err))
		}
	}
	technologies := mappers.ProjectTechnologiesToDB(projectID, project.Stack())
	if len(technologies) > 0 {
		if err := tx.Create(&technologies).Error; err != nil {
			return fmt.Errorf("create resume project technology: %w", mapDBError(err))
		}
	}
	contributors := mappers.ProjectContributorsToDB(projectID, project.Contributors())
	if len(contributors) > 0 {
		if err := tx.Create(&contributors).Error; err != nil {
			return fmt.Errorf("create resume project contributors: %w", mapDBError(err))
		}
	}
	return nil
}
