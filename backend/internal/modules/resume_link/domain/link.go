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

type ResumeLinkParams struct {
	ID       sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
	Name     string
	Version  int64
	Status   ResumeLinkStatus
}

type ResumeLink struct {
	id       sharedvalueobjects.ID
	resumeID sharedvalueobjects.ID
	name     string
	status   ResumeLinkStatus
	version  int64
}

func NewResumeLink(params ResumeLinkParams) (*ResumeLink, error) {
	name := strings.TrimSpace(params.Name)
	if params.ID.IsZero() || params.ResumeID.IsZero() {
		return nil, fmt.Errorf("%w: identity is empty", ErrInvalidResumeLink)
	}
	if name == "" {
		return nil, fmt.Errorf("%w: name is empty", ErrInvalidResumeLinkName)
	}
	if len([]rune(name)) > MaxResumeLinkNameLength {
		return nil, fmt.Errorf("%w: name exceeds %d characters", ErrInvalidResumeLinkName, MaxResumeLinkNameLength)
	}

	version := params.Version
	if version == 0 {
		version = 1
	}
	if version < 1 {
		return nil, fmt.Errorf("%w: version must be positive", ErrInvalidResumeLink)
	}

	status := params.Status
	if status == "" {
		status = ResumeLinkStatusActive
	}
	if !status.IsValid() {
		return nil, fmt.Errorf("%w: invalid status %q", ErrInvalidResumeLink, status)
	}

	return &ResumeLink{
		id:       params.ID,
		resumeID: params.ResumeID,
		name:     name,
		status:   status,
		version:  version,
	}, nil
}

func (l *ResumeLink) ID() sharedvalueobjects.ID {
	if l == nil {
		return sharedvalueobjects.ID{}
	}
	return l.id
}

func (l *ResumeLink) ResumeID() sharedvalueobjects.ID {
	if l == nil {
		return sharedvalueobjects.ID{}
	}
	return l.resumeID
}

func (l *ResumeLink) Name() string {
	if l == nil {
		return ""
	}
	return l.name
}

func (l *ResumeLink) Status() ResumeLinkStatus {
	if l == nil {
		return ""
	}
	return l.status
}

func (l *ResumeLink) IsActive() bool {
	return l != nil && l.status == ResumeLinkStatusActive
}

func (l *ResumeLink) Version() int64 {
	if l == nil {
		return 0
	}
	return l.version
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
