package repository

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	if err := tx.WithContext(ctx).Where("id = ? AND resume_id = ?", c.Project.ID().UUID(), c.ResumeID.UUID()).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return err
	}
	if err := tx.WithContext(ctx).Model(&gormmodel.ResumeProject{}).Where("id = ? AND resume_id = ?", c.Project.ID().UUID(), c.ResumeID.UUID()).Updates(map[string]any{
		"name": model.Name, "description": model.Description, "role": model.Role, "link": model.Link, "git_link": model.GitLink,
	}).Error; err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Where("project_id = ?", c.Project.ID().UUID()).Delete(&gormmodel.ResumeProjectScreenshot{}).Error; err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Where("project_id = ?", c.Project.ID().UUID()).Delete(&gormmodel.ResumeProjectTechnology{}).Error; err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Where("project_id = ?", c.Project.ID().UUID()).Delete(&gormmodel.ResumeProjectContributor{}).Error; err != nil {
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

func writeProjectChildrenFromDomain(ctx context.Context, tx *gorm.DB, model *gormmodel.ResumeProject, project *domain.Project) error {
	for position, screenshot := range project.Screenshots() {
		if err := tx.WithContext(ctx).Create(&gormmodel.ResumeProjectScreenshot{ProjectID: model.ID, Position: position, Path: screenshot.String()}).Error; err != nil {
			return err
		}
	}
	for position, technology := range project.Stack() {
		if err := tx.WithContext(ctx).Create(&gormmodel.ResumeProjectTechnology{ProjectID: model.ID, Position: position, Name: technology.Name(), Grade: projectGradePtr(technology.Grade())}).Error; err != nil {
			return err
		}
	}
	for position, contributor := range project.Contributors() {
		item := gormmodel.ResumeProjectContributor{ID: uuid.New(), ProjectID: model.ID, Position: position, Name: contributor.Name(), Surname: contributor.Surname(), Role: contributor.Role(), Phone: contributor.Contacts().Phone()}
		if value := contributor.Contacts().Telegram(); value != nil {
			text := value.String()
			item.Telegram = &text
		}
		if value := contributor.Contacts().Email(); value != nil {
			text := value.String()
			item.Email = &text
		}
		if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectRepository) mapProject(ctx context.Context, tx *gorm.DB, model gormmodel.ResumeProject) (*domain.Project, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}

	var screenshots []gormmodel.ResumeProjectScreenshot
	if err := tx.WithContext(ctx).Where("project_id = ?", model.ID).Order("position").Find(&screenshots).Error; err != nil {
		return nil, err
	}
	screenshotValues := make([]valueobjects.AssetPath, 0, len(screenshots))
	for _, screenshot := range screenshots {
		value, err := valueobjects.NewAssetPath(screenshot.Path)
		if err != nil {
			return nil, err
		}
		screenshotValues = append(screenshotValues, value)
	}

	var technologies []gormmodel.ResumeProjectTechnology
	if err := tx.WithContext(ctx).Where("project_id = ?", model.ID).Order("position").Find(&technologies).Error; err != nil {
		return nil, err
	}
	stack := make([]entities.Technology, 0, len(technologies))
	for _, technology := range technologies {
		var grade *enums.TechnologyGrade
		if technology.Grade != nil {
			value := enums.TechnologyGrade(*technology.Grade)
			grade = &value
		}
		item, err := entities.NewTechnology(entities.TechnologyParams{Name: technology.Name, Grade: grade})
		if err != nil {
			return nil, err
		}
		stack = append(stack, item)
	}

	var contributors []gormmodel.ResumeProjectContributor
	if err := tx.WithContext(ctx).Where("project_id = ?", model.ID).Order("position").Find(&contributors).Error; err != nil {
		return nil, err
	}
	contributorValues := make([]entities.Contributor, 0, len(contributors))
	for _, contributor := range contributors {
		contacts := entities.EmptyContacts()
		if contributor.Phone != nil || contributor.Telegram != nil || contributor.Email != nil {
			var telegram *valueobjects.Telegram
			if contributor.Telegram != nil {
				value, err := valueobjects.NewTelegram(*contributor.Telegram)
				if err != nil {
					return nil, err
				}
				telegram = &value
			}
			var email *valueobjects.Email
			if contributor.Email != nil {
				value, err := valueobjects.NewEmail(*contributor.Email)
				if err != nil {
					return nil, err
				}
				email = &value
			}
			phone := ""
			if contributor.Phone != nil {
				phone = *contributor.Phone
			}
			var err error
			contacts, err = entities.NewContacts(entities.ContactsParams{Phone: phone, Telegram: telegram, Email: email})
			if err != nil {
				return nil, err
			}
		}
		item, err := entities.NewContributor(entities.ContributorParams{Name: contributor.Name, Surname: contributor.Surname, Role: contributor.Role, Contacts: contacts})
		if err != nil {
			return nil, err
		}
		contributorValues = append(contributorValues, item)
	}

	var link, gitLink *valueobjects.ExternalURL
	if model.Link != nil {
		value, err := valueobjects.NewExternalURL(*model.Link)
		if err != nil {
			return nil, err
		}
		link = &value
	}
	if model.GitLink != nil {
		value, err := valueobjects.NewExternalURL(*model.GitLink)
		if err != nil {
			return nil, err
		}
		gitLink = &value
	}
	return domain.NewProject(domain.ProjectParams{ID: id, ResumeID: resumeID, Name: model.Name, Description: model.Description, Role: model.Role, Screenshots: screenshotValues, Link: link, GitLink: gitLink, Stack: stack, Contributors: contributorValues})
}
