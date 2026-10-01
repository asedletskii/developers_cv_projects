package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ProfileToDBModels(
	p *domain.Profile,
) (gormmodel.ResumeProfile, []gormmodel.ResumeProfileWorkFormat, error) {
	if p == nil {
		return gormmodel.ResumeProfile{}, nil, fmt.Errorf("profile is nil")
	}
	m := gormmodel.ResumeProfile{
		ID:          p.ID().UUID(),
		ResumeID:    p.ResumeID().UUID(),
		Name:        p.Name(),
		Surname:     p.Surname(),
		Patronymic:  p.Patronymic(),
		Position:    p.Position(),
		PhotoPath:   p.PhotoPath().String(),
		DateOfBirth: p.DateOfBirth().Time(),
		Location:    p.Location(),
		Relocation:  p.Relocation(),
	}
	fs := make([]gormmodel.ResumeProfileWorkFormat, 0, len(p.WorkFormats()))
	for i, v := range p.WorkFormats() {
		fs = append(fs, gormmodel.ResumeProfileWorkFormat{ProfileID: m.ID, Value: string(v), Position: i})
	}
	return m, fs, nil
}

func ProfileToDomain(
	m gormmodel.ResumeProfile,
	fs []gormmodel.ResumeProfileWorkFormat,
) (*domain.Profile, error) {
	id, e := sharedvalueobjects.IDFromUUID(m.ID)
	if e != nil {
		return nil, e
	}
	rid, e := sharedvalueobjects.IDFromUUID(m.ResumeID)
	if e != nil {
		return nil, e
	}
	photo, e := domain.NewExternalURL(m.PhotoPath)
	if e != nil {
		return nil, e
	}
	dob, e := domain.NewDate(m.DateOfBirth)
	if e != nil {
		return nil, e
	}
	formats := make([]domain.WorkFormat, len(fs))
	for i, v := range fs {
		formats[i] = domain.WorkFormat(v.Value)
	}
	return domain.NewProfile(domain.ProfileParams{
		ID:          id,
		ResumeID:    rid,
		Name:        m.Name,
		Surname:     m.Surname,
		Patronymic:  m.Patronymic,
		Position:    m.Position,
		PhotoPath:   photo,
		DateOfBirth: dob,
		Location:    m.Location,
		Relocation:  m.Relocation,
		WorkFormats: formats,
	})
}
