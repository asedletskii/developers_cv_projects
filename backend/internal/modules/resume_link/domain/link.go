package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeLinkStatus string

const (
	ResumeLinkStatusActive   ResumeLinkStatus = "active"
	ResumeLinkStatusInactive ResumeLinkStatus = "inactive"
)

func (v ResumeLinkStatus) IsValid() bool {
	switch v {
	case ResumeLinkStatusActive, ResumeLinkStatusInactive:
		return true
	default:
		return false
	}
}

const (
	MaxResumeLinkNameLength = 128
	MaxResumeLinksPerResume = 50
)

type ResumeLink struct {
	id       sharedvalueobjects.ID
	resumeID sharedvalueobjects.ID
	name     string
	status   ResumeLinkStatus
	version  int64
}

type ResumeLinkParams struct {
	ResumeID sharedvalueobjects.ID
	Name     string
}

type ResumeLinkSnapshot struct {
	ID       sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
	Name     string
	Status   ResumeLinkStatus
	Version  int64
}

func NewResumeLink(params ResumeLinkParams) (*ResumeLink, error) {
	id := sharedvalueobjects.NewID()
	name := strings.TrimSpace(params.Name)
	if err := validateResumeLink(id, params.ResumeID, name, ResumeLinkStatusActive, 1); err != nil {
		return nil, err
	}

	return &ResumeLink{
		id:       id,
		resumeID: params.ResumeID,
		name:     name,
		status:   ResumeLinkStatusActive,
		version:  1,
	}, nil
}

func RestoreResumeLink(snapshot ResumeLinkSnapshot) (*ResumeLink, error) {
	if err := validateResumeLink(snapshot.ID, snapshot.ResumeID, snapshot.Name, snapshot.Status, snapshot.Version); err != nil {
		return nil, err
	}
	return &ResumeLink{
		id: snapshot.ID, resumeID: snapshot.ResumeID, name: snapshot.Name,
		status: snapshot.Status, version: snapshot.Version,
	}, nil
}

func EnsureResumeLinkCapacity(existingCount int64) error {
	if existingCount >= MaxResumeLinksPerResume {
		return ErrResumeLinkLimitExceeded
	}
	return nil
}

func ReconstructResumeLink(snapshot ResumeLinkSnapshot) (*ResumeLink, error) {
	return RestoreResumeLink(snapshot)
}

func validateResumeLink(
	id, resumeID sharedvalueobjects.ID,
	name string,
	status ResumeLinkStatus,
	version int64,
) error {
	if id.IsZero() || resumeID.IsZero() {
		return fmt.Errorf("%w: resume link identity is empty", ErrInvalidResumeLink)
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is empty", ErrInvalidResumeLinkName)
	}
	if len([]rune(name)) > MaxResumeLinkNameLength {
		return fmt.Errorf("%w: name exceeds %d characters", ErrInvalidResumeLinkName, MaxResumeLinkNameLength)
	}
	if version < 1 {
		return fmt.Errorf("%w: version must be positive", ErrInvalidResumeLink)
	}
	if !status.IsValid() {
		return fmt.Errorf("%w: invalid status %q", ErrInvalidResumeLink, status)
	}
	return nil
}

func (l *ResumeLink) ID() sharedvalueobjects.ID { return l.id }
func (l *ResumeLink) ResumeID() sharedvalueobjects.ID {
	return l.resumeID
}
func (l *ResumeLink) Name() string             { return l.name }
func (l *ResumeLink) Status() ResumeLinkStatus { return l.status }
func (l *ResumeLink) IsActive() bool           { return l.status == ResumeLinkStatusActive }
func (l *ResumeLink) Version() int64           { return l.version }

func (l *ResumeLink) Rename(name string) error {
	if l == nil {
		return fmt.Errorf("%w: link is nil", ErrInvalidResumeLink)
	}
	trimmed := strings.TrimSpace(name)
	if err := validateResumeLink(l.id, l.resumeID, trimmed, l.status, l.version); err != nil {
		return err
	}
	l.name = trimmed
	return nil
}

func (l *ResumeLink) Activate() error {
	if l == nil {
		return fmt.Errorf("%w: link is nil", ErrInvalidResumeLink)
	}
	l.status = ResumeLinkStatusActive
	return nil
}

func (l *ResumeLink) Deactivate() error {
	if l == nil {
		return fmt.Errorf("%w: link is nil", ErrInvalidResumeLink)
	}
	l.status = ResumeLinkStatusInactive
	return nil
}
