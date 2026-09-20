package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TechnologiesToDomain(rows []gormmodel.ResumeTechnology) (*domain.ResumeTechnologies, error) {
	if len(rows) == 0 {
		return nil, nil
	}
	entries := make([]domain.ResumeTechnologyEntryParams, 0, len(rows))
	for i, row := range rows {
		id, err := sharedvalueobjects.IDFromUUID(row.ID)
		if err != nil {
			return nil, fmt.Errorf("map technology at index %d: %w", i, err)
		}
		entries = append(entries, domain.ResumeTechnologyEntryParams{
			ID:       id,
			Name:     row.Name,
			Position: row.Position,
		})
	}
	return domain.RestoreResumeTechnologies(domain.ResumeTechnologiesParams{
		Entries: entries,
	}), nil
}

func TechnologyToDB(
	x domain.ResumeTechnologyEntry,
	resumeID sharedvalueobjects.ID,
) gormmodel.ResumeTechnology {
	return gormmodel.ResumeTechnology{
		ID:       x.ID().UUID(),
		ResumeID: resumeID.UUID(),
		Name:     x.Name(),
		Position: x.Position(),
	}
}

func ResumeTechnologiesToDBModels(
	collection *domain.ResumeTechnologies,
	resumeID sharedvalueobjects.ID,
) []gormmodel.ResumeTechnology {
	if collection == nil {
		return nil
	}
	entries := collection.Entries()
	models := make([]gormmodel.ResumeTechnology, 0, len(entries))
	for _, entry := range entries {
		models = append(models, TechnologyToDB(entry, resumeID))
	}
	return models
}
