package domain

import (
	"fmt"
	"strings"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

// ProjectParams describes requested state. A zero ID creates a new item; a non-zero ID selects an existing item owned by the aggregate.
type ProjectParams struct {
	ID           sharedvalueobjects.ID
	Name         string
	Description  *string
	Role         *string
	Screenshots  []ExternalURL
	Link         *ExternalURL
	GitLink      *ExternalURL
	Stack        []Technology
	Contributors []Contributor
}

type ProjectSnapshot struct {
	ID           sharedvalueobjects.ID
	ResumeID     sharedvalueobjects.ID
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
	id           sharedvalueobjects.ID
	resumeID     sharedvalueobjects.ID
	name         string
	description  *string
	role         *string
	screenshots  []ExternalURL
	link         *ExternalURL
	gitLink      *ExternalURL
	stack        []Technology
	contributors []Contributor
}

func newProject(id, resumeID sharedvalueobjects.ID, params ProjectParams) (*Project, error) {
	stack := make([]Technology, 0, len(params.Stack))
	for _, item := range params.Stack {
		technology, err := NewTechnology(TechnologyParams{Name: item.Name()})
		if err != nil {
			return nil, fmt.Errorf("%w: project contains invalid technology", ErrInvalidEntity)
		}
		stack = append(stack, technology)
	}
	project := &Project{
		id: id, resumeID: resumeID, name: strings.TrimSpace(params.Name),
		description: normalizeOptionalString(params.Description), role: normalizeOptionalString(params.Role),
		screenshots: cloneExternalURLs(params.Screenshots), link: cloneExternalURL(params.Link),
		gitLink: cloneExternalURL(params.GitLink), stack: stack,
		contributors: append([]Contributor(nil), params.Contributors...),
	}
	if err := project.validate(); err != nil {
		return nil, err
	}
	return project, nil
}

func RestoreProject(snapshot ProjectSnapshot) *Project {
	return &Project{
		id: snapshot.ID, resumeID: snapshot.ResumeID, name: snapshot.Name,
		description: copyString(snapshot.Description), role: copyString(snapshot.Role),
		screenshots: cloneExternalURLs(snapshot.Screenshots), link: cloneExternalURL(snapshot.Link),
		gitLink: cloneExternalURL(snapshot.GitLink), stack: append([]Technology(nil), snapshot.Stack...),
		contributors: append([]Contributor(nil), snapshot.Contributors...),
	}
}

func (p *Project) validate() error {
	if p == nil || p.id.IsZero() || p.resumeID.IsZero() || p.name == "" {
		return fmt.Errorf("%w: invalid project identity or required fields", ErrInvalidEntity)
	}
	if !textLengthWithin(p.name, MaxShortTextLength) || !optionalTextLengthWithin(p.description, MaxDescriptionLength) ||
		!optionalTextLengthWithin(p.role, MaxShortTextLength) || len(p.screenshots) > MaxProjectScreenshots ||
		len(p.stack) > MaxProjectStack || len(p.contributors) > MaxProjectContributors {
		return fmt.Errorf("%w: project exceeds configured limits", ErrInvalidEntity)
	}
	if p.link != nil {
		if err := p.link.validate(); err != nil {
			return fmt.Errorf("%w: invalid project URL", ErrInvalidEntity)
		}
	}
	if p.gitLink != nil {
		if err := p.gitLink.validate(); err != nil {
			return fmt.Errorf("%w: invalid project URL", ErrInvalidEntity)
		}
	}
	for _, screenshot := range p.screenshots {
		if err := screenshot.validate(); err != nil {
			return fmt.Errorf("%w: project screenshot URL is invalid", ErrInvalidEntity)
		}
	}
	seenTechnologies := make(map[string]struct{}, len(p.stack))
	for _, technology := range p.stack {
		if err := technology.validate(); err != nil {
			return fmt.Errorf("%w: project contains invalid technology", ErrInvalidEntity)
		}
		if _, exists := seenTechnologies[technology.Name()]; exists {
			return fmt.Errorf("%w: project contains duplicate technology", ErrDuplicateEntity)
		}
		seenTechnologies[technology.Name()] = struct{}{}
	}
	for _, contributor := range p.contributors {
		if err := contributor.validate(); err != nil {
			return fmt.Errorf("%w: project contains invalid contributor", ErrInvalidEntity)
		}
	}
	return nil
}

func (p *Project) ID() sharedvalueobjects.ID       { return p.id }
func (p *Project) ResumeID() sharedvalueobjects.ID { return p.resumeID }
func (p *Project) Name() string                    { return p.name }
func (p *Project) Description() *string            { return copyString(p.description) }
func (p *Project) Role() *string                   { return copyString(p.role) }
func (p *Project) Screenshots() []ExternalURL      { return cloneExternalURLs(p.screenshots) }
func (p *Project) Link() *ExternalURL              { return cloneExternalURL(p.link) }
func (p *Project) GitLink() *ExternalURL           { return cloneExternalURL(p.gitLink) }
func (p *Project) Stack() []Technology {
	return append([]Technology(nil), p.stack...)
}
func (p *Project) Contributors() []Contributor {
	contributors := append([]Contributor(nil), p.contributors...)
	for i := range contributors {
		contributors[i].contacts = cloneContactMethods(contributors[i].contacts)
	}
	return contributors
}

func (r *Resume) ReplaceProjects(params []ProjectParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if len(params) > MaxProjectsPerResume {
		return fmt.Errorf("%w: resume cannot contain more than %d projects", ErrInvalidEntity, MaxProjectsPerResume)
	}
	existing := make(map[sharedvalueobjects.ID]struct{}, len(r.projects))
	for _, current := range r.projects {
		if current == nil {
			return fmt.Errorf("%w: nil project", ErrResumeInvariant)
		}
		existing[current.ID()] = struct{}{}
	}
	seen := make(map[sharedvalueobjects.ID]struct{}, len(params))
	replaced := make([]*Project, 0, len(params))
	for _, item := range params {
		resolvedID, err := resolveCollectionItemID(item.ID, existing, "project")
		if err != nil {
			return err
		}
		if _, ok := seen[resolvedID]; ok {
			return fmt.Errorf("%w: duplicate project %s", ErrInvalidEntity, item.ID.String())
		}
		seen[resolvedID] = struct{}{}
		project, err := newProject(resolvedID, r.id, item)
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
	return cloneProjects(r.projects)
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
	contributor := Contributor{
		name: strings.TrimSpace(params.Name), surname: strings.TrimSpace(params.Surname),
		contacts: params.Contacts, role: strings.TrimSpace(params.Role),
	}
	if err := contributor.validate(); err != nil {
		return Contributor{}, err
	}
	return contributor, nil
}

func RestoreContributor(params ContributorParams) Contributor {
	return Contributor{name: params.Name, surname: params.Surname, contacts: params.Contacts, role: params.Role}
}

func (c Contributor) validate() error {
	if c.name == "" || c.surname == "" || c.role == "" || !textLengthWithin(c.name, MaxProfileNameLength) ||
		!textLengthWithin(c.surname, MaxProfileNameLength) || !textLengthWithin(c.role, MaxProfilePositionLength) {
		return fmt.Errorf("%w: contributor fields are required", ErrInvalidValue)
	}
	if !c.contacts.hasAny() {
		return fmt.Errorf("%w: contributor requires at least one contact", ErrInvalidValue)
	}
	if err := c.contacts.validate(); err != nil {
		return fmt.Errorf("%w: contributor requires valid contacts", ErrInvalidValue)
	}
	return nil
}

func (c Contributor) Name() string             { return c.name }
func (c Contributor) Surname() string          { return c.surname }
func (c Contributor) Role() string             { return c.role }
func (c Contributor) Contacts() ContactMethods { return c.contacts }

func cloneProjects(values []*Project) []*Project {
	if values == nil {
		return nil
	}
	result := make([]*Project, len(values))
	for i, value := range values {
		if value == nil {
			continue
		}
		copy := *value
		copy.description = copyString(value.description)
		copy.role = copyString(value.role)
		copy.screenshots = cloneExternalURLs(value.screenshots)
		copy.link = cloneExternalURL(value.link)
		copy.gitLink = cloneExternalURL(value.gitLink)
		copy.stack = append([]Technology(nil), value.stack...)
		copy.contributors = append([]Contributor(nil), value.contributors...)
		for j := range copy.contributors {
			copy.contributors[j].contacts = cloneContactMethods(value.contributors[j].contacts)
		}
		result[i] = &copy
	}
	return result
}
