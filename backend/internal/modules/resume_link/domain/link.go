package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type ResumeLinkParams struct {
	ID       valueobjects.ID
	ResumeID valueobjects.ID
	Name     string
	Version  int64
}

const (
	MaxResumeLinkNameLength = 128
	MaxResumeLinksPerResume = 50

	ResumeLinkStatusActive   = "active"
	ResumeLinkStatusInactive = "inactive"
)

type ResumeLink struct {
	id       valueobjects.ID
	resumeID valueobjects.ID
	name     string
	active   bool
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

	return &ResumeLink{
		id:       params.ID,
		resumeID: params.ResumeID,
		name:     name,
		active:   true,
		version:  version,
	}, nil
}

func (l *ResumeLink) ID() valueobjects.ID {
	if l == nil {
		return valueobjects.ID{}
	}
	return l.id
}

func (l *ResumeLink) ResumeID() valueobjects.ID {
	if l == nil {
		return valueobjects.ID{}
	}
	return l.resumeID
}

func (l *ResumeLink) Name() string {
	if l == nil {
		return ""
	}
	return l.name
}

func (l *ResumeLink) IsActive() bool {
	return l != nil && l.active
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
	l.active = true
	return nil
}

func (l *ResumeLink) Deactivate() error {
	if l == nil {
		return fmt.Errorf("%w: link is nil", ErrInvalidResumeLink)
	}
	l.active = false
	return nil
}
