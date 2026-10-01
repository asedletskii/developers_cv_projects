package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TechnologiesToDomain(
	rows []gormmodel.ResumeTechnology,
	resumeID sharedvalueobjects.ID,
) (*domain.ResumeTechnologies, error) {
	if len(rows) == 0 {
		return nil, nil
	}
	collectionID, err := sharedvalueobjects.IDFromUUID(rows[0].ID)
	if err != nil {
		return nil, err
	}
	entries := make([]domain.ResumeTechnologyEntryParams, 0, len(rows))
	for _, row := range rows {
		id, err := sharedvalueobjects.IDFromUUID(row.ID)
		if err != nil {
			return nil, err
		}
		var grade *domain.TechnologyGrade
		if row.Grade != nil {
			value := domain.TechnologyGrade(*row.Grade)
			if !value.IsValid() {
				return nil, domain.ErrInvalidEntity
			}
			grade = &value
		}
		entries = append(entries, domain.ResumeTechnologyEntryParams{
			ID:       id,
			Name:     row.Name,
			Grade:    grade,
			Position: row.Position,
		})
	}
	return domain.NewResumeTechnologies(domain.ResumeTechnologiesParams{
		ID:       collectionID,
		ResumeID: resumeID,
		Entries:  entries,
	})
}

func TechnologyToDB(
	x domain.ResumeTechnologyEntry,
	resumeID sharedvalueobjects.ID,
) (gormmodel.ResumeTechnology, error) {
	if x.ID().IsZero() {
		return gormmodel.ResumeTechnology{}, fmt.Errorf("technology id is required")
	}
	m := gormmodel.ResumeTechnology{
		ID:       x.ID().UUID(),
		ResumeID: resumeID.UUID(),
		Name:     x.Name(),
		Position: x.Position(),
	}
	if g := x.Grade(); g != nil {
		if !g.IsValid() {
			return gormmodel.ResumeTechnology{}, fmt.Errorf("invalid technology grade")
		}
		value := string(*g)
		m.Grade = &value
	}
	return m, nil
}
