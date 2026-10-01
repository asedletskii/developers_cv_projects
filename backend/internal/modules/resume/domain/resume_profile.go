package domain

import (
	"fmt"
	"strings"
	"time"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type WorkFormat string

const (
	WorkFormatFullTime WorkFormat = "fullTime"
	WorkFormatPartTime WorkFormat = "partTime"
	WorkFormatOffice   WorkFormat = "office"
	WorkFormatRemote   WorkFormat = "remote"
	WorkFormatHybrid   WorkFormat = "hybrid"
)

func (v WorkFormat) IsValid() bool {
	switch v {
	case WorkFormatFullTime, WorkFormatPartTime, WorkFormatOffice, WorkFormatRemote, WorkFormatHybrid:
		return true
	default:
		return false
	}
}

type Profile struct {
	id          identity.ID
	resumeID    identity.ID
	name        string
	surname     string
	patronymic  *string
	position    string
	photoPath   ExternalURL
	dateOfBirth Date
	location    string
	relocation  bool
	workFormats []WorkFormat
}

type ProfileParams struct {
	ID          identity.ID
	ResumeID    identity.ID
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoPath   ExternalURL
	DateOfBirth Date
	Location    string
	Relocation  bool
	WorkFormats []WorkFormat
}

func NewProfile(p ProfileParams) (*Profile, error) {
	name, surname, position, location := strings.TrimSpace(p.Name), strings.TrimSpace(p.Surname), strings.TrimSpace(p.Position), strings.TrimSpace(p.Location)
	if p.ID.IsZero() || p.ResumeID.IsZero() || name == "" || surname == "" || position == "" || location == "" || p.PhotoPath.IsZero() || p.DateOfBirth.IsZero() {
		return nil, fmt.Errorf("%w: required profile fields are empty", ErrInvalidEntity)
	}
	if p.DateOfBirth.Time().After(time.Now().UTC()) {
		return nil, fmt.Errorf("%w: date of birth is in the future", ErrInvalidEntity)
	}
	formats := make([]WorkFormat, len(p.WorkFormats))
	seen := make(map[WorkFormat]struct{}, len(p.WorkFormats))
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
func (p *Profile) Replace(params ProfileParams) error {
	params.ID, params.ResumeID = p.id, p.resumeID
	updated, err := NewProfile(params)
	if err != nil {
		return err
	}
	*p = *updated
	return nil
}

func (p *Profile) ID() identity.ID       { return p.id }
func (p *Profile) ResumeID() identity.ID { return p.resumeID }
func (p *Profile) Name() string          { return p.name }
func (p *Profile) Surname() string       { return p.surname }
func (p *Profile) Patronymic() *string {
	if p.patronymic == nil {
		return nil
	}
	v := *p.patronymic
	return &v
}
func (p *Profile) Position() string       { return p.position }
func (p *Profile) PhotoPath() ExternalURL { return p.photoPath }
func (p *Profile) DateOfBirth() Date      { return p.dateOfBirth }
func (p *Profile) Location() string       { return p.location }
func (p *Profile) Relocation() bool       { return p.relocation }
func (p *Profile) WorkFormats() []WorkFormat {
	return append([]WorkFormat(nil), p.workFormats...)
}

func (r *Resume) ReplaceProfile(params ProfileParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.profile != nil {
			params.ID = r.profile.ID()
		} else {
			params.ID = identity.NewID()
		}
	} else if r.profile == nil || params.ID != r.profile.ID() {
		return fmt.Errorf("%w: profile does not belong to resume", ErrResumeInvariant)
	}
	params.ResumeID = r.id
	profile, err := NewProfile(params)
	if err != nil {
		return err
	}
	r.profile = profile
	return nil
}

func (r *Resume) Profile() *Profile {
	if r == nil || r.profile == nil {
		return nil
	}
	copy := *r.profile
	return &copy
}
