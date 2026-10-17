package analytics

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
)

type LayoutSnapshotBuilder struct{}

func NewLayoutSnapshotBuilder() *LayoutSnapshotBuilder { return &LayoutSnapshotBuilder{} }

func (b *LayoutSnapshotBuilder) Build(resume resumequeries.PublicResumeView) (domain.LayoutSnapshot, error) {
	sections := resume.Overview.Sections
	if len(sections) == 0 {
		sections = defaultSections()
	}
	ordered := make([]resumequeries.PublicSectionSettings, 0, len(sections)+1)
	ordered = append(ordered, resumequeries.PublicSectionSettings{ID: "about", Enabled: true})
	for _, section := range sections {
		if section.ID != "about" {
			ordered = append(ordered, section)
		}
	}
	order := make([]resumedomain.SectionID, 0, len(ordered))
	hashes := make(map[resumedomain.SectionID]string)
	seen := make(map[resumedomain.SectionID]struct{})
	for _, section := range ordered {
		id := resumedomain.SectionID(section.ID)
		if !section.Enabled || id == resumedomain.SectionTechStack || !hasContent(id, resume) {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		order = append(order, id)
		hash, err := sectionHash(id, resume)
		if err != nil {
			return domain.LayoutSnapshot{}, err
		}
		hashes[id] = hash
	}
	if len(order) == 0 {
		return domain.LayoutSnapshot{}, fmt.Errorf("public resume has no visible analytics sections")
	}
	return domain.NewLayoutSnapshot(order, hashes, publicTechnologies(resume))
}

func publicTechnologies(resume resumequeries.PublicResumeView) []string {
	values := make([]string, 0, len(resume.Technologies))
	for _, technology := range resume.Technologies {
		values = append(values, technology.Name)
	}
	if resume.Skills != nil {
		values = append(values, resume.Skills.HardSkills...)
	}
	for _, experience := range resume.Experiences {
		for _, technology := range experience.Stack {
			values = append(values, technology.Name)
		}
	}
	for _, project := range resume.Projects {
		for _, technology := range project.Stack {
			values = append(values, technology.Name)
		}
	}
	return values
}

func defaultSections() []resumequeries.PublicSectionSettings {
	return []resumequeries.PublicSectionSettings{
		{ID: "about", Enabled: true}, {ID: "techStack", Enabled: true}, {ID: "experience", Enabled: true},
		{ID: "skills", Enabled: true}, {ID: "education", Enabled: true}, {ID: "courses", Enabled: true},
		{ID: "projects", Enabled: true}, {ID: "activities", Enabled: true}, {ID: "contacts", Enabled: true},
	}
}

func hasContent(id resumedomain.SectionID, r resumequeries.PublicResumeView) bool {
	switch id {
	case resumedomain.SectionAbout:
		return true
	case resumedomain.SectionExperience:
		return len(r.Experiences) > 0
	case resumedomain.SectionSkills:
		return r.Skills != nil && (len(r.Skills.HardSkills) > 0 || len(r.Skills.SoftSkills) > 0)
	case resumedomain.SectionEducation:
		return len(r.Educations) > 0
	case resumedomain.SectionCourses:
		return len(r.Courses) > 0
	case resumedomain.SectionProjects:
		return len(r.Projects) > 0
	case resumedomain.SectionActivities:
		return len(r.Activities) > 0
	case resumedomain.SectionContacts:
		return r.Contacts != nil && (r.Contacts.Phone != nil || r.Contacts.Telegram != nil || r.Contacts.Email != nil || len(r.Contacts.HiddenContacts) > 0)
	default:
		return false
	}
}

func sectionHash(id resumedomain.SectionID, r resumequeries.PublicResumeView) (string, error) {
	var value any
	switch id {
	case resumedomain.SectionAbout:
		value = struct {
			Profile   any                                `json:"profile"`
			Languages []resumequeries.PublicLanguageView `json:"languages"`
		}{canonicalProfileProjection(r.Profile), r.Languages}
	case resumedomain.SectionExperience:
		value = canonicalExperienceProjection(r.Experiences)
	case resumedomain.SectionSkills:
		value = canonicalSkillsProjection(r.Skills)
	case resumedomain.SectionEducation:
		value = canonicalEducationProjection(r.Educations)
	case resumedomain.SectionCourses:
		value = canonicalCourseProjection(r.Courses)
	case resumedomain.SectionProjects:
		value = canonicalProjectProjection(r.Projects)
	case resumedomain.SectionActivities:
		value = canonicalActivityProjection(r.Activities)
	case resumedomain.SectionContacts:
		value = canonicalContactsProjection(r.Contacts)
	default:
		return "", fmt.Errorf("unsupported layout section %q", id)
	}
	canonical, err := json.Marshal(struct {
		Language string `json:"language"`
		Content  any    `json:"content"`
	}{r.Language, value})
	if err != nil {
		return "", fmt.Errorf("marshal section %q: %w", id, err)
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:8]), nil
}

func canonicalProfileProjection(v *resumequeries.PublicProfileView) any {
	if v == nil {
		return nil
	}
	return struct {
		Name        string   `json:"name"`
		Surname     string   `json:"surname"`
		Patronymic  *string  `json:"patronymic,omitempty"`
		Position    string   `json:"position"`
		PhotoURL    string   `json:"photoPath"`
		DateOfBirth string   `json:"dateOfBirth"`
		Location    string   `json:"location"`
		Relocation  bool     `json:"relocation"`
		WorkFormat  []string `json:"workFormat"`
	}{v.Name, v.Surname, v.Patronymic, v.Position, v.PhotoURL, v.DateOfBirth, v.Location, v.Relocation, v.WorkFormat}
}

func canonicalExperienceProjection(v []resumequeries.PublicExperienceView) []any {
	out := make([]any, len(v))
	for i, x := range v {
		out[i] = struct {
			Company      string                                        `json:"company"`
			Position     string                                        `json:"position"`
			Location     *string                                       `json:"location,omitempty"`
			StartDate    string                                        `json:"startDate"`
			EndDate      *string                                       `json:"endDate,omitempty"`
			Description  *string                                       `json:"description,omitempty"`
			Achievements []string                                      `json:"achievements"`
			Tasks        []string                                      `json:"tasks"`
			Stack        []resumequeries.PublicExperienceStackItemView `json:"stack"`
		}{x.Company, x.Position, x.Location, x.StartDate, x.EndDate, x.Description, x.Achievements, x.Tasks, x.Stack}
	}
	return out
}
func canonicalSkillsProjection(v *resumequeries.PublicSkillsView) any {
	if v == nil {
		return nil
	}
	return struct {
		HardSkills []string `json:"hardSkills"`
		SoftSkills []string `json:"softSkills"`
	}{v.HardSkills, v.SoftSkills}
}
func canonicalEducationProjection(v []resumequeries.PublicEducationView) []any {
	out := make([]any, len(v))
	for i, x := range v {
		out[i] = struct {
			Location    string  `json:"location"`
			University  string  `json:"university"`
			Grade       int     `json:"grade"`
			ProgramName string  `json:"programName"`
			ProfileName *string `json:"profileName,omitempty"`
			Faculty     string  `json:"faculty"`
			Course      *int    `json:"course,omitempty"`
			EndDate     *string `json:"endDate,omitempty"`
		}{x.Location, x.University, x.Grade, x.ProgramName, x.ProfileName, x.Faculty, x.Course, x.EndDate}
	}
	return out
}
func canonicalCourseProjection(v []resumequeries.PublicCourseView) []any {
	out := make([]any, len(v))
	for i, x := range v {
		out[i] = struct {
			Name            string  `json:"name"`
			Organization    string  `json:"organization"`
			EndDate         string  `json:"endDate"`
			Hours           *int    `json:"hours,omitempty"`
			CertificateLink *string `json:"certificateLink,omitempty"`
		}{x.Name, x.Organization, x.EndDate, x.Hours, x.CertificateLink}
	}
	return out
}
func canonicalActivityProjection(v []resumequeries.PublicActivityView) []any {
	out := make([]any, len(v))
	for i, x := range v {
		out[i] = struct {
			Type        string  `json:"type"`
			Title       string  `json:"title"`
			Description *string `json:"description,omitempty"`
			Date        *string `json:"date,omitempty"`
			Link        *string `json:"link,omitempty"`
		}{x.Type, x.Title, x.Description, x.Date, x.Link}
	}
	return out
}
func canonicalProjectProjection(v []resumequeries.PublicProjectView) []any {
	out := make([]any, len(v))
	for i, x := range v {
		out[i] = struct {
			Name         string                                       `json:"name"`
			Description  *string                                      `json:"description,omitempty"`
			Screenshots  []string                                     `json:"screenshots"`
			Link         *string                                      `json:"link,omitempty"`
			Stack        []resumequeries.PublicProjectTechnologyView  `json:"stack"`
			Contributors []resumequeries.PublicProjectContributorView `json:"contributors"`
		}{x.Name, x.Description, x.Screenshots, x.Link, x.Stack, x.Contributors}
	}
	return out
}
func canonicalContactsProjection(v *resumequeries.PublicContactsView) any {
	if v == nil {
		return nil
	}
	return struct {
		Phone          *string                               `json:"phone,omitempty"`
		Telegram       *string                               `json:"telegram,omitempty"`
		Email          *string                               `json:"email,omitempty"`
		HiddenContacts []string                              `json:"hiddenContacts"`
		ContactsText   *resumequeries.PublicContactsTextView `json:"contactsText,omitempty"`
	}{v.Phone, v.Telegram, v.Email, v.HiddenContacts, v.ContactsText}
}
