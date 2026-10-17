package domain

import (
	"fmt"
	"strings"
	"time"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type WorkFormat string

const (
	WorkFormatFullTime WorkFormat = "fullTime"
	WorkFormatPartTime WorkFormat = "partTime"
	WorkFormatOffice   WorkFormat = "office"
	WorkFormatRemote   WorkFormat = "remote"
	WorkFormatHybrid   WorkFormat = "hybrid"
)

func (v WorkFormat) isValid() bool {
	switch v {
	case WorkFormatFullTime, WorkFormatPartTime, WorkFormatOffice, WorkFormatRemote, WorkFormatHybrid:
		return true
	default:
		return false
	}
}

type Profile struct {
	id          sharedvalueobjects.ID
	resumeID    sharedvalueobjects.ID
	name        string
	surname     string
	patronymic  *string
	position    string
	photoURL    ExternalURL
	dateOfBirth Date
	location    string
	relocation  bool
	workFormats []WorkFormat
}

type ProfileParams struct {
	ID          sharedvalueobjects.ID
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoURL    ExternalURL
	DateOfBirth Date
	Location    string
	Relocation  bool
	WorkFormats []WorkFormat
}

type ProfileSnapshot struct {
	ID          sharedvalueobjects.ID
	ResumeID    sharedvalueobjects.ID
	Name        string
	Surname     string
	Patronymic  *string
	Position    string
	PhotoURL    ExternalURL
	DateOfBirth Date
	Location    string
	Relocation  bool
	WorkFormats []WorkFormat
}

func newProfile(params ProfileParams, resumeID sharedvalueobjects.ID) (*Profile, error) {
	profile := &Profile{
		id: params.ID, resumeID: resumeID, name: strings.TrimSpace(params.Name),
		surname: strings.TrimSpace(params.Surname), patronymic: normalizeOptionalString(params.Patronymic),
		position: strings.TrimSpace(params.Position), photoURL: params.PhotoURL, dateOfBirth: params.DateOfBirth,
		location: strings.TrimSpace(params.Location), relocation: params.Relocation,
		workFormats: append([]WorkFormat(nil), params.WorkFormats...),
	}
	if err := profile.validate(); err != nil {
		return nil, err
	}
	return profile, nil
}

func RestoreProfile(snapshot ProfileSnapshot) *Profile {
	return &Profile{
		id: snapshot.ID, resumeID: snapshot.ResumeID, name: snapshot.Name, surname: snapshot.Surname,
		patronymic: copyString(snapshot.Patronymic), position: snapshot.Position, photoURL: snapshot.PhotoURL,
		dateOfBirth: snapshot.DateOfBirth, location: snapshot.Location, relocation: snapshot.Relocation,
		workFormats: append([]WorkFormat(nil), snapshot.WorkFormats...),
	}
}

func (p *Profile) validate() error {
	if p == nil || p.id.IsZero() || p.resumeID.IsZero() || p.name == "" || p.surname == "" ||
		p.position == "" || p.location == "" || p.photoURL.IsZero() || p.dateOfBirth.IsZero() {
		return fmt.Errorf("%w: required profile fields are empty", ErrInvalidEntity)
	}
	if !textLengthWithin(p.name, MaxProfileNameLength) || !textLengthWithin(p.surname, MaxProfileNameLength) ||
		!optionalTextLengthWithin(p.patronymic, MaxProfileNameLength) || !textLengthWithin(p.position, MaxProfilePositionLength) ||
		!textLengthWithin(p.location, MaxLocationLength) {
		return fmt.Errorf("%w: profile text exceeds configured limits", ErrInvalidEntity)
	}
	if err := p.photoURL.validate(); err != nil {
		return fmt.Errorf("%w: invalid profile photo URL", ErrInvalidEntity)
	}
	if p.dateOfBirth.Time().After(time.Now().UTC()) {
		return fmt.Errorf("%w: date of birth is in the future", ErrInvalidEntity)
	}
	seen := make(map[WorkFormat]struct{}, len(p.workFormats))
	for _, format := range p.workFormats {
		if !format.isValid() {
			return fmt.Errorf("%w: invalid work format", ErrInvalidEntity)
		}
		if _, exists := seen[format]; exists {
			return fmt.Errorf("%w: duplicate work format", ErrDuplicateEntity)
		}
		seen[format] = struct{}{}
	}
	return nil
}

func (p *Profile) ID() sharedvalueobjects.ID       { return p.id }
func (p *Profile) ResumeID() sharedvalueobjects.ID { return p.resumeID }
func (p *Profile) Name() string                    { return p.name }
func (p *Profile) Surname() string                 { return p.surname }
func (p *Profile) Patronymic() *string             { return copyString(p.patronymic) }
func (p *Profile) Position() string                { return p.position }
func (p *Profile) PhotoURL() ExternalURL           { return p.photoURL }
func (p *Profile) DateOfBirth() Date               { return p.dateOfBirth }
func (p *Profile) Location() string                { return p.location }
func (p *Profile) Relocation() bool                { return p.relocation }
func (p *Profile) WorkFormats() []WorkFormat {
	return append([]WorkFormat(nil), p.workFormats...)
}

func (r *Resume) ReplaceProfile(params ProfileParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	var currentID sharedvalueobjects.ID
	if r.profile != nil {
		currentID = r.profile.ID()
	}
	resolvedID, err := resolveSingletonID(params.ID, currentID, "profile")
	if err != nil {
		return err
	}
	params.ID = resolvedID
	profile, err := newProfile(params, r.id)
	if err != nil {
		return err
	}
	r.profile = profile
	return nil
}

func (r *Resume) Profile() *Profile {
	if r == nil {
		return nil
	}
	return cloneProfile(r.profile)
}

func cloneProfile(value *Profile) *Profile {
	if value == nil {
		return nil
	}
	copy := *value
	copy.patronymic = copyString(value.patronymic)
	copy.workFormats = append([]WorkFormat(nil), value.workFormats...)
	return &copy
}
