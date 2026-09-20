package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeStatus string

const (
	ResumeStatusDraft  ResumeStatus = "draft"
	ResumeStatusOpened ResumeStatus = "opened"
	ResumeStatusClosed ResumeStatus = "closed"
)

func (v ResumeStatus) isValid() bool {
	switch v {
	case ResumeStatusDraft, ResumeStatusOpened, ResumeStatusClosed:
		return true
	default:
		return false
	}
}

func ParseResumeStatus(value string) (ResumeStatus, error) {
	status := ResumeStatus(value)
	if !status.isValid() {
		return "", fmt.Errorf("%w: invalid resume status", ErrInvalidValue)
	}
	return status, nil
}

const MaxResumeTitleLength = 200

type ResumeLanguage string

const (
	ResumeLanguageRU ResumeLanguage = "ru"
	ResumeLanguageEN ResumeLanguage = "en"
)

func (v ResumeLanguage) isValid() bool {
	switch v {
	case ResumeLanguageRU, ResumeLanguageEN:
		return true
	default:
		return false
	}
}

type Resume struct {
	id              sharedvalueobjects.ID
	userID          sharedvalueobjects.ID
	title           string
	language        ResumeLanguage
	status          ResumeStatus
	version         int64
	overview        *ResumeOverview
	profile         *Profile
	contacts        *Contacts
	skills          *Skills
	spokenLanguages *ResumeSpokenLanguages
	technologies    *ResumeTechnologies
	experiences     []*Experience
	educations      []*Education
	courses         []*Course
	activities      []*Activity
	projects        []*Project
}

type ResumeParams struct {
	UserID   sharedvalueobjects.ID
	Title    string
	Language ResumeLanguage
}

type ResumeSnapshot struct {
	ID              sharedvalueobjects.ID
	UserID          sharedvalueobjects.ID
	Title           string
	Language        ResumeLanguage
	Status          ResumeStatus
	Version         int64
	Overview        *ResumeOverview
	Profile         *Profile
	Contacts        *Contacts
	Skills          *Skills
	SpokenLanguages *ResumeSpokenLanguages
	Technologies    *ResumeTechnologies
	Experiences     []*Experience
	Educations      []*Education
	Courses         []*Course
	Activities      []*Activity
	Projects        []*Project
}

func NewResume(params ResumeParams) (*Resume, error) {
	resume := &Resume{
		id: sharedvalueobjects.NewID(), userID: params.UserID, title: strings.TrimSpace(params.Title),
		language: params.Language, status: ResumeStatusDraft, version: 1,
	}
	if err := resume.validate(); err != nil {
		return nil, err
	}
	return resume, nil
}

func (r *Resume) validate() error {
	if r == nil || r.id.IsZero() || r.userID.IsZero() {
		return fmt.Errorf("%w: invalid resume identity", ErrInvalidResume)
	}
	if !r.language.isValid() {
		return fmt.Errorf("%w: invalid resume language", ErrInvalidResume)
	}
	if err := validateResumeTitle(r.title); err != nil {
		return err
	}
	if !r.status.isValid() {
		return fmt.Errorf("%w: invalid resume status", ErrInvalidResume)
	}
	if r.version < 1 {
		return fmt.Errorf("%w: resume version must be positive", ErrInvalidResume)
	}
	return nil
}

func RestoreResume(snapshot ResumeSnapshot) *Resume {
	return &Resume{
		id:              snapshot.ID,
		userID:          snapshot.UserID,
		title:           snapshot.Title,
		language:        snapshot.Language,
		status:          snapshot.Status,
		version:         snapshot.Version,
		overview:        cloneResumeOverview(snapshot.Overview),
		profile:         cloneProfile(snapshot.Profile),
		contacts:        cloneContacts(snapshot.Contacts),
		skills:          cloneSkills(snapshot.Skills),
		spokenLanguages: cloneResumeSpokenLanguages(snapshot.SpokenLanguages),
		technologies:    cloneResumeTechnologies(snapshot.Technologies),
		experiences:     cloneExperiences(snapshot.Experiences),
		educations:      cloneEducations(snapshot.Educations),
		courses:         cloneCourses(snapshot.Courses),
		activities:      cloneActivities(snapshot.Activities),
		projects:        cloneProjects(snapshot.Projects),
	}
}

func (r *Resume) ID() sharedvalueobjects.ID     { return r.id }
func (r *Resume) UserID() sharedvalueobjects.ID { return r.userID }
func (r *Resume) Title() string                 { return r.title }
func (r *Resume) Language() ResumeLanguage      { return r.language }
func (r *Resume) Version() int64                { return r.version }
func (r *Resume) Status() ResumeStatus          { return r.status }

func (r *Resume) ReplaceStatus(status ResumeStatus) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if err := validateResumeStatusTransition(r.status, status); err != nil {
		return err
	}
	r.status = status
	return nil
}

func (r *Resume) ReplaceTitle(title string) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	normalized, err := normalizeResumeTitle(title)
	if err != nil {
		return err
	}
	r.title = normalized
	return nil
}

func normalizeResumeTitle(title string) (string, error) {
	normalized := strings.TrimSpace(title)
	if err := validateResumeTitle(normalized); err != nil {
		return "", err
	}
	return normalized, nil
}

func validateResumeTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("%w: resume title is empty", ErrInvalidResume)
	}
	if utf8.RuneCountInString(title) > MaxResumeTitleLength {
		return fmt.Errorf("%w: resume title exceeds %d characters", ErrInvalidResume, MaxResumeTitleLength)
	}
	return nil
}

func validateResumeStatusTransition(current, target ResumeStatus) error {
	switch target {
	case ResumeStatusOpened:
		if current == ResumeStatusOpened {
			return fmt.Errorf("%w: resume is already opened", ErrInvalidResumeStatusTransition)
		}
		if current != ResumeStatusDraft && current != ResumeStatusClosed {
			return fmt.Errorf("%w: unsupported current status %q", ErrInvalidResumeStatusTransition, current)
		}
	case ResumeStatusClosed:
		if current != ResumeStatusOpened {
			return fmt.Errorf("%w: cannot close resume with status %q", ErrInvalidResumeStatusTransition, current)
		}
	default:
		return fmt.Errorf("%w: unsupported target status %q", ErrInvalidResumeStatusTransition, target)
	}
	return nil
}
