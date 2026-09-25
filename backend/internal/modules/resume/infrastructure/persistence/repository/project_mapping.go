package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

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
