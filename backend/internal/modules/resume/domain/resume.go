package domain

import (
	"fmt"
	"strings"

	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Resume struct {
	id           valueobjects.ID
	userID       useridentity.ID
	title        string
	language     enums.ResumeLanguage
	version      int64
	overview     *ResumeOverview
	profile      *Profile
	contacts     *Contacts
	skills       *Skills
	languages    *ResumeLanguages
	technologies *ResumeTechnologies
	experiences  []*Experience
	educations   []*Education
	courses      []*Course
	activities   []*Activity
	projects     []*Project
}

type ResumeParams struct {
	ID       valueobjects.ID
	UserID   useridentity.ID
	Title    string
	Language enums.ResumeLanguage
	Version  int64
}

type ResumeSnapshot struct {
	ID           valueobjects.ID
	UserID       useridentity.ID
	Title        string
	Language     enums.ResumeLanguage
	Version      int64
	Overview     *ResumeOverview
	Profile      *Profile
	Contacts     *Contacts
	Skills       *Skills
	Languages    *ResumeLanguages
	Technologies *ResumeTechnologies
	Experiences  []*Experience
	Educations   []*Education
	Courses      []*Course
	Activities   []*Activity
	Projects     []*Project
}

func NewResume(params ResumeParams) (*Resume, error) {
	title := strings.TrimSpace(params.Title)
	if params.ID.IsZero() {
		params.ID = valueobjects.NewID()
	}
	if params.UserID.IsZero() || title == "" {
		return nil, fmt.Errorf("%w: invalid resume identity or title", ErrInvalidResume)
	}
	if !params.Language.IsValid() {
		return nil, fmt.Errorf("%w: invalid resume language", ErrInvalidResume)
	}
	if err := ValidateResumeTitle(title); err != nil {
		return nil, err
	}
	version := params.Version
	if version == 0 {
		version = 1
	}
	if version < 1 {
		return nil, fmt.Errorf("%w: resume version must be positive", ErrInvalidResume)
	}
	return &Resume{id: params.ID, userID: params.UserID, title: title, language: params.Language, version: version}, nil
}

func ReconstructResume(snapshot ResumeSnapshot) (*Resume, error) {
	resume, err := NewResume(ResumeParams{
		ID: snapshot.ID, UserID: snapshot.UserID, Title: snapshot.Title,
		Language: snapshot.Language, Version: snapshot.Version,
	})
	if err != nil {
		return nil, err
	}
	if err := resume.validateOwnedState(snapshot); err != nil {
		return nil, err
	}
	resume.overview = snapshot.Overview
	resume.profile = snapshot.Profile
	resume.contacts = snapshot.Contacts
	resume.skills = snapshot.Skills
	resume.languages = snapshot.Languages
	resume.technologies = snapshot.Technologies
	resume.experiences = append([]*Experience(nil), snapshot.Experiences...)
	resume.educations = append([]*Education(nil), snapshot.Educations...)
	resume.courses = append([]*Course(nil), snapshot.Courses...)
	resume.activities = append([]*Activity(nil), snapshot.Activities...)
	resume.projects = append([]*Project(nil), snapshot.Projects...)
	return resume, nil
}

func (r *Resume) validateOwnedState(snapshot ResumeSnapshot) error {
	check := func(id valueobjects.ID) error {
		if id.IsZero() || id != r.id {
			return fmt.Errorf("%w: child does not belong to resume", ErrResumeInvariant)
		}
		return nil
	}
	for _, item := range snapshot.Experiences {
		if item == nil {
			return fmt.Errorf("%w: nil experience", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Educations {
		if item == nil {
			return fmt.Errorf("%w: nil education", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Courses {
		if item == nil {
			return fmt.Errorf("%w: nil course", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Activities {
		if item == nil {
			return fmt.Errorf("%w: nil activity", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	for _, item := range snapshot.Projects {
		if item == nil {
			return fmt.Errorf("%w: nil project", ErrResumeInvariant)
		}
		if err := check(item.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Profile != nil {
		if err := check(snapshot.Profile.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Contacts != nil {
		if err := check(snapshot.Contacts.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Skills != nil {
		if err := check(snapshot.Skills.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Languages != nil {
		if err := check(snapshot.Languages.ResumeID()); err != nil {
			return err
		}
	}
	if snapshot.Technologies != nil {
		if err := check(snapshot.Technologies.ResumeID()); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resume) ID() valueobjects.ID            { return r.id }
func (r *Resume) UserID() useridentity.ID        { return r.userID }
func (r *Resume) Title() string                  { return r.title }
func (r *Resume) Language() enums.ResumeLanguage { return r.language }
func (r *Resume) Version() int64                 { return r.version }

func (r *Resume) ReplaceTitle(title string) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	title = strings.TrimSpace(title)
	if err := ValidateResumeTitle(title); err != nil {
		return err
	}
	r.title = title
	return nil
}

func (r *Resume) Snapshot() ResumeSnapshot {
	return ResumeSnapshot{
		ID: r.id, UserID: r.userID, Title: r.title, Language: r.language, Version: r.version,
		Overview: r.overview, Profile: r.profile, Contacts: r.contacts, Skills: r.skills,
		Languages: r.languages, Technologies: r.technologies,
		Experiences: append([]*Experience(nil), r.experiences...), Educations: append([]*Education(nil), r.educations...),
		Courses: append([]*Course(nil), r.courses...), Activities: append([]*Activity(nil), r.activities...), Projects: append([]*Project(nil), r.projects...),
	}
}

func (r *Resume) ReplaceExperiences(params []ExperienceParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[valueobjects.ID]*Experience, len(r.experiences))
	for _, item := range r.experiences {
		if item == nil {
			return fmt.Errorf("%w: nil experience", ErrResumeInvariant)
		}
		existing[item.ID()] = item
	}
	seen := make(map[valueobjects.ID]struct{}, len(params))
	replaced := make([]*Experience, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = valueobjects.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: experience does not belong to resume", ErrResumeInvariant)
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate experience id", ErrInvalidEntity)
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		experience, err := NewExperience(item)
		if err != nil {
			return err
		}
		replaced = append(replaced, experience)
	}
	r.experiences = replaced
	return nil
}

func (r *Resume) ReplaceEducations(params []EducationParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[valueobjects.ID]struct{}, len(r.educations))
	for _, item := range r.educations {
		if item == nil {
			return fmt.Errorf("%w: nil education", ErrResumeInvariant)
		}
		existing[item.ID()] = struct{}{}
	}
	seen := make(map[valueobjects.ID]struct{}, len(params))
	replaced := make([]*Education, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = valueobjects.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: education does not belong to resume", ErrResumeInvariant)
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate education id", ErrInvalidEntity)
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		education, err := NewEducation(item)
		if err != nil {
			return err
		}
		replaced = append(replaced, education)
	}
	r.educations = replaced
	return nil
}

func (r *Resume) ReplaceCourses(params []CourseParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[valueobjects.ID]struct{}, len(r.courses))
	for _, item := range r.courses {
		if item == nil {
			return fmt.Errorf("%w: nil course", ErrResumeInvariant)
		}
		existing[item.ID()] = struct{}{}
	}
	seen := make(map[valueobjects.ID]struct{}, len(params))
	replaced := make([]*Course, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = valueobjects.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: course does not belong to resume", ErrResumeInvariant)
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate course id", ErrInvalidEntity)
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		course, err := NewCourse(item)
		if err != nil {
			return err
		}
		replaced = append(replaced, course)
	}
	r.courses = replaced
	return nil
}

func (r *Resume) ReplaceProjects(params []ProjectParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[valueobjects.ID]struct{}, len(r.projects))
	for _, item := range r.projects {
		if item == nil {
			return fmt.Errorf("%w: nil project", ErrResumeInvariant)
		}
		existing[item.ID()] = struct{}{}
	}
	seen := make(map[valueobjects.ID]struct{}, len(params))
	replaced := make([]*Project, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = valueobjects.NewID()
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

func (r *Resume) ReplaceActivities(params []ActivityParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	existing := make(map[valueobjects.ID]*Activity, len(r.activities))
	for _, item := range r.activities {
		if item == nil {
			return fmt.Errorf("%w: nil activity", ErrResumeInvariant)
		}
		existing[item.ID()] = item
	}
	seen := make(map[valueobjects.ID]struct{}, len(params))
	replaced := make([]*Activity, 0, len(params))
	for _, item := range params {
		if item.ID.IsZero() {
			item.ID = valueobjects.NewID()
		} else if _, ok := existing[item.ID]; !ok {
			return fmt.Errorf("%w: activity %s does not belong to resume", ErrResumeInvariant, item.ID.String())
		}
		if _, ok := seen[item.ID]; ok {
			return fmt.Errorf("%w: duplicate activity %s", ErrInvalidEntity, item.ID.String())
		}
		seen[item.ID] = struct{}{}
		item.ResumeID = r.id
		activity, err := NewActivity(item)
		if err != nil {
			return err
		}
		replaced = append(replaced, activity)
	}
	r.activities = replaced
	return nil
}

func (r *Resume) Activities() []*Activity {
	if r == nil {
		return nil
	}
	return append([]*Activity(nil), r.activities...)
}

func (r *Resume) Courses() []*Course {
	if r == nil {
		return nil
	}
	return append([]*Course(nil), r.courses...)
}

func (r *Resume) Educations() []*Education {
	if r == nil {
		return nil
	}
	return append([]*Education(nil), r.educations...)
}

func (r *Resume) ReplaceOverview(params ResumeOverviewParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	o, err := NewResumeOverview(params)
	if err != nil {
		return err
	}
	r.overview = o
	return nil
}

func (r *Resume) ReplaceContacts(params ContactsParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.contacts != nil {
			params.ID = r.contacts.ID()
		} else {
			params.ID = valueobjects.NewID()
		}
	} else if r.contacts == nil || params.ID != r.contacts.ID() {
		return fmt.Errorf("%w: contacts do not belong to resume", ErrResumeInvariant)
	}
	params.ResumeID = r.id
	contacts, err := NewContacts(params)
	if err != nil {
		return err
	}
	r.contacts = contacts
	return nil
}

func (r *Resume) ReplaceSkills(params SkillsParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.skills != nil {
			params.ID = r.skills.ID()
		} else {
			params.ID = valueobjects.NewID()
		}
	} else if r.skills == nil || params.ID != r.skills.ID() {
		return fmt.Errorf("%w: skills do not belong to resume", ErrResumeInvariant)
	}
	params.ResumeID = r.id
	skills, err := NewSkills(params)
	if err != nil {
		return err
	}
	r.skills = skills
	return nil
}

func (r *Resume) Skills() *Skills {
	if r == nil || r.skills == nil {
		return nil
	}
	copy := *r.skills
	return &copy
}

func (r *Resume) ReplaceLanguages(params ResumeLanguagesParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.languages != nil {
			params.ID = r.languages.ID()
		} else {
			params.ID = valueobjects.NewID()
		}
	} else if r.languages == nil || params.ID != r.languages.ID() {
		return fmt.Errorf("%w: languages do not belong to resume", ErrResumeInvariant)
	}

	existing := make(map[valueobjects.ID]struct{})
	if r.languages != nil {
		for _, entry := range r.languages.Entries() {
			existing[entry.ID()] = struct{}{}
		}
	}
	seen := make(map[valueobjects.ID]struct{}, len(params.Entries))
	entries := make([]ResumeLanguageEntryParams, 0, len(params.Entries))
	for _, entry := range params.Entries {
		if entry.ID.IsZero() {
			entry.ID = valueobjects.NewID()
		} else {
			if _, ok := existing[entry.ID]; !ok {
				return fmt.Errorf("%w: language does not belong to resume", ErrResumeInvariant)
			}
		}
		if _, ok := seen[entry.ID]; ok {
			return fmt.Errorf("%w: duplicate language id", ErrInvalidEntity)
		}
		seen[entry.ID] = struct{}{}
		entries = append(entries, entry)
	}

	params.ResumeID = r.id
	params.Entries = entries
	languages, err := NewResumeLanguages(params)
	if err != nil {
		return err
	}
	r.languages = languages
	return nil
}

func (r *Resume) Languages() *ResumeLanguages {
	if r == nil || r.languages == nil {
		return nil
	}
	copy := *r.languages
	copy.entries = append([]ResumeLanguageEntry(nil), r.languages.entries...)
	return &copy
}

func (r *Resume) ReplaceTechnologies(params ResumeTechnologiesParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.technologies != nil {
			params.ID = r.technologies.ID()
		} else {
			params.ID = valueobjects.NewID()
		}
	} else if r.technologies == nil || params.ID != r.technologies.ID() {
		return fmt.Errorf("%w: technologies do not belong to resume", ErrResumeInvariant)
	}

	existing := make(map[valueobjects.ID]struct{})
	if r.technologies != nil {
		for _, entry := range r.technologies.Entries() {
			existing[entry.ID()] = struct{}{}
		}
	}
	seen := make(map[valueobjects.ID]struct{}, len(params.Entries))
	entries := make([]ResumeTechnologyEntryParams, 0, len(params.Entries))
	for _, entry := range params.Entries {
		if entry.ID.IsZero() {
			entry.ID = valueobjects.NewID()
		} else {
			if _, ok := existing[entry.ID]; !ok {
				return fmt.Errorf("%w: technology does not belong to resume", ErrResumeInvariant)
			}
		}
		if _, ok := seen[entry.ID]; ok {
			return fmt.Errorf("%w: duplicate technology id", ErrInvalidEntity)
		}
		seen[entry.ID] = struct{}{}
		entries = append(entries, entry)
	}

	params.ResumeID = r.id
	params.Entries = entries
	technologies, err := NewResumeTechnologies(params)
	if err != nil {
		return err
	}
	r.technologies = technologies
	return nil
}

func (r *Resume) Technologies() *ResumeTechnologies {
	if r == nil || r.technologies == nil {
		return nil
	}
	copy := *r.technologies
	copy.entries = append([]ResumeTechnologyEntry(nil), r.technologies.entries...)
	return &copy
}

func (r *Resume) Experiences() []*Experience {
	return append([]*Experience(nil), r.experiences...)
}

func (r *Resume) ReplaceProfile(params ProfileParams) error {
	if r == nil {
		return fmt.Errorf("%w: resume is nil", ErrResumeInvariant)
	}
	if params.ID.IsZero() {
		if r.profile != nil {
			params.ID = r.profile.ID()
		} else {
			params.ID = valueobjects.NewID()
		}
	} else if r.profile == nil || params.ID != r.profile.ID() {
		return fmt.Errorf("%w: profile does not belong to resume", ErrResumeInvariant)
	}
	params.ResumeID = r.id
	profile, err := NewProfile(params)
	if err != nil {
		return err
	}
	r.profile = profile
	return nil
}

func (r *Resume) Contacts() *Contacts {
	if r == nil || r.contacts == nil {
		return nil
	}
	copy := *r.contacts
	return &copy
}

func (r *Resume) Profile() *Profile {
	if r == nil || r.profile == nil {
		return nil
	}
	copy := *r.profile
	return &copy
}

func (r *Resume) Overview() *ResumeOverview {
	if r == nil || r.overview == nil {
		return nil
	}
	copy := *r.overview
	return &copy
}

func ValidateResumeTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w: resume title is empty", ErrInvalidResume)
	}
	return nil
}
