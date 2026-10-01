package domain

import (
	"fmt"
	"strings"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type ProjectParams struct {
	ID           identity.ID
	ResumeID     identity.ID
	Name         string
	Description  *string
	Role         *string
	Screenshots  []ExternalURL
	Link         *ExternalURL
	GitLink      *ExternalURL
	Stack        []Technology
	Contributors []Contributor
}

type Project struct {
	id           identity.ID
	resumeID     identity.ID
	name         string
	description  *string
	role         *string
	screenshots  []ExternalURL
	link         *ExternalURL
	gitLink      *ExternalURL
	stack        []Technology
	contributors []Contributor
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
		screenshots:  cloneExternalURLs(params.Screenshots),
		link:         cloneExternalURL(params.Link),
		gitLink:      cloneExternalURL(params.GitLink),
		stack:        append([]Technology(nil), params.Stack...),
		contributors: append([]Contributor(nil), params.Contributors...),
	}, nil
}

func (p *Project) Replace(params ProjectParams) error {
	params.ID, params.ResumeID = p.id, p.resumeID
	updated, err := NewProject(params)
	if err != nil {
		return err
	}
	*p = *updated
	return nil
}

func (p *Project) ID() identity.ID            { return p.id }
func (p *Project) ResumeID() identity.ID      { return p.resumeID }
func (p *Project) Name() string               { return p.name }
func (p *Project) Description() *string       { return cloneString(p.description) }
func (p *Project) Role() *string              { return cloneString(p.role) }
func (p *Project) Screenshots() []ExternalURL { return cloneExternalURLs(p.screenshots) }
func (p *Project) Link() *ExternalURL         { return cloneExternalURL(p.link) }
func (p *Project) GitLink() *ExternalURL      { return cloneExternalURL(p.gitLink) }
func (p *Project) Stack() []Technology {
	return append([]Technology(nil), p.stack...)
}
func (p *Project) Contributors() []Contributor {
	return append([]Contributor(nil), p.contributors...)
}

func (r *Resume) ReplaceProjects(params []ProjectParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[identity.ID]struct{}, len(r.projects))
	for _, item := range r.projects {
		if item == nil {
			return fmt.Errorf("%w: nil project", ErrResumeInvariant)
		}
		existing[item.ID()] = struct{}{}
	}
	seen := make(map[identity.ID]struct{}, len(params))
	replaced := make([]*Project, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = identity.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: project %s does not belong to resume", ErrResumeInvariant, item.ID.String())
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate project %s", ErrInvalidEntity, item.ID.String())
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		project, err := NewProject(item)
		if err != nil {
			return err
		}
		replaced = append(replaced, project)
	}
	r.projects = replaced
	return nil
}

func (r *Resume) Projects() []*Project {
	if r == nil {
		return nil
	}
	return append([]*Project(nil), r.projects...)
}

type ContributorParams struct {
	Name     string
	Surname  string
	Contacts ContactMethods
	Role     string
}

type Contributor struct {
	name     string
	surname  string
	contacts ContactMethods
	role     string
}

func NewContributor(params ContributorParams) (Contributor, error) {
	name := strings.TrimSpace(params.Name)
	surname := strings.TrimSpace(params.Surname)
	role := strings.TrimSpace(params.Role)
	if name == "" || surname == "" || role == "" {
		return Contributor{}, fmt.Errorf("%w: contributor fields are required", ErrInvalidValue)
	}
	if !params.Contacts.HasAny() {
		return Contributor{}, fmt.Errorf("%w: contributor requires at least one contact", ErrInvalidValue)
	}
	return Contributor{name: name, surname: surname, contacts: params.Contacts, role: role}, nil
}

func (c Contributor) Name() string             { return c.name }
func (c Contributor) Surname() string          { return c.surname }
func (c Contributor) Role() string             { return c.role }
func (c Contributor) Contacts() ContactMethods { return c.contacts }
