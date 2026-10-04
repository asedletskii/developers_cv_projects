package mappers

import (
	"fmt"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ProjectToDB(project *domain.Project) (gormmodel.ResumeProject, error) {
	if project == nil {
		return gormmodel.ResumeProject{}, fmt.Errorf("project is nil")
	}
	return gormmodel.ResumeProject{
		ID:          project.ID().UUID(),
		ResumeID:    project.ResumeID().UUID(),
		Name:        project.Name(),
		Description: project.Description(),
		Role:        project.Role(),
		Link:        ExternalURLToDB(project.Link()),
		GitLink:     ExternalURLToDB(project.GitLink()),
	}, nil
}

func ProjectToDomain(
	model gormmodel.ResumeProject,
	screenshots []gormmodel.ResumeProjectScreenshot,
	technologies []gormmodel.ResumeProjectTechnology,
	contributors []gormmodel.ResumeProjectContributor,
) (*domain.Project, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	screenshotValues := make([]domain.ExternalURL, 0, len(screenshots))
	for _, screenshot := range screenshots {
		value, err := domain.NewExternalURL(screenshot.Path)
		if err != nil {
			return nil, err
		}
		screenshotValues = append(screenshotValues, value)
	}
	stack := make([]domain.Technology, 0, len(technologies))
	for _, technology := range technologies {
		item, err := domain.NewTechnology(domain.TechnologyParams{Name: technology.Name})
		if err != nil {
			return nil, err
		}
		stack = append(stack, item)
	}
	contributorValues := make([]domain.Contributor, 0, len(contributors))
	for _, contributor := range contributors {
		var telegram *domain.Telegram
		if contributor.Telegram != nil {
			value, err := domain.NewTelegram(*contributor.Telegram)
			if err != nil {
				return nil, err
			}
			telegram = &value
		}
		var email *domain.Email
		if contributor.Email != nil {
			value, err := domain.NewEmail(*contributor.Email)
			if err != nil {
				return nil, err
			}
			email = &value
		}
		phone := ""
		if contributor.Phone != nil {
			phone = *contributor.Phone
		}
		contacts, err := domain.NewContactMethods(domain.ContactMethodsParams{
			Phone:    phone,
			Telegram: telegram,
			Email:    email,
		})
		if err != nil {
			return nil, err
		}
		item, err := domain.NewContributor(domain.ContributorParams{
			Name:     contributor.Name,
			Surname:  contributor.Surname,
			Role:     contributor.Role,
			Contacts: contacts,
		})
		if err != nil {
			return nil, err
		}
		contributorValues = append(contributorValues, item)
	}
	var link, gitLink *domain.ExternalURL
	if model.Link != nil {
		value, err := domain.NewExternalURL(*model.Link)
		if err != nil {
			return nil, err
		}
		link = &value
	}
	if model.GitLink != nil {
		value, err := domain.NewExternalURL(*model.GitLink)
		if err != nil {
			return nil, err
		}
		gitLink = &value
	}
	return domain.NewProject(domain.ProjectParams{
		ID:           id,
		ResumeID:     resumeID,
		Name:         model.Name,
		Description:  model.Description,
		Role:         model.Role,
		Screenshots:  screenshotValues,
		Link:         link,
		GitLink:      gitLink,
		Stack:        stack,
		Contributors: contributorValues,
	})
}

func ProjectScreenshotsToDB(projectID sharedvalueobjects.ID, values []domain.ExternalURL) []gormmodel.ResumeProjectScreenshot {
	models := make([]gormmodel.ResumeProjectScreenshot, 0, len(values))
	for position, value := range values {
		models = append(models, gormmodel.ResumeProjectScreenshot{
			ProjectID: projectID.UUID(),
			Position:  position,
			Path:      value.String(),
		})
	}
	return models
}

func ProjectTechnologiesToDB(projectID sharedvalueobjects.ID, values []domain.Technology) []gormmodel.ResumeProjectTechnology {
	models := make([]gormmodel.ResumeProjectTechnology, 0, len(values))
	for position, value := range values {
		models = append(models, gormmodel.ResumeProjectTechnology{
			ProjectID: projectID.UUID(),
			Position:  position,
			Name:      value.Name(),
		})
	}
	return models
}

func ProjectContributorsToDB(projectID sharedvalueobjects.ID, values []domain.Contributor) []gormmodel.ResumeProjectContributor {
	models := make([]gormmodel.ResumeProjectContributor, 0, len(values))
	for position, contributor := range values {
		item := gormmodel.ResumeProjectContributor{
			ID:        uuid.New(),
			ProjectID: projectID.UUID(),
			Position:  position,
			Name:      contributor.Name(),
			Surname:   contributor.Surname(),
			Role:      contributor.Role(),
			Phone:     contributor.Contacts().Phone(),
		}
		if value := contributor.Contacts().Telegram(); value != nil {
			text := value.String()
			item.Telegram = &text
		}
		if value := contributor.Contacts().Email(); value != nil {
			text := value.String()
			item.Email = &text
		}
		models = append(models, item)
	}
	return models
}
