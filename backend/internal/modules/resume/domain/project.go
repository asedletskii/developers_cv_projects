package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ProjectParams struct {
	ID           valueobjects.ID
	ResumeID     valueobjects.ID
	Name         string
	Description  *string
	Role         *string
	Screenshots  []valueobjects.AssetPath
	Link         *valueobjects.ExternalURL
	GitLink      *valueobjects.ExternalURL
	Stack        []entities.Technology
	Contributors []entities.Contributor
}

type Project struct {
	id           valueobjects.ID
	resumeID     valueobjects.ID
	name         string
	description  *string
	role         *string
	screenshots  []valueobjects.AssetPath
	link         *valueobjects.ExternalURL
	gitLink      *valueobjects.ExternalURL
	stack        []entities.Technology
	contributors []entities.Contributor
}

func NewProject(params ProjectParams) (*Project, error) {
	name := strings.TrimSpace(params.Name)
	if params.ID.IsZero() || params.ResumeID.IsZero() || name == "" {
		return nil, fmt.Errorf("%w: invalid project identity or required fields", ErrInvalidEntity)
	}
	return &Project{
		id:           params.ID,
		resumeID:     params.ResumeID,
		name:         name,
		description:  cloneString(params.Description),
		role:         cloneString(params.Role),
		screenshots:  cloneAssetPaths(params.Screenshots),
		link:         cloneExternalURL(params.Link),
		gitLink:      cloneExternalURL(params.GitLink),
		stack:        append([]entities.Technology(nil), params.Stack...),
		contributors: append([]entities.Contributor(nil), params.Contributors...),
	}, nil
}

func (p *Project) ID() valueobjects.ID                   { return p.id }
func (p *Project) ResumeID() valueobjects.ID             { return p.resumeID }
func (p *Project) Name() string                          { return p.name }
func (p *Project) Description() *string                  { return cloneString(p.description) }
func (p *Project) Role() *string                         { return cloneString(p.role) }
func (p *Project) Screenshots() []valueobjects.AssetPath { return cloneAssetPaths(p.screenshots) }
func (p *Project) Link() *valueobjects.ExternalURL       { return cloneExternalURL(p.link) }
func (p *Project) GitLink() *valueobjects.ExternalURL    { return cloneExternalURL(p.gitLink) }
func (p *Project) Stack() []entities.Technology {
	return append([]entities.Technology(nil), p.stack...)
}
func (p *Project) Contributors() []entities.Contributor {
	return append([]entities.Contributor(nil), p.contributors...)
}
