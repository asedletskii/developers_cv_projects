package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type ResumeReader struct {
	db *gorm.DB
}

func NewResumeReader(db *gorm.DB) *ResumeReader {
	return &ResumeReader{db: db}
}

func (r *ResumeReader) GetOwned(ctx context.Context, query queries.GetResumeQuery) (queries.ResumeView, error) {
	if r == nil {
		return queries.ResumeView{}, fmt.Errorf("get owned resume: reader is nil")
	}
	if query.ResumeID.IsZero() || query.UserID.IsZero() {
		return queries.ResumeView{}, fmt.Errorf("get owned resume: invalid identifiers")
	}
	db, err := resolveDB(ctx, r.db, "get owned resume")
	if err != nil {
		return queries.ResumeView{}, err
	}

	userID := query.UserID.UUID()
	metadata, err := r.loadMetadata(ctx, db, query.ResumeID.UUID(), &userID)
	if err != nil {
		return queries.ResumeView{}, err
	}
	return r.loadView(ctx, db, query.ResumeID.UUID(), metadata)
}

func (r *ResumeReader) GetPublic(ctx context.Context, query queries.PublicResumeQuery) (queries.ResumeView, error) {
	if r == nil {
		return queries.ResumeView{}, fmt.Errorf("get public resume: reader is nil")
	}
	if query.ID.IsZero() {
		return queries.ResumeView{}, fmt.Errorf("get public resume: invalid resume id")
	}
	db, err := resolveDB(ctx, r.db, "get public resume")
	if err != nil {
		return queries.ResumeView{}, err
	}

	metadata, err := r.loadMetadata(ctx, db, query.ID.UUID(), nil)
	if err != nil {
		return queries.ResumeView{}, err
	}
	return r.loadView(ctx, db, query.ID.UUID(), metadata)
}

func (r *ResumeReader) loadView(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, metadata queries.ResumeMetadataView) (queries.ResumeView, error) {
	view := queries.ResumeView{
		ID: metadata.ID, Title: metadata.Title, Language: metadata.Language,
		CreatedAt: metadata.CreatedAt, UpdatedAt: metadata.UpdatedAt, Version: int64Ptr(metadata.Version),
		Languages: make([]queries.LanguageView, 0), Technologies: make([]queries.TechnologyView, 0),
		Experiences: make([]queries.ExperienceView, 0), Educations: make([]queries.EducationView, 0),
		Courses: make([]queries.CourseView, 0), Activities: make([]queries.ActivityView, 0), Projects: make([]queries.ProjectView, 0),
	}

	if err := r.loadOverview(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadProfile(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadContacts(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadSkills(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadLanguages(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadTechnologies(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadExperiences(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadEducations(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadCourses(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadActivities(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	if err := r.loadProjects(ctx, db, resumeID, &view); err != nil {
		return queries.ResumeView{}, err
	}
	return view, nil
}

func (r *ResumeReader) loadMetadata(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, userID *uuid.UUID) (queries.ResumeMetadataView, error) {
	var model gormmodel.Resume
	readDB := db.Where("id = ?", resumeID)
	if userID != nil {
		readDB = readDB.Where("user_id = ?", *userID)
	}
	if err := readDB.WithContext(ctx).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ResumeMetadataView{}, queries.ErrResumeNotFound
		}
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: %w", err)
	}
	return queries.ResumeMetadataView{
		ID: model.ID.String(), Title: model.Title, Language: model.Language,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt, Version: model.Version,
	}, nil
}

func (r *ResumeReader) loadOverview(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeOverview
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ErrResumeNotFound
		}
		return fmt.Errorf("get overview: %w", err)
	}
	result := queries.ResumeOverviewView{
		Greeting:   queries.OverviewGreeting{First: model.GreetingFirst, Second: model.GreetingSecond, SecondHighlight: model.GreetingHighlight, Third: model.GreetingThird},
		ResumeLink: model.ResumeLink,
		Sections:   make([]queries.OverviewSectionSettings, 0),
	}
	if model.AvailabilityStatus != nil {
		result.Availability = &queries.OverviewAvailability{Status: *model.AvailabilityStatus, Note: model.AvailabilityNote}
		if model.AvailabilityFrom != nil {
			value := model.AvailabilityFrom.Format("2006-01-02")
			result.Availability.From = &value
		}
	}
	if model.ThemePrimary != nil || model.ThemeSecondary != nil {
		result.Theme = &queries.OverviewThemeSettings{Primary: valueOrEmpty(model.ThemePrimary), Secondary: valueOrEmpty(model.ThemeSecondary)}
	}
	var sections []gormmodel.ResumeOverviewSection
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("position").Find(&sections).Error; err != nil {
		return fmt.Errorf("get overview sections: %w", err)
	}
	for _, section := range sections {
		result.Sections = append(result.Sections, queries.OverviewSectionSettings{ID: section.SectionID, Enabled: section.Enabled, Collapsed: section.Collapsed})
	}
	view.Overview = result
	return nil
}

func (r *ResumeReader) loadProfile(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeProfile
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get resume profile: %w", err)
	}
	var formats []gormmodel.ResumeProfileWorkFormat
	if err := db.WithContext(ctx).Where("profile_id = ?", model.ID).Order("position").Find(&formats).Error; err != nil {
		return fmt.Errorf("get profile work formats: %w", err)
	}
	workFormat := make([]string, 0, len(formats))
	for _, format := range formats {
		workFormat = append(workFormat, format.Value)
	}
	view.Profile = &queries.ResumeProfileView{
		ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Surname: model.Surname,
		Patronymic: model.Patronymic, Position: model.PositionValue(), PhotoPath: model.PhotoPathValue(),
		DateOfBirth: model.DateOfBirth.Format("2006-01-02"), Location: model.LocationValue(), Relocation: model.Relocation,
		WorkFormat: workFormat,
	}
	return nil
}

func (r *ResumeReader) loadContacts(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeContacts
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get resume contacts: %w", err)
	}

	var hidden []gormmodel.ResumeContactHidden
	if err := db.WithContext(ctx).Where("contacts_id = ?", model.ID).Order("position").Find(&hidden).Error; err != nil {
		return fmt.Errorf("get hidden contacts: %w", err)
	}

	hiddenContacts := make([]string, 0, len(hidden))
	for _, item := range hidden {
		hiddenContacts = append(hiddenContacts, item.Value)
	}

	result := queries.ResumeContactsView{
		ID:             model.ID.String(),
		ResumeID:       model.ResumeID.String(),
		Phone:          model.Phone,
		Telegram:       model.Telegram,
		Email:          model.Email,
		HiddenContacts: hiddenContacts,
	}
	if err := r.loadContactsText(ctx, db, model.ID, &result); err != nil {
		return err
	}
	view.Contacts = &result
	return nil
}

func (r *ResumeReader) loadContactsText(ctx context.Context, db *gorm.DB, contactsID uuid.UUID, result *queries.ResumeContactsView) error {
	var text gormmodel.ResumeContactsText
	if err := db.WithContext(ctx).Where("contacts_id = ?", contactsID).First(&text).Error; err == nil {
		result.ContactsText = &queries.ResumeContactsTextView{Title: text.Title, Subtitle: text.Subtitle}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("get contacts text: %w", err)
	}
	return nil
}

func (r *ResumeReader) loadSkills(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var model gormmodel.ResumeSkills
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("get skills: %w", err)
	}
	var hard []gormmodel.ResumeSkillHard
	if err := db.WithContext(ctx).Where("skills_id = ?", model.ID).Order("position").Find(&hard).Error; err != nil {
		return fmt.Errorf("get hard skills: %w", err)
	}
	var soft []gormmodel.ResumeSkillSoft
	if err := db.WithContext(ctx).Where("skills_id = ?", model.ID).Order("position").Find(&soft).Error; err != nil {
		return fmt.Errorf("get soft skills: %w", err)
	}
	result := queries.SkillsView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), HardSkills: make([]string, 0, len(hard)), SoftSkills: make([]string, 0, len(soft))}
	for _, item := range hard {
		result.HardSkills = append(result.HardSkills, item.Value)
	}
	for _, item := range soft {
		result.SoftSkills = append(result.SoftSkills, item.Value)
	}
	view.Skills = &result
	return nil
}

func (r *ResumeReader) loadLanguages(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeLanguage
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("position").Find(&models).Error; err != nil {
		return fmt.Errorf("list languages: %w", err)
	}
	result := make([]queries.LanguageView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.LanguageView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Language: model.Language, Grade: model.Grade, Position: model.Position})
	}
	view.Languages = result
	return nil
}

func (r *ResumeReader) loadTechnologies(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeTechnology
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("position").Find(&models).Error; err != nil {
		return fmt.Errorf("list technologies: %w", err)
	}
	result := make([]queries.TechnologyView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.TechnologyView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Grade: model.Grade, Position: model.Position})
	}
	view.Technologies = result
	return nil
}

func (r *ResumeReader) loadExperiences(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var rows []experienceRow
	if err := db.WithContext(ctx).Table("resume_experiences AS e").Select("e.id::text AS id, e.resume_id::text AS resume_id, e.company, e.position, e.location, to_char(e.start_date, 'YYYY-MM') AS start_date, CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY-MM') END AS end_date, e.description").Where("e.resume_id = ?", resumeID).Order("e.sort_order").Scan(&rows).Error; err != nil {
		return fmt.Errorf("list experiences: %w", err)
	}
	result := make([]queries.ExperienceView, 0, len(rows))
	if len(rows) == 0 {
		view.Experiences = result
		return nil
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return fmt.Errorf("list experiences: invalid experience id: %w", err)
		}
		ids = append(ids, id)
	}
	achievements := make(map[uuid.UUID][]string, len(ids))
	var achievementRows []gormmodel.ResumeExperienceAchievement
	if err := db.WithContext(ctx).Where("experience_id IN ?", ids).Order("experience_id, position").Find(&achievementRows).Error; err != nil {
		return fmt.Errorf("get experience achievements: %w", err)
	}
	for _, item := range achievementRows {
		achievements[item.ExperienceID] = append(achievements[item.ExperienceID], item.Value)
	}
	stacks := make(map[uuid.UUID][]queries.ExperienceStackItemView, len(ids))
	var stackRows []gormmodel.ResumeExperienceStack
	if err := db.WithContext(ctx).Where("experience_id IN ?", ids).Order("experience_id, position").Find(&stackRows).Error; err != nil {
		return fmt.Errorf("get experience stack: %w", err)
	}
	for _, item := range stackRows {
		stacks[item.ExperienceID] = append(stacks[item.ExperienceID], queries.ExperienceStackItemView{Name: item.Name, Grade: item.Grade})
	}
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return err
		}
		result = append(result, queries.ExperienceView{ID: row.ID, ResumeID: row.ResumeID, Company: row.Company, Position: row.Position, Location: row.Location, StartDate: row.StartDate, EndDate: row.EndDate, Description: row.Description, Achievements: append([]string(nil), achievements[id]...), Stack: append([]queries.ExperienceStackItemView(nil), stacks[id]...)})
	}
	for i := range result {
		if result[i].Achievements == nil {
			result[i].Achievements = make([]string, 0)
		}
		if result[i].Stack == nil {
			result[i].Stack = make([]queries.ExperienceStackItemView, 0)
		}
	}
	view.Experiences = result
	return nil
}

type experienceRow struct {
	ID, ResumeID, Company, Position string
	Location                        *string
	StartDate                       string
	EndDate                         *string
	Description                     *string
}

func (r *ResumeReader) loadEducations(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeEducation
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("sort_order").Find(&models).Error; err != nil {
		return fmt.Errorf("list educations: %w", err)
	}
	result := make([]queries.EducationView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.EducationView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Location: model.Location, University: model.University, Grade: model.Grade, ProgramName: model.ProgramName, ProfileName: model.ProfileName, Faculty: model.Faculty, Course: model.Course, EndDate: optionalYear(model.EndDate)})
	}
	view.Educations = result
	return nil
}

func (r *ResumeReader) loadCourses(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeCourse
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("sort_order").Find(&models).Error; err != nil {
		return fmt.Errorf("list courses: %w", err)
	}
	result := make([]queries.CourseView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.CourseView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Organization: model.Organization, EndDate: yearMonth(model.EndDate), Hours: model.Hours, CertificateLink: model.CertificateLink})
	}
	view.Courses = result
	return nil
}

func (r *ResumeReader) loadActivities(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var models []gormmodel.ResumeActivity
	if err := db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("sort_order").Find(&models).Error; err != nil {
		return fmt.Errorf("list activities: %w", err)
	}
	result := make([]queries.ActivityView, 0, len(models))
	for _, model := range models {
		result = append(result, queries.ActivityView{ID: model.ID.String(), ResumeID: model.ResumeID.String(), Type: model.Type, Title: model.Title, Description: model.Description, Date: optionalYearMonth(model.Date), Link: model.Link})
	}
	view.Activities = result
	return nil
}

func (r *ResumeReader) loadProjects(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, view *queries.ResumeView) error {
	var rows []projectRow
	if err := db.WithContext(ctx).Table("resume_projects AS p").Select("p.id::text AS id, p.resume_id::text AS resume_id, p.name, p.description, p.role, p.link, p.git_link").Where("p.resume_id = ?", resumeID).Order("p.sort_order").Scan(&rows).Error; err != nil {
		return fmt.Errorf("list projects: %w", err)
	}
	result := make([]queries.ProjectView, 0, len(rows))
	if len(rows) == 0 {
		view.Projects = result
		return nil
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return fmt.Errorf("list projects: invalid project id: %w", err)
		}
		ids = append(ids, id)
	}
	screenshots := make(map[uuid.UUID][]string, len(ids))
	var screenshotRows []gormmodel.ResumeProjectScreenshot
	if err := db.WithContext(ctx).Where("project_id IN ?", ids).Order("project_id, position").Find(&screenshotRows).Error; err != nil {
		return fmt.Errorf("list project screenshots: %w", err)
	}
	for _, item := range screenshotRows {
		screenshots[item.ProjectID] = append(screenshots[item.ProjectID], item.Path)
	}
	stacks := make(map[uuid.UUID][]queries.ProjectTechnologyView, len(ids))
	var stackRows []gormmodel.ResumeProjectTechnology
	if err := db.WithContext(ctx).Where("project_id IN ?", ids).Order("project_id, position").Find(&stackRows).Error; err != nil {
		return fmt.Errorf("list project technologies: %w", err)
	}
	for _, item := range stackRows {
		stacks[item.ProjectID] = append(stacks[item.ProjectID], queries.ProjectTechnologyView{Name: item.Name, Grade: item.Grade})
	}
	contributors := make(map[uuid.UUID][]queries.ProjectContributorView, len(ids))
	var contributorRows []gormmodel.ResumeProjectContributor
	if err := db.WithContext(ctx).Where("project_id IN ?", ids).Order("project_id, position").Find(&contributorRows).Error; err != nil {
		return fmt.Errorf("list project contributors: %w", err)
	}
	for _, item := range contributorRows {
		contributors[item.ProjectID] = append(contributors[item.ProjectID], queries.ProjectContributorView{Name: item.Name, Surname: item.Surname, Role: item.Role, Contacts: queries.ProjectContactsView{Phone: item.Phone, Telegram: item.Telegram, Email: item.Email}})
	}
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return err
		}
		result = append(result, queries.ProjectView{ID: row.ID, ResumeID: row.ResumeID, Name: row.Name, Description: row.Description, Role: row.Role, Link: row.Link, GitLink: row.GitLink, Screenshots: append([]string(nil), screenshots[id]...), Stack: append([]queries.ProjectTechnologyView(nil), stacks[id]...), Contributors: append([]queries.ProjectContributorView(nil), contributors[id]...)})
	}
	for i := range result {
		if result[i].Screenshots == nil {
			result[i].Screenshots = make([]string, 0)
		}
		if result[i].Stack == nil {
			result[i].Stack = make([]queries.ProjectTechnologyView, 0)
		}
		if result[i].Contributors == nil {
			result[i].Contributors = make([]queries.ProjectContributorView, 0)
		}
	}
	view.Projects = result
	return nil
}

type projectRow struct {
	ID, ResumeID, Name               string
	Description, Role, Link, GitLink *string
}

func int64Ptr(value int64) *int64 { return &value }
func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (r *ResumeReader) GetContactForView(ctx context.Context, query queries.ContactViewQuery) (queries.ContactView, error) {
	if r == nil {
		return queries.ContactView{}, fmt.Errorf("get contact for view: reader is nil")
	}
	if query.ResumeID.IsZero() || !query.Contact.IsValid() {
		return queries.ContactView{}, fmt.Errorf("get contact for view: invalid query")
	}
	db, err := resolveDB(ctx, r.db, "get contact for view")
	if err != nil {
		return queries.ContactView{}, err
	}

	var resume gormmodel.Resume
	if err := db.WithContext(ctx).Where("id = ?", query.ResumeID.UUID()).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ContactView{}, queries.ErrResumeNotFound
		}
		return queries.ContactView{}, fmt.Errorf("get resume for contact view: %w", err)
	}

	var model gormmodel.ResumeContacts
	if err := db.WithContext(ctx).
		Where("resume_id = ?", query.ResumeID.UUID()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ContactView{}, queries.ErrContactNotFound
		}
		return queries.ContactView{}, fmt.Errorf("get contact for view: %w", err)
	}

	var hidden gormmodel.ResumeContactHidden
	if err := db.WithContext(ctx).
		Where("contacts_id = ? AND value = ?", model.ID, string(query.Contact)).
		First(&hidden).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.ContactView{}, queries.ErrContactNotHidden
		}
		return queries.ContactView{}, fmt.Errorf("get hidden contact for view: %w", err)
	}

	var value *string
	switch query.Contact {
	case enums.ContactKeyPhone:
		value = model.Phone
	case enums.ContactKeyTelegram:
		value = model.Telegram
	case enums.ContactKeyEmail:
		value = model.Email
	}
	if value == nil || *value == "" {
		return queries.ContactView{}, queries.ErrContactNotFound
	}

	return queries.ContactView{Contact: query.Contact, Value: *value}, nil
}
