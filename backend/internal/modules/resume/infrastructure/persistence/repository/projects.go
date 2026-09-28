package repository

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
		model, err := projectDomainModel(project)
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

func loadProjects(tx *gorm.DB, resumeID valueobjects.ID) ([]*domain.Project, error) {
	var rows []gormmodel.ResumeProject
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get resume projects: %w", mapDBError(err))
	}
	out := make([]*domain.Project, 0, len(rows))
	for _, row := range rows {
		project, err := mapProject(tx, row)
		if err != nil {
			return nil, err
		}
		out = append(out, project)
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out, nil
}

func projectDomainModel(project *domain.Project) (gormmodel.ResumeProject, error) {
	if project == nil {
		return gormmodel.ResumeProject{}, fmt.Errorf("project is nil")
	}
	return gormmodel.ResumeProject{
		ID: project.ID().UUID(), ResumeID: project.ResumeID().UUID(),
		Name: project.Name(), Description: project.Description(), Role: project.Role(),
		Link: externalURLString(project.Link()), GitLink: externalURLString(project.GitLink()),
	}, nil
}

func writeProjectChildrenFromDomain(tx *gorm.DB, model *gormmodel.ResumeProject, project *domain.Project) error {
	for position, screenshot := range project.Screenshots() {
		if err := tx.Create(&gormmodel.ResumeProjectScreenshot{
			ProjectID: model.ID, Position: position, Path: screenshot.String(),
		}).Error; err != nil {
			return fmt.Errorf("create resume project screenshot: %w", mapDBError(err))
		}
	}
	for position, technology := range project.Stack() {
		if err := tx.Create(&gormmodel.ResumeProjectTechnology{
			ProjectID: model.ID, Position: position, Name: technology.Name(), Grade: projectGradePtr(technology.Grade()),
		}).Error; err != nil {
			return fmt.Errorf("create resume project technology: %w", mapDBError(err))
		}
	}
	for position, contributor := range project.Contributors() {
		item := gormmodel.ResumeProjectContributor{
			ID: uuid.New(), ProjectID: model.ID, Position: position,
			Name: contributor.Name(), Surname: contributor.Surname(), Role: contributor.Role(),
			Phone: contributor.Contacts().Phone(),
		}
		if value := contributor.Contacts().Telegram(); value != nil {
			text := value.String()
			item.Telegram = &text
		}
		if value := contributor.Contacts().Email(); value != nil {
			text := value.String()
			item.Email = &text
		}
		if err := tx.Create(&item).Error; err != nil {
			return fmt.Errorf("create resume project contributor: %w", mapDBError(err))
		}
	}
	return nil
}

func projectGradePtr(value *enums.TechnologyGrade) *string {
	if value == nil {
		return nil
	}
	text := string(*value)
	return &text
}

func mapProject(tx *gorm.DB, model gormmodel.ResumeProject) (*domain.Project, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}

	var screenshots []gormmodel.ResumeProjectScreenshot
	if err := tx.Where("project_id = ?", model.ID).Order("position").Find(&screenshots).Error; err != nil {
		return nil, err
	}
	screenshotValues := make([]valueobjects.ExternalURL, 0, len(screenshots))
	for _, screenshot := range screenshots {
		value, err := valueobjects.NewExternalURL(screenshot.Path)
		if err != nil {
			return nil, err
		}
		screenshotValues = append(screenshotValues, value)
	}

	var technologies []gormmodel.ResumeProjectTechnology
	if err := tx.Where("project_id = ?", model.ID).Order("position").Find(&technologies).Error; err != nil {
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
	if err := tx.Where("project_id = ?", model.ID).Order("position").Find(&contributors).Error; err != nil {
		return nil, err
	}
	contributorValues := make([]entities.Contributor, 0, len(contributors))
	for _, contributor := range contributors {
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
		contacts, err := entities.NewContacts(entities.ContactsParams{Phone: phone, Telegram: telegram, Email: email})
		if err != nil {
			return nil, err
		}
		item, err := entities.NewContributor(entities.ContributorParams{
			Name: contributor.Name, Surname: contributor.Surname, Role: contributor.Role, Contacts: contacts,
		})
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

	return domain.NewProject(domain.ProjectParams{
		ID: id, ResumeID: resumeID, Name: model.Name, Description: model.Description,
		Role: model.Role, Screenshots: screenshotValues, Link: link, GitLink: gitLink,
		Stack: stack, Contributors: contributorValues,
	})
}
