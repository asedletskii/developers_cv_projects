package mappers

import (
	"fmt"

	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ProjectsToDomain(
	rows []gormmodel.ResumeProject,
	screenshots map[uuid.UUID][]gormmodel.ResumeProjectScreenshot,
	technologies map[uuid.UUID][]gormmodel.ResumeProjectTechnology,
	contributors map[uuid.UUID][]gormmodel.ResumeProjectContributor,
) ([]*domain.Project, error) {
	projects := make([]*domain.Project, 0, len(rows))
	for i, row := range rows {
		project, err := ProjectToDomain(row, screenshots[row.ID], technologies[row.ID], contributors[row.ID])
		if err != nil {
			return nil, fmt.Errorf("map project at index %d: %w", i, err)
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func ProjectToDB(project *domain.Project) gormmodel.ResumeProject {
	return gormmodel.ResumeProject{
		ID:          project.ID().UUID(),
		ResumeID:    project.ResumeID().UUID(),
		Name:        project.Name(),
		Description: project.Description(),
		Role:        project.Role(),
		Link:        ExternalURLToDB(project.Link()),
		GitLink:     ExternalURLToDB(project.GitLink()),
	}
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
		screenshotValues = append(screenshotValues, domain.RestoreExternalURL(screenshot.Path))
	}
	stack := make([]domain.Technology, 0, len(technologies))
	for _, technology := range technologies {
		stack = append(stack, domain.RestoreTechnology(domain.TechnologyParams{Name: technology.Name}))
	}
	contributorValues := make([]domain.Contributor, 0, len(contributors))
	for _, contributor := range contributors {
		var telegram *domain.Telegram
		if contributor.Telegram != nil {
			value := domain.RestoreTelegram(*contributor.Telegram)
			telegram = &value
		}
		var email *domain.Email
		if contributor.Email != nil {
			value := domain.RestoreEmail(*contributor.Email)
			email = &value
		}
		contacts := domain.RestoreContactMethods(domain.ContactMethodsParams{
			Phone: PointerStringValue(contributor.Phone), Telegram: telegram, Email: email,
		})
		contributorValues = append(contributorValues, domain.RestoreContributor(domain.ContributorParams{
			Name: contributor.Name, Surname: contributor.Surname, Role: contributor.Role, Contacts: contacts,
		}))
	}
	var link, gitLink *domain.ExternalURL
	if model.Link != nil {
		value := domain.RestoreExternalURL(*model.Link)
		link = &value
	}
	if model.GitLink != nil {
		value := domain.RestoreExternalURL(*model.GitLink)
		gitLink = &value
	}
	return domain.RestoreProject(domain.ProjectSnapshot{
		ID: id, ResumeID: resumeID, Name: model.Name, Description: model.Description, Role: model.Role,
		Screenshots: screenshotValues, Link: link, GitLink: gitLink, Stack: stack, Contributors: contributorValues,
	}), nil
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

func ProjectsToDBModels(items []*domain.Project) (
	[]gormmodel.ResumeProject,
	[]gormmodel.ResumeProjectScreenshot,
	[]gormmodel.ResumeProjectTechnology,
	[]gormmodel.ResumeProjectContributor,
) {
	models := make([]gormmodel.ResumeProject, 0, len(items))
	screenshots := make([]gormmodel.ResumeProjectScreenshot, 0)
	technologies := make([]gormmodel.ResumeProjectTechnology, 0)
	contributors := make([]gormmodel.ResumeProjectContributor, 0)
	for position, item := range items {
		model := ProjectToDB(item)
		model.SortOrder = position
		models = append(models, model)
		screenshots = append(screenshots, ProjectScreenshotsToDB(item.ID(), item.Screenshots())...)
		technologies = append(technologies, ProjectTechnologiesToDB(item.ID(), item.Stack())...)
		contributors = append(contributors, ProjectContributorsToDB(item.ID(), item.Contributors())...)
	}
	return models, screenshots, technologies, contributors
}
