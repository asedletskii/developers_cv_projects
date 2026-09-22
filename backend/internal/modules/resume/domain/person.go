package domain

import (
	"fmt"
	"strings"
	"time"

	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type Person struct {
	id             valueobjects.ID
	name           string
	surname        string
	patronymic     *string
	position       string
	availability   *valueobjects.Availability
	resumeLink     *valueobjects.AssetPath
	theme          *valueobjects.ThemeSettings
	sections       []valueobjects.SectionSettings
	hiddenContacts []enums.ContactKey
	contactsText   *valueobjects.ContactsText
	greeting       valueobjects.Greeting
	photoPath      valueobjects.AssetPath
	dateOfBirth    valueobjects.Date
	location       string
	relocation     bool
	workFormats    []enums.WorkFormat
	contacts       entities.Contacts
	languages      []entities.PersonLanguage
	experience     []entities.Experience
	education      []entities.Education
	courses        []entities.Course
	activities     []entities.Activity
	hardSkills     []enums.HardSkill
	softSkills     []enums.SoftSkill
	technologies   []entities.Technology
	projects       []entities.Project
}

type PersonParams struct {
	ID             valueobjects.ID
	Name           string
	Surname        string
	Patronymic     *string
	Position       string
	Availability   *valueobjects.Availability
	ResumeLink     *valueobjects.AssetPath
	Theme          *valueobjects.ThemeSettings
	Sections       []valueobjects.SectionSettings
	HiddenContacts []enums.ContactKey
	ContactsText   *valueobjects.ContactsText
	Greeting       valueobjects.Greeting
	PhotoPath      valueobjects.AssetPath
	DateOfBirth    valueobjects.Date
	Location       string
	Relocation     bool
	WorkFormats    []enums.WorkFormat
	Contacts       entities.Contacts
	Languages      []entities.PersonLanguage
	Experience     []entities.Experience
	Education      []entities.Education
	Courses        []entities.Course
	Activities     []entities.Activity
	HardSkills     []enums.HardSkill
	SoftSkills     []enums.SoftSkill
	Technologies   []entities.Technology
	Projects       []entities.Project
}

type Skills struct {
	Hard []enums.HardSkill
	Soft []enums.SoftSkill
}

func NewPerson(params PersonParams) (*Person, error) {
	name := strings.TrimSpace(params.Name)
	surname := strings.TrimSpace(params.Surname)
	position := strings.TrimSpace(params.Position)
	location := strings.TrimSpace(params.Location)
	if params.ID.IsZero() || name == "" ||
		surname == "" || position == "" ||
		location == "" || params.PhotoPath.IsZero() {
		return nil, fmt.Errorf("%w: required person fields are empty", ErrInvalidPerson)
	}
	if params.DateOfBirth.IsZero() || params.DateOfBirth.Time().After(time.Now().UTC()) {
		return nil, fmt.Errorf("%w: date of birth is in the future", ErrInvalidPerson)
	}
	if err := validateSections(params.Sections); err != nil {
		return nil, err
	}
	if err := validateContactKeys(params.HiddenContacts); err != nil {
		return nil, err
	}
	if err := validateWorkFormats(params.WorkFormats); err != nil {
		return nil, err
	}
	if err := validateSkills(params.HardSkills, params.SoftSkills); err != nil {
		return nil, err
	}
	if err := validateUniqueEntityIDs(experienceIDs(params.Experience), "experience"); err != nil {
		return nil, err
	}
	if err := validateUniqueEntityIDs(educationIDs(params.Education), "education"); err != nil {
		return nil, err
	}
	if err := validateUniqueEntityIDs(courseIDs(params.Courses), "courses"); err != nil {
		return nil, err
	}
	if err := validateUniqueEntityIDs(activityIDs(params.Activities), "activities"); err != nil {
		return nil, err
	}
	if err := validateUniqueEntityIDs(projectIDs(params.Projects), "projects"); err != nil {
		return nil, err
	}

	return &Person{
		id:             params.ID,
		name:           name,
		surname:        surname,
		patronymic:     cloneString(params.Patronymic),
		position:       position,
		availability:   cloneAvailability(params.Availability),
		resumeLink:     cloneAssetPath(params.ResumeLink),
		theme:          cloneTheme(params.Theme),
		sections:       cloneSections(params.Sections),
		hiddenContacts: cloneContactKeys(params.HiddenContacts),
		contactsText:   cloneContactsText(params.ContactsText),
		greeting:       params.Greeting,
		photoPath:      params.PhotoPath,
		dateOfBirth:    params.DateOfBirth,
		location:       location,
		relocation:     params.Relocation,
		workFormats:    cloneWorkFormats(params.WorkFormats),
		contacts:       params.Contacts,
		languages:      cloneLanguages(params.Languages),
		experience:     cloneExperiences(params.Experience),
		education:      cloneEducation(params.Education),
		courses:        cloneCourses(params.Courses),
		activities:     cloneActivities(params.Activities),
		hardSkills:     cloneHardSkills(params.HardSkills),
		softSkills:     cloneSoftSkills(params.SoftSkills),
		technologies:   cloneTechnologies(params.Technologies),
		projects:       cloneProjects(params.Projects),
	}, nil
}

func (p *Person) ID() valueobjects.ID                      { return p.id }
func (p *Person) Name() string                             { return p.name }
func (p *Person) Surname() string                          { return p.surname }
func (p *Person) Patronymic() *string                      { return cloneString(p.patronymic) }
func (p *Person) Position() string                         { return p.position }
func (p *Person) Availability() *valueobjects.Availability { return cloneAvailability(p.availability) }
func (p *Person) ResumeLink() *valueobjects.AssetPath      { return cloneAssetPath(p.resumeLink) }
func (p *Person) Theme() *valueobjects.ThemeSettings       { return cloneTheme(p.theme) }
func (p *Person) Sections() []valueobjects.SectionSettings { return cloneSections(p.sections) }
func (p *Person) HiddenContacts() []enums.ContactKey       { return cloneContactKeys(p.hiddenContacts) }
func (p *Person) ContactsText() *valueobjects.ContactsText { return cloneContactsText(p.contactsText) }
func (p *Person) Greeting() valueobjects.Greeting          { return p.greeting }
func (p *Person) PhotoPath() valueobjects.AssetPath        { return p.photoPath }
func (p *Person) DateOfBirth() valueobjects.Date           { return p.dateOfBirth }
func (p *Person) Location() string                         { return p.location }
func (p *Person) Relocation() bool                         { return p.relocation }
func (p *Person) WorkFormats() []enums.WorkFormat          { return cloneWorkFormats(p.workFormats) }
func (p *Person) Contacts() entities.Contacts              { return p.contacts }
func (p *Person) Languages() []entities.PersonLanguage     { return cloneLanguages(p.languages) }
func (p *Person) Experience() []entities.Experience        { return cloneExperiences(p.experience) }
func (p *Person) Education() []entities.Education          { return cloneEducation(p.education) }
func (p *Person) Courses() []entities.Course               { return cloneCourses(p.courses) }
func (p *Person) Activities() []entities.Activity          { return cloneActivities(p.activities) }
func (p *Person) HardSkills() []enums.HardSkill            { return cloneHardSkills(p.hardSkills) }
func (p *Person) SoftSkills() []enums.SoftSkill            { return cloneSoftSkills(p.softSkills) }
func (p *Person) Technologies() []entities.Technology      { return cloneTechnologies(p.technologies) }
func (p *Person) Projects() []entities.Project             { return cloneProjects(p.projects) }

func (p *Person) ChangePosition(position string) error {
	position = strings.TrimSpace(position)
	if position == "" {
		return fmt.Errorf("%w: position is empty", ErrPersonInvariant)
	}
	p.position = position
	return nil
}

func (p *Person) ChangeLocation(location string) error {
	location = strings.TrimSpace(location)
	if location == "" {
		return fmt.Errorf("%w: location is empty", ErrPersonInvariant)
	}
	p.location = location
	return nil
}

func (p *Person) ChangePatronymic(patronymic *string) {
	p.patronymic = cloneString(patronymic)
}

func (p *Person) SetAvailability(availability *valueobjects.Availability) {
	p.availability = cloneAvailability(availability)
}

func (p *Person) SetResumeLink(resumeLink *valueobjects.AssetPath) {
	p.resumeLink = cloneAssetPath(resumeLink)
}

func (p *Person) SetTheme(theme *valueobjects.ThemeSettings) {
	p.theme = cloneTheme(theme)
}

func (p *Person) SetContactsText(contactsText *valueobjects.ContactsText) {
	p.contactsText = cloneContactsText(contactsText)
}

func (p *Person) SetContacts(contacts entities.Contacts) {
	p.contacts = contacts
}

func (p *Person) SetRelocation(relocation bool) {
	p.relocation = relocation
}

func (p *Person) ReplaceSections(sections []valueobjects.SectionSettings) error {
	if err := validateSections(sections); err != nil {
		return err
	}
	p.sections = cloneSections(sections)
	return nil
}

func (p *Person) SetHiddenContacts(keys []enums.ContactKey) error {
	if err := validateContactKeys(keys); err != nil {
		return err
	}
	p.hiddenContacts = cloneContactKeys(keys)
	return nil
}

func (p *Person) SetWorkFormats(formats []enums.WorkFormat) error {
	if err := validateWorkFormats(formats); err != nil {
		return err
	}
	p.workFormats = cloneWorkFormats(formats)
	return nil
}

func (p *Person) SetSkills(skills Skills) error {
	if err := validateSkills(skills.Hard, skills.Soft); err != nil {
		return err
	}
	p.hardSkills = cloneHardSkills(skills.Hard)
	p.softSkills = cloneSoftSkills(skills.Soft)
	return nil
}

func (p *Person) AddExperience(experience entities.Experience) error {
	for _, existing := range p.experience {
		if existing.ID().String() == experience.ID().String() {
			return fmt.Errorf("%w: experience", ErrDuplicateEntity)
		}
	}
	p.experience = append(p.experience, experience)
	return nil
}

func (p *Person) RemoveExperience(id valueobjects.ID) error {
	return removeEntity(id, &p.experience)
}

func (p *Person) AddEducation(education entities.Education) error {
	for _, existing := range p.education {
		if existing.ID().String() == education.ID().String() {
			return fmt.Errorf("%w: education", ErrDuplicateEntity)
		}
	}
	p.education = append(p.education, education)
	return nil
}

func (p *Person) RemoveEducation(id valueobjects.ID) error {
	return removeEntity(id, &p.education)
}

func (p *Person) AddCourse(course entities.Course) error {
	for _, existing := range p.courses {
		if existing.ID().String() == course.ID().String() {
			return fmt.Errorf("%w: course", ErrDuplicateEntity)
		}
	}
	p.courses = append(p.courses, course)
	return nil
}

func (p *Person) RemoveCourse(id valueobjects.ID) error {
	return removeEntity(id, &p.courses)
}

func (p *Person) AddActivity(activity entities.Activity) error {
	for _, existing := range p.activities {
		if existing.ID().String() == activity.ID().String() {
			return fmt.Errorf("%w: activity", ErrDuplicateEntity)
		}
	}
	p.activities = append(p.activities, activity)
	return nil
}

func (p *Person) RemoveActivity(id valueobjects.ID) error {
	return removeEntity(id, &p.activities)
}

func (p *Person) AddProject(project entities.Project) error {
	for _, existing := range p.projects {
		if existing.ID().String() == project.ID().String() {
			return fmt.Errorf("%w: project", ErrDuplicateEntity)
		}
	}
	p.projects = append(p.projects, project)
	return nil
}

func (p *Person) RemoveProject(id valueobjects.ID) error {
	return removeEntity(id, &p.projects)
}

func (p *Person) AddTechnology(technology entities.Technology) error {
	for _, existing := range p.technologies {
		if strings.EqualFold(existing.Name(), technology.Name()) {
			return fmt.Errorf("%w: technology %q", ErrDuplicateEntity, technology.Name())
		}
	}
	p.technologies = append(p.technologies, technology)
	return nil
}

func (p *Person) RemoveTechnology(name string) error {
	name = strings.TrimSpace(name)
	for index, technology := range p.technologies {
		if strings.EqualFold(technology.Name(), name) {
			p.technologies = append(p.technologies[:index], p.technologies[index+1:]...)
			return nil
		}
	}
	return ErrEntityNotFound
}

func (p *Person) ReplaceLanguages(languages []entities.PersonLanguage) {
	p.languages = cloneLanguages(languages)
}

func (p *Person) SetPhotoPath(photoPath valueobjects.AssetPath) {
	p.photoPath = photoPath
}

func (p *Person) SetDateOfBirth(dateOfBirth valueobjects.Date) error {
	if dateOfBirth.IsZero() || dateOfBirth.Time().After(time.Now().UTC()) {
		return fmt.Errorf("%w: date of birth is in the future", ErrPersonInvariant)
	}
	p.dateOfBirth = dateOfBirth
	return nil
}

func removeEntity[T interface{ ID() valueobjects.ID }](id valueobjects.ID, values *[]T) error {
	for index, entity := range *values {
		if entity.ID().String() == id.String() {
			*values = append((*values)[:index], (*values)[index+1:]...)
			return nil
		}
	}
	return ErrEntityNotFound
}

func validateSections(values []valueobjects.SectionSettings) error {
	seen := make(map[enums.SectionID]struct{}, len(values))
	for _, value := range values {
		if !value.ID().IsValid() {
			return fmt.Errorf("%w: invalid section", ErrPersonInvariant)
		}
		if _, ok := seen[value.ID()]; ok {
			return fmt.Errorf("%w: section %q", ErrDuplicateEntity, value.ID())
		}
		seen[value.ID()] = struct{}{}
	}
	return nil
}

func validateContactKeys(values []enums.ContactKey) error {
	seen := make(map[enums.ContactKey]struct{}, len(values))
	for _, value := range values {
		if !value.IsValid() {
			return fmt.Errorf("%w: invalid hidden contact", ErrPersonInvariant)
		}
		if _, ok := seen[value]; ok {
			return fmt.Errorf("%w: hidden contact %q", ErrDuplicateEntity, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validateWorkFormats(values []enums.WorkFormat) error {
	seen := make(map[enums.WorkFormat]struct{}, len(values))
	for _, value := range values {
		if !value.IsValid() {
			return fmt.Errorf("%w: invalid work format", ErrPersonInvariant)
		}
		if _, ok := seen[value]; ok {
			return fmt.Errorf("%w: work format %q", ErrDuplicateEntity, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validateSkills(hard []enums.HardSkill, soft []enums.SoftSkill) error {
	for _, value := range hard {
		if !value.IsValid() {
			return fmt.Errorf("%w: invalid hard skill", ErrPersonInvariant)
		}
	}
	for _, value := range soft {
		if !value.IsValid() {
			return fmt.Errorf("%w: invalid soft skill", ErrPersonInvariant)
		}
	}
	return nil
}

func validateUniqueEntityIDs(values []valueobjects.ID, entityName string) error {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value.IsZero() {
			return fmt.Errorf("%w: %s id is empty", ErrInvalidEntity, entityName)
		}
		if _, ok := seen[value.String()]; ok {
			return fmt.Errorf("%w: %s", ErrDuplicateEntity, entityName)
		}
		seen[value.String()] = struct{}{}
	}
	return nil
}

func experienceIDs(values []entities.Experience) []valueobjects.ID {
	result := make([]valueobjects.ID, 0, len(values))
	for _, value := range values {
		result = append(result, value.ID())
	}
	return result
}
func educationIDs(values []entities.Education) []valueobjects.ID {
	result := make([]valueobjects.ID, 0, len(values))
	for _, value := range values {
		result = append(result, value.ID())
	}
	return result
}
func courseIDs(values []entities.Course) []valueobjects.ID {
	result := make([]valueobjects.ID, 0, len(values))
	for _, value := range values {
		result = append(result, value.ID())
	}
	return result
}
func activityIDs(values []entities.Activity) []valueobjects.ID {
	result := make([]valueobjects.ID, 0, len(values))
	for _, value := range values {
		result = append(result, value.ID())
	}
	return result
}
func projectIDs(values []entities.Project) []valueobjects.ID {
	result := make([]valueobjects.ID, 0, len(values))
	for _, value := range values {
		result = append(result, value.ID())
	}
	return result
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
func cloneAvailability(value *valueobjects.Availability) *valueobjects.Availability {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
func cloneAssetPath(value *valueobjects.AssetPath) *valueobjects.AssetPath {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
func cloneTheme(value *valueobjects.ThemeSettings) *valueobjects.ThemeSettings {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
func cloneSections(values []valueobjects.SectionSettings) []valueobjects.SectionSettings {
	if values == nil {
		return nil
	}
	return append([]valueobjects.SectionSettings{}, values...)
}
func cloneContactKeys(values []enums.ContactKey) []enums.ContactKey {
	if values == nil {
		return nil
	}
	return append([]enums.ContactKey{}, values...)
}
func cloneContactsText(value *valueobjects.ContactsText) *valueobjects.ContactsText {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
func cloneWorkFormats(values []enums.WorkFormat) []enums.WorkFormat {
	if values == nil {
		return nil
	}
	return append([]enums.WorkFormat{}, values...)
}
func cloneLanguages(values []entities.PersonLanguage) []entities.PersonLanguage {
	if values == nil {
		return nil
	}
	return append([]entities.PersonLanguage{}, values...)
}
func cloneExperiences(values []entities.Experience) []entities.Experience {
	if values == nil {
		return nil
	}
	return append([]entities.Experience{}, values...)
}
func cloneEducation(values []entities.Education) []entities.Education {
	if values == nil {
		return nil
	}
	return append([]entities.Education{}, values...)
}
func cloneCourses(values []entities.Course) []entities.Course {
	if values == nil {
		return nil
	}
	return append([]entities.Course{}, values...)
}
func cloneActivities(values []entities.Activity) []entities.Activity {
	if values == nil {
		return nil
	}
	return append([]entities.Activity{}, values...)
}
func cloneHardSkills(values []enums.HardSkill) []enums.HardSkill {
	if values == nil {
		return nil
	}
	return append([]enums.HardSkill{}, values...)
}
func cloneSoftSkills(values []enums.SoftSkill) []enums.SoftSkill {
	if values == nil {
		return nil
	}
	return append([]enums.SoftSkill{}, values...)
}
func cloneTechnologies(values []entities.Technology) []entities.Technology {
	if values == nil {
		return nil
	}
	return append([]entities.Technology{}, values...)
}
func cloneProjects(values []entities.Project) []entities.Project {
	if values == nil {
		return nil
	}
	return append([]entities.Project{}, values...)
}
