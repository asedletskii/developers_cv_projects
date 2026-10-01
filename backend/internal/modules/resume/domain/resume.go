package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeStatus string

const (
	ResumeStatusDraft  ResumeStatus = "draft"
	ResumeStatusOpened ResumeStatus = "opened"
	ResumeStatusClosed ResumeStatus = "closed"
)

func (v ResumeStatus) IsValid() bool {
	switch v {
	case ResumeStatusDraft, ResumeStatusOpened, ResumeStatusClosed:
		return true
	default:
		return false
	}
}

type ResumeLanguage string

const (
	ResumeLanguageRU ResumeLanguage = "ru"
	ResumeLanguageEN ResumeLanguage = "en"
	ResumeLanguageFR ResumeLanguage = "fr"
	ResumeLanguageDE ResumeLanguage = "de"
	ResumeLanguageKZ ResumeLanguage = "kz"
	ResumeLanguageZH ResumeLanguage = "zh"
	ResumeLanguageES ResumeLanguage = "es"
)

func (v ResumeLanguage) IsValid() bool {
	switch v {
	case ResumeLanguageRU, ResumeLanguageEN, ResumeLanguageFR, ResumeLanguageDE, ResumeLanguageKZ, ResumeLanguageZH, ResumeLanguageES:
		return true
	default:
		return false
	}
}

type Resume struct {
	id           sharedvalueobjects.ID
	userID       sharedvalueobjects.ID
	title        string
	language     ResumeLanguage
	status       ResumeStatus
	version      int64
	overview     *ResumeOverview
	profile      *Profile
	contacts     *Contacts
	skills       *Skills
	languages    *ResumeLanguages
	technologies *ResumeTechnologies
	experiences  []*Experience
	educations   []*Education
	courses      []*Course
	activities   []*Activity
	projects     []*Project
}

type ResumeParams struct {
	ID       sharedvalueobjects.ID
	UserID   sharedvalueobjects.ID
	Title    string
	Language ResumeLanguage
	Version  int64
}

type ResumeSnapshot struct {
	ID           sharedvalueobjects.ID
	UserID       sharedvalueobjects.ID
	Title        string
	Language     ResumeLanguage
	Status       ResumeStatus
	Version      int64
	Overview     *ResumeOverview
	Profile      *Profile
	Contacts     *Contacts
	Skills       *Skills
	Languages    *ResumeLanguages
	Technologies *ResumeTechnologies
	Experiences  []*Experience
	Educations   []*Education
	Courses      []*Course
	Activities   []*Activity
	Projects     []*Project
}

func NewResume(params ResumeParams) (*Resume, error) {
	title := strings.TrimSpace(params.Title)
	if params.ID.IsZero() {
		params.ID = sharedvalueobjects.NewID()
	}
	if params.UserID.IsZero() || title == "" {
		return nil, fmt.Errorf("%w: invalid resume identity or title", ErrInvalidResume)
	}
	if !params.Language.IsValid() {
		return nil, fmt.Errorf("%w: invalid resume language", ErrInvalidResume)
	}
	if err := ValidateResumeTitle(title); err != nil {
		return nil, err
	}
	version := params.Version
	if version == 0 {
		version = 1
	}
	if version < 1 {
		return nil, fmt.Errorf("%w: resume version must be positive", ErrInvalidResume)
	}
	return &Resume{id: params.ID, userID: params.UserID, title: title, language: params.Language, status: ResumeStatusDraft, version: version}, nil
}

func ReconstructResume(snapshot ResumeSnapshot) (*Resume, error) {
	resume, err := NewResume(ResumeParams{
		ID: snapshot.ID, UserID: snapshot.UserID, Title: snapshot.Title,
		Language: snapshot.Language, Version: snapshot.Version,
	})
	if err != nil {
		return nil, err
	}
	if !snapshot.Status.IsValid() {
		return nil, fmt.Errorf("%w: invalid resume status", ErrInvalidResume)
	}
	resume.status = snapshot.Status
	if err := resume.validateOwnedState(snapshot); err != nil {
		return nil, err
	}
	resume.overview = snapshot.Overview
	resume.profile = snapshot.Profile
	resume.contacts = snapshot.Contacts
	resume.skills = snapshot.Skills
	resume.languages = snapshot.Languages
	resume.technologies = snapshot.Technologies
	resume.experiences = append([]*Experience(nil), snapshot.Experiences...)
	resume.educations = append([]*Education(nil), snapshot.Educations...)
	resume.courses = append([]*Course(nil), snapshot.Courses...)
	resume.activities = append([]*Activity(nil), snapshot.Activities...)
	resume.projects = append([]*Project(nil), snapshot.Projects...)
	return resume, nil
}

func (r *Resume) validateOwnedState(snapshot ResumeSnapshot) error {
	check := func(id sharedvalueobjects.ID) error {
		if id.IsZero() || id != r.id {
			return fmt.Errorf("%w: child does not belong to resume", ErrResumeInvariant)
		}
		return nil
	}
	for _, item := range snapshot.Experiences {
		if item == nil {
			return fmt.Errorf("%w: nil experience", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Educations {
		if item == nil {
			return fmt.Errorf("%w: nil education", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Courses {
		if item == nil {
			return fmt.Errorf("%w: nil course", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Activities {
		if item == nil {
			return fmt.Errorf("%w: nil activity", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Projects {
		if item == nil {
			return fmt.Errorf("%w: nil project", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Profile != nil {
		if err := check(snapshot.Profile.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Contacts != nil {
		if err := check(snapshot.Contacts.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Skills != nil {
		if err := check(snapshot.Skills.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Languages != nil {
		if err := check(snapshot.Languages.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Technologies != nil {
		if err := check(snapshot.Technologies.ResumeID()); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resume) ID() sharedvalueobjects.ID     { return r.id }
func (r *Resume) UserID() sharedvalueobjects.ID { return r.userID }
func (r *Resume) Title() string                 { return r.title }
func (r *Resume) Language() ResumeLanguage      { return r.language }
func (r *Resume) Version() int64                { return r.version }
func (r *Resume) Status() ResumeStatus          { return r.status }

func (r *Resume) Publish() error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if r.status == ResumeStatusOpened {
		return fmt.Errorf("%w: resume is already opened", ErrInvalidResumeStatusTransition)
	}
	if r.status != ResumeStatusDraft && r.status != ResumeStatusClosed {
		return fmt.Errorf("%w: unsupported current status %q", ErrInvalidResumeStatusTransition, r.status)
	}
	r.status = ResumeStatusOpened
	return nil
}

func (r *Resume) Close() error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if r.status != ResumeStatusOpened {
		return fmt.Errorf("%w: cannot close resume with status %q", ErrInvalidResumeStatusTransition, r.status)
	}
	r.status = ResumeStatusClosed
	return nil
}

func (r *Resume) ReplaceTitle(title string) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	title = strings.TrimSpace(title)
	if err := ValidateResumeTitle(title); err != nil {
		return err
	}
	r.title = title
	return nil
}

func ValidateResumeTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w: resume title is empty", ErrInvalidResume)
	}
	return nil
}
