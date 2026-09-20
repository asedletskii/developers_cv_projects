package mappers

import (
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ProfileToDBModels(
	p *domain.Profile,
) (gormmodel.ResumeProfile, []gormmodel.ResumeProfileWorkFormat) {
	m := gormmodel.ResumeProfile{
		ID:          p.ID().UUID(),
		ResumeID:    p.ResumeID().UUID(),
		Name:        p.Name(),
		Surname:     p.Surname(),
		Patronymic:  p.Patronymic(),
		Position:    p.Position(),
		PhotoURL:    p.PhotoURL().String(),
		DateOfBirth: p.DateOfBirth().Time(),
		Location:    p.Location(),
		Relocation:  p.Relocation(),
	}
	fs := make([]gormmodel.ResumeProfileWorkFormat, 0, len(p.WorkFormats()))
	for i, v := range p.WorkFormats() {
		fs = append(fs, gormmodel.ResumeProfileWorkFormat{ProfileID: m.ID, Value: string(v), Position: i})
	}
	return m, fs
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
	photo := domain.RestoreExternalURL(m.PhotoURL)
	dob := domain.RestoreDate(m.DateOfBirth)
	formats := make([]domain.WorkFormat, len(fs))
	for i, v := range fs {
		formats[i] = domain.WorkFormat(v.Value)
	}
	return domain.RestoreProfile(domain.ProfileSnapshot{
		ID: id, ResumeID: rid, Name: m.Name, Surname: m.Surname, Patronymic: m.Patronymic,
		Position: m.Position, PhotoURL: photo, DateOfBirth: dob, Location: m.Location,
		Relocation: m.Relocation, WorkFormats: formats,
	}), nil
}
