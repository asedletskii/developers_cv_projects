package entities

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type ProjectParams struct {
	ID           valueobjects.ID
	Name         string
	Description  *string
	Role         *string
	Screenshots  []valueobjects.AssetPath
	Link         *valueobjects.ExternalURL
	GitLink      *valueobjects.ExternalURL
	Stack        []Technology
	Contributors []Contributor
}

type Project struct {
	id           valueobjects.ID
	name         string
	description  *string
	role         *string
	screenshots  []valueobjects.AssetPath
	link         *valueobjects.ExternalURL
	gitLink      *valueobjects.ExternalURL
	stack        []Technology
	contributors []Contributor
}

func NewProject(params ProjectParams) (Project, error) {
	name := strings.TrimSpace(params.Name)
	if params.ID.IsZero() || name == "" {
		return Project{}, fmt.Errorf("%w: invalid project fields", valueobjects.ErrInvalidValue)
	}
	return Project{
		id:           params.ID,
		name:         name,
		description:  cloneString(params.Description),
		role:         cloneString(params.Role),
		screenshots:  cloneAssetPaths(params.Screenshots),
		link:         cloneExternalURL(params.Link),
		gitLink:      cloneExternalURL(params.GitLink),
		stack:        cloneTechnologies(params.Stack),
		contributors: cloneContributors(params.Contributors),
	}, nil
}

func (p Project) ID() valueobjects.ID                   { return p.id }
func (p Project) Name() string                          { return p.name }
func (p Project) Description() *string                  { return cloneString(p.description) }
func (p Project) Role() *string                         { return cloneString(p.role) }
func (p Project) Screenshots() []valueobjects.AssetPath { return cloneAssetPaths(p.screenshots) }
func (p Project) Link() *valueobjects.ExternalURL       { return cloneExternalURL(p.link) }
func (p Project) GitLink() *valueobjects.ExternalURL    { return cloneExternalURL(p.gitLink) }
func (p Project) Stack() []Technology                   { return cloneTechnologies(p.stack) }
func (p Project) Contributors() []Contributor           { return cloneContributors(p.contributors) }
