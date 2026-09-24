package domain

import (
	"fmt"
	"strings"
	"time"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type Profile struct {
	id          valueobjects.ID
	resumeID    valueobjects.ID
	name        string
	surname     string
	patronymic  *string
	position    string
	photoPath   valueobjects.AssetPath
	dateOfBirth valueobjects.Date
	location    string
	relocation  bool
	workFormats []enums.WorkFormat
}

type ProfileParams struct {
	ID          valueobjects.ID
	ResumeID    valueobjects.ID
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoPath   valueobjects.AssetPath
	DateOfBirth valueobjects.Date
	Location    string
	Relocation  bool
	WorkFormats []enums.WorkFormat
}

func NewProfile(p ProfileParams) (*Profile, error) {
	name, surname, position, location := strings.TrimSpace(p.Name), strings.TrimSpace(p.Surname), strings.TrimSpace(p.Position), strings.TrimSpace(p.Location)
	if p.ID.IsZero() || p.ResumeID.IsZero() || name == "" || surname == "" || position == "" || location == "" || p.PhotoPath.IsZero() || p.DateOfBirth.IsZero() {
		return nil, fmt.Errorf("%w: required profile fields are empty", ErrInvalidEntity)
	}
	if p.DateOfBirth.Time().After(time.Now().UTC()) {
		return nil, fmt.Errorf("%w: date of birth is in the future", ErrInvalidEntity)
	}
	formats := make([]enums.WorkFormat, len(p.WorkFormats))
	seen := make(map[enums.WorkFormat]struct{}, len(p.WorkFormats))
	for i, f := range p.WorkFormats {
		if !f.IsValid() {
			return nil, fmt.Errorf("%w: invalid work format", ErrInvalidEntity)
		}
		if _, ok := seen[f]; ok {
			return nil, fmt.Errorf("%w: duplicate work format", ErrDuplicateEntity)
		}
		seen[f] = struct{}{}
		formats[i] = f
	}
	var patronymic *string
	if p.Patronymic != nil {
		v := strings.TrimSpace(*p.Patronymic)
		if v != "" {
			patronymic = &v
		}
	}
	return &Profile{id: p.ID, resumeID: p.ResumeID, name: name, surname: surname, patronymic: patronymic, position: position, photoPath: p.PhotoPath, dateOfBirth: p.DateOfBirth, location: location, relocation: p.Relocation, workFormats: formats}, nil
}
func (p *Profile) ID() valueobjects.ID       { return p.id }
func (p *Profile) ResumeID() valueobjects.ID { return p.resumeID }
func (p *Profile) Name() string              { return p.name }
func (p *Profile) Surname() string           { return p.surname }
func (p *Profile) Patronymic() *string {
	if p.patronymic == nil {
		return nil
	}
	v := *p.patronymic
	return &v
}
func (p *Profile) Position() string                  { return p.position }
func (p *Profile) PhotoPath() valueobjects.AssetPath { return p.photoPath }
func (p *Profile) DateOfBirth() valueobjects.Date    { return p.dateOfBirth }
func (p *Profile) Location() string                  { return p.location }
func (p *Profile) Relocation() bool                  { return p.relocation }
func (p *Profile) WorkFormats() []enums.WorkFormat {
	return append([]enums.WorkFormat(nil), p.workFormats...)
}

type ProfileOwner struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
