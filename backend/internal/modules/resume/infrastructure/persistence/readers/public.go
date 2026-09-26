package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type PublicResumeReader struct{ db *gorm.DB }

func NewPublicResumeReader(db *gorm.DB) *PublicResumeReader { return &PublicResumeReader{db: db} }

func (r *PublicResumeReader) GetPublicResume(ctx context.Context, query queries.PublicResumeQuery) (queries.PublicResumeView, error) {
	db, err := resolveDB(ctx, r.db, "get public resume")
	if err != nil {
		return queries.PublicResumeView{}, err
	}
	if query.ID.IsZero() {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume: invalid resume id")
	}

	var resume gormmodel.Resume
	if err := db.Where("id = ?", query.ID.UUID()).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.PublicResumeView{}, queries.ErrResumeNotFound
		}
		return queries.PublicResumeView{}, fmt.Errorf("get public resume: %w", err)
	}

	var profile gormmodel.ResumeProfile
	if err := db.Where("resume_id = ?", resume.ID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.PublicResumeView{}, queries.ErrResumeNotFound
		}
		return queries.PublicResumeView{}, fmt.Errorf("get public resume profile: %w", err)
	}
	var overview gormmodel.ResumeOverview
	if err := db.Where("resume_id = ?", resume.ID).First(&overview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.PublicResumeView{}, queries.ErrResumeNotFound
		}
		return queries.PublicResumeView{}, fmt.Errorf("get public resume overview: %w", err)
	}

	view := queries.PublicResumeView{
		ID: resume.ID.String(), Name: profile.Name, Surname: profile.Surname, Patronymic: profile.Patronymic,
		Position: profile.PositionValue(), Title: resume.Title, Language: resume.Language, UpdatedAt: &resume.UpdatedAt,
		ResumeLink: overview.ResumeLink, PhotoPath: profile.PhotoPathValue(), DateOfBirth: profile.DateOfBirth.Format("2006-01-02"),
		Location: profile.LocationValue(), Relocation: profile.Relocation, Greeting: queries.Greeting{
			First: overview.GreetingFirst, Second: overview.GreetingSecond, SecondHighlight: overview.GreetingHighlight, Third: overview.GreetingThird,
		},
		WorkFormat: make([]string, 0), HiddenContacts: make([]string, 0), Languages: make([]queries.ResumeLanguage, 0),
		Experience: make([]queries.Experience, 0), Education: make([]queries.Education, 0), Courses: make([]queries.Course, 0),
		Activities: make([]queries.Activity, 0), HardSkills: make([]string, 0), SoftSkills: make([]string, 0), Technologies: make([]queries.Technology, 0),
		Projects: make([]queries.ProjectSummary, 0), Sections: make([]queries.SectionSettings, 0),
	}

	if overview.AvailabilityStatus != nil {
		availability := &queries.Availability{Status: *overview.AvailabilityStatus, Note: overview.AvailabilityNote}
		if overview.AvailabilityFrom != nil {
			value := overview.AvailabilityFrom.Format("2006-01-02")
			availability.From = &value
		}
		view.Availability = availability
	}
	if overview.ThemePrimary != nil || overview.ThemeSecondary != nil {
		view.Theme = &queries.ThemeSettings{Primary: valueOrEmpty(overview.ThemePrimary), Secondary: valueOrEmpty(overview.ThemeSecondary)}
	}

	var formats []gormmodel.ResumeProfileWorkFormat
	if err := db.Where("profile_id = ?", profile.ID).Order("position").Find(&formats).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume work formats: %w", err)
	}
	for _, item := range formats {
		view.WorkFormat = append(view.WorkFormat, item.Value)
	}

	var sections []gormmodel.ResumeOverviewSection
	if err := db.Where("resume_id = ?", resume.ID).Order("position").Find(&sections).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume sections: %w", err)
	}
	for _, item := range sections {
		view.Sections = append(view.Sections, queries.SectionSettings{ID: item.SectionID, Enabled: item.Enabled, Collapsed: item.Collapsed})
	}

	var contacts gormmodel.ResumeContacts
	if err := db.Where("resume_id = ?", resume.ID).First(&contacts).Error; err == nil {
		view.Contacts = queries.Contacts{Phone: contacts.Phone, Telegram: contacts.Telegram, Email: contacts.Email}
		var hidden []gormmodel.ResumeContactHidden
		if err := db.Where("contacts_id = ?", contacts.ID).Order("position").Find(&hidden).Error; err != nil {
			return queries.PublicResumeView{}, fmt.Errorf("get public resume hidden contacts: %w", err)
		}
		for _, item := range hidden {
			view.HiddenContacts = append(view.HiddenContacts, item.Value)
		}
		var contactsText gormmodel.ResumeContactsText
		if err := db.Where("contacts_id = ?", contacts.ID).First(&contactsText).Error; err == nil {
			view.ContactsText = &queries.ContactsText{Title: contactsText.Title, Subtitle: contactsText.Subtitle}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.PublicResumeView{}, fmt.Errorf("get public resume contacts text: %w", err)
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume contacts: %w", err)
	}

	var languages []gormmodel.ResumeLanguage
	if err := db.Where("resume_id = ?", resume.ID).Order("position").Find(&languages).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume languages: %w", err)
	}
	for _, item := range languages {
		view.Languages = append(view.Languages, queries.ResumeLanguage{Language: item.Language, Grade: item.Grade})
	}

	var skills gormmodel.ResumeSkills
	if err := db.Where("resume_id = ?", resume.ID).First(&skills).Error; err == nil {
		var hard []gormmodel.ResumeSkillHard
		if err := db.Where("skills_id = ?", skills.ID).Order("position").Find(&hard).Error; err != nil {
			return queries.PublicResumeView{}, fmt.Errorf("get public resume hard skills: %w", err)
		}
		for _, item := range hard {
			view.HardSkills = append(view.HardSkills, item.Value)
		}
		var soft []gormmodel.ResumeSkillSoft
		if err := db.Where("skills_id = ?", skills.ID).Order("position").Find(&soft).Error; err != nil {
			return queries.PublicResumeView{}, fmt.Errorf("get public resume soft skills: %w", err)
		}
		for _, item := range soft {
			view.SoftSkills = append(view.SoftSkills, item.Value)
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume skills: %w", err)
	}

	var technologies []gormmodel.ResumeTechnology
	if err := db.Where("resume_id = ?", resume.ID).Order("position").Find(&technologies).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume technologies: %w", err)
	}
	for _, item := range technologies {
		view.Technologies = append(view.Technologies, queries.Technology{Name: item.Name})
	}

	var experiences []gormmodel.ResumeExperience
	if err := db.Where("resume_id = ?", resume.ID).Order("sort_order").Find(&experiences).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume experiences: %w", err)
	}
	for _, item := range experiences {
		var achievements []gormmodel.ResumeExperienceAchievement
		if err := db.Where("experience_id = ?", item.ID).Order("position").Find(&achievements).Error; err != nil {
			return queries.PublicResumeView{}, fmt.Errorf("get public resume experience achievements: %w", err)
		}
		stack := make([]queries.Technology, 0)
		var experienceTech []gormmodel.ResumeExperienceTechnology
		if err := db.Where("experience_id = ?", item.ID).Order("position").Find(&experienceTech).Error; err != nil {
			return queries.PublicResumeView{}, fmt.Errorf("get public resume experience technologies: %w", err)
		}
		for _, tech := range experienceTech {
			stack = append(stack, queries.Technology{Name: tech.Name})
		}
		achievementValues := make([]string, 0, len(achievements))
		for _, achievement := range achievements {
			achievementValues = append(achievementValues, achievement.Value)
		}
		endDate := item.EndDate
		view.Experience = append(view.Experience, queries.Experience{ID: item.ID.String(), Company: item.Company, Position: item.Position, Location: item.Location, StartDate: yearMonth(item.StartDate), EndDate: optionalYearMonth(endDate), Description: item.Description, Achievements: achievementValues, Stack: stack})
	}

	var educations []gormmodel.ResumeEducation
	if err := db.Where("resume_id = ?", resume.ID).Order("sort_order").Find(&educations).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume educations: %w", err)
	}
	for _, item := range educations {
		view.Education = append(view.Education, queries.Education{ID: item.ID.String(), Location: item.Location, University: item.University, Grade: item.Grade, ProgramName: item.ProgramName, ProfileName: item.ProfileName, Faculty: item.Faculty, Course: item.Course, EndDate: optionalYear(item.EndDate)})
	}

	var courses []gormmodel.ResumeCourse
	if err := db.Where("resume_id = ?", resume.ID).Order("sort_order").Find(&courses).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume courses: %w", err)
	}
	for _, item := range courses {
		view.Courses = append(view.Courses, queries.Course{ID: item.ID.String(), Name: item.Name, Organization: item.Organization, EndDate: yearMonth(item.EndDate), Hours: item.Hours, CertificateLink: item.CertificateLink})
	}

	var activities []gormmodel.ResumeActivity
	if err := db.Where("resume_id = ?", resume.ID).Order("sort_order").Find(&activities).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume activities: %w", err)
	}
	for _, item := range activities {
		view.Activities = append(view.Activities, queries.Activity{ID: item.ID.String(), Type: item.Type, Title: item.Title, Description: item.Description, Date: optionalYearMonth(item.Date), Link: item.Link})
	}

	var projects []gormmodel.ResumeProject
	if err := db.Where("resume_id = ?", resume.ID).Order("sort_order").Find(&projects).Error; err != nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume projects: %w", err)
	}
	for _, item := range projects {
		var technologies []gormmodel.ResumeProjectTechnology
		if err := db.Where("project_id = ?", item.ID).Order("position").Find(&technologies).Error; err != nil {
			return queries.PublicResumeView{}, fmt.Errorf("get public resume project technologies: %w", err)
		}
		stack := make([]queries.Technology, 0, len(technologies))
		for _, tech := range technologies {
			stack = append(stack, queries.Technology{Name: tech.Name})
		}
		view.Projects = append(view.Projects, queries.ProjectSummary{ID: item.ID.String(), Name: item.Name, Description: item.Description, Stack: stack, Link: item.Link})
	}
	return view, nil
}
