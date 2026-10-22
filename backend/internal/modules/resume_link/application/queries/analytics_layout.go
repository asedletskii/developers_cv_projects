package queries

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
)

func BuildAnalyticsLayout(view resumecontracts.PublicResumeView) (analyticscontracts.ResumeAnalyticsLayout, error) {
	sections := view.Overview.Sections
	if len(sections) == 0 {
		sections = defaultAnalyticsSections()
	}
	ordered := make([]resumecontracts.PublicSectionSettings, 0, len(sections)+1)
	ordered = append(ordered, resumecontracts.PublicSectionSettings{ID: "about", Enabled: true})
	for _, section := range sections {
		if section.ID != "about" {
			ordered = append(ordered, section)
		}
	}

	layout := analyticscontracts.ResumeAnalyticsLayout{
		Sections:     make([]analyticscontracts.ResumeAnalyticsSection, 0, len(ordered)),
		Technologies: publicAnalyticsTechnologies(view),
	}
	seen := make(map[string]struct{}, len(ordered))
	for _, section := range ordered {
		if !section.Enabled || section.ID == "techStack" || !analyticsSectionHasContent(section.ID, view) {
			continue
		}
		if _, exists := seen[section.ID]; exists {
			continue
		}
		hash, err := analyticsSectionHash(section.ID, view)
		if err != nil {
			return analyticscontracts.ResumeAnalyticsLayout{}, err
		}
		seen[section.ID] = struct{}{}
		layout.Sections = append(layout.Sections, analyticscontracts.ResumeAnalyticsSection{ID: section.ID, Hash: hash})
	}
	if len(layout.Sections) == 0 {
		return analyticscontracts.ResumeAnalyticsLayout{}, fmt.Errorf("public resume has no visible analytics sections")
	}
	return layout, nil
}

func publicAnalyticsTechnologies(view resumecontracts.PublicResumeView) []string {
	values := make([]string, 0, len(view.Technologies))
	for _, technology := range view.Technologies {
		values = append(values, technology.Name)
	}
	if view.Skills != nil {
		values = append(values, view.Skills.HardSkills...)
	}
	for _, experience := range view.Experiences {
		for _, technology := range experience.Stack {
			values = append(values, technology.Name)
		}
	}
	for _, project := range view.Projects {
		for _, technology := range project.Stack {
			values = append(values, technology.Name)
		}
	}
	return values
}

func defaultAnalyticsSections() []resumecontracts.PublicSectionSettings {
	return []resumecontracts.PublicSectionSettings{
		{ID: "about", Enabled: true}, {ID: "techStack", Enabled: true}, {ID: "experience", Enabled: true},
		{ID: "skills", Enabled: true}, {ID: "education", Enabled: true}, {ID: "courses", Enabled: true},
		{ID: "projects", Enabled: true}, {ID: "activities", Enabled: true}, {ID: "contacts", Enabled: true},
	}
}

func analyticsSectionHasContent(id string, view resumecontracts.PublicResumeView) bool {
	switch id {
	case "about":
		return true
	case "experience":
		return len(view.Experiences) > 0
	case "skills":
		return view.Skills != nil && (len(view.Skills.HardSkills) > 0 || len(view.Skills.SoftSkills) > 0)
	case "education":
		return len(view.Educations) > 0
	case "courses":
		return len(view.Courses) > 0
	case "projects":
		return len(view.Projects) > 0
	case "activities":
		return len(view.Activities) > 0
	case "contacts":
		return view.Contacts != nil && (view.Contacts.Phone != nil || view.Contacts.Telegram != nil ||
			view.Contacts.Email != nil || len(view.Contacts.HiddenContacts) > 0)
	default:
		return false
	}
}

func analyticsSectionHash(id string, view resumecontracts.PublicResumeView) (string, error) {
	var value any
	switch id {
	case "about":
		value = struct {
			Profile   any                                  `json:"profile"`
			Languages []resumecontracts.PublicLanguageView `json:"languages"`
		}{canonicalProfileProjection(view.Profile), view.Languages}
	case "experience":
		value = canonicalExperienceProjection(view.Experiences)
	case "skills":
		value = canonicalSkillsProjection(view.Skills)
	case "education":
		value = canonicalEducationProjection(view.Educations)
	case "courses":
		value = canonicalCourseProjection(view.Courses)
	case "projects":
		value = canonicalProjectProjection(view.Projects)
	case "activities":
		value = canonicalActivityProjection(view.Activities)
	case "contacts":
		value = canonicalContactsProjection(view.Contacts)
	default:
		return "", fmt.Errorf("unsupported layout section %q", id)
	}
	canonical, err := json.Marshal(struct {
		Language string `json:"language"`
		Content  any    `json:"content"`
	}{view.Language, value})
	if err != nil {
		return "", fmt.Errorf("marshal section %q: %w", id, err)
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:8]), nil
}

func canonicalProfileProjection(view *resumecontracts.PublicProfileView) any {
	if view == nil {
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
	}{view.Name, view.Surname, view.Patronymic, view.Position, view.PhotoURL, view.DateOfBirth, view.Location, view.Relocation, view.WorkFormat}
}

func canonicalExperienceProjection(values []resumecontracts.PublicExperienceView) []any {
	out := make([]any, len(values))
	for i, value := range values {
		out[i] = struct {
			Company      string                                          `json:"company"`
			Position     string                                          `json:"position"`
			Location     *string                                         `json:"location,omitempty"`
			StartDate    string                                          `json:"startDate"`
			EndDate      *string                                         `json:"endDate,omitempty"`
			Description  *string                                         `json:"description,omitempty"`
			Achievements []string                                        `json:"achievements"`
			Tasks        []string                                        `json:"tasks"`
			Stack        []resumecontracts.PublicExperienceStackItemView `json:"stack"`
		}{value.Company, value.Position, value.Location, value.StartDate, value.EndDate, value.Description, value.Achievements, value.Tasks, value.Stack}
	}
	return out
}

func canonicalSkillsProjection(view *resumecontracts.PublicSkillsView) any {
	if view == nil {
		return nil
	}
	return struct {
		HardSkills []string `json:"hardSkills"`
		SoftSkills []string `json:"softSkills"`
	}{view.HardSkills, view.SoftSkills}
}

func canonicalEducationProjection(values []resumecontracts.PublicEducationView) []any {
	out := make([]any, len(values))
	for i, value := range values {
		out[i] = struct {
			Location    string  `json:"location"`
			University  string  `json:"university"`
			Grade       int     `json:"grade"`
			ProgramName string  `json:"programName"`
			ProfileName *string `json:"profileName,omitempty"`
			Faculty     string  `json:"faculty"`
			Course      *int    `json:"course,omitempty"`
			EndDate     *string `json:"endDate,omitempty"`
		}{value.Location, value.University, value.Grade, value.ProgramName, value.ProfileName, value.Faculty, value.Course, value.EndDate}
	}
	return out
}

func canonicalCourseProjection(values []resumecontracts.PublicCourseView) []any {
	out := make([]any, len(values))
	for i, value := range values {
		out[i] = struct {
			Name            string  `json:"name"`
			Organization    string  `json:"organization"`
			EndDate         string  `json:"endDate"`
			Hours           *int    `json:"hours,omitempty"`
			CertificateLink *string `json:"certificateLink,omitempty"`
		}{value.Name, value.Organization, value.EndDate, value.Hours, value.CertificateLink}
	}
	return out
}

func canonicalActivityProjection(values []resumecontracts.PublicActivityView) []any {
	out := make([]any, len(values))
	for i, value := range values {
		out[i] = struct {
			Type        string  `json:"type"`
			Title       string  `json:"title"`
			Description *string `json:"description,omitempty"`
			Date        *string `json:"date,omitempty"`
			Link        *string `json:"link,omitempty"`
		}{value.Type, value.Title, value.Description, value.Date, value.Link}
	}
	return out
}

func canonicalProjectProjection(values []resumecontracts.PublicProjectView) []any {
	out := make([]any, len(values))
	for i, value := range values {
		out[i] = struct {
			Name         string                                         `json:"name"`
			Description  *string                                        `json:"description,omitempty"`
			Screenshots  []string                                       `json:"screenshots"`
			Link         *string                                        `json:"link,omitempty"`
			Stack        []resumecontracts.PublicProjectTechnologyView  `json:"stack"`
			Contributors []resumecontracts.PublicProjectContributorView `json:"contributors"`
		}{value.Name, value.Description, value.Screenshots, value.Link, value.Stack, value.Contributors}
	}
	return out
}

func canonicalContactsProjection(view *resumecontracts.PublicContactsView) any {
	if view == nil {
		return nil
	}
	return struct {
		Phone          *string                                 `json:"phone,omitempty"`
		Telegram       *string                                 `json:"telegram,omitempty"`
		Email          *string                                 `json:"email,omitempty"`
		HiddenContacts []string                                `json:"hiddenContacts"`
		ContactsText   *resumecontracts.PublicContactsTextView `json:"contactsText,omitempty"`
	}{view.Phone, view.Telegram, view.Email, view.HiddenContacts, view.ContactsText}
}
