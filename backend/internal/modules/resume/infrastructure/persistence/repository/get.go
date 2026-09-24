package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/database"
)

func (r *ResumeRepository) Get(ctx context.Context, query commands.GetResumeRepositoryQuery) (*domain.Resume, error) {
	userID := query.UserID
	id := query.ResumeID
	var resumeModel gormmodel.Resume
	db := database.DBFromContext(ctx, r.db)
	if userID.IsZero() {
		return nil, fmt.Errorf("get resume: invalid user id")
	}
	if id.IsZero() {
		return nil, fmt.Errorf("get resume: invalid resume id")
	}
	if err := db.WithContext(ctx).First(&resumeModel, "id = ? AND user_id = ?", id.UUID(), userID.UUID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume: %w", err)
	}
	return r.loadAggregate(ctx, db, resumeModel)
}
func (r *ResumeRepository) loadAggregate(
	ctx context.Context,
	db *gorm.DB, pm gormmodel.Resume,
) (*domain.Resume, error) {
	id, err := valueobjects.IDFromUUID(pm.ID)
	if err != nil {
		return nil, err
	}
	date, err := valueobjects.NewDate(pm.DateOfBirth)
	if err != nil {
		return nil, err
	}
	photo, err := valueobjects.NewAssetPath(pm.PhotoPath)
	if err != nil {
		return nil, err
	}
	var resume *valueobjects.AssetPath
	if pm.ResumeLink != nil {
		value, e := valueobjects.NewAssetPath(*pm.ResumeLink)
		if e != nil {
			return nil, e
		}
		resume = &value
	}
	var patronymic = cloneString(pm.Patronymic)
	var greeting gormmodel.Greeting
	if err := db.WithContext(ctx).First(&greeting, "resume_id = ?", pm.ID).Error; err != nil {
		return nil, err
	}
	gt, err := newGreeting(greeting)
	if err != nil {
		return nil, err
	}
	contacts, err := r.loadContacts(ctx, db, pm.ID)
	if err != nil {
		return nil, err
	}
	userID, err := useridentity.IDFromUUID(pm.UserID)
	if err != nil {
		return nil, fmt.Errorf("map resume user id: %w", err)
	}
	params := domain.ResumeParams{
		ID:          id,
		UserID:      userID,
		Name:        pm.Name,
		Surname:     pm.Surname,
		Patronymic:  patronymic,
		Position:    pm.Position,
		Title:       pm.Title,
		ResumeLink:  resume,
		Greeting:    gt,
		PhotoPath:   photo,
		DateOfBirth: date,
		Location:    pm.Location,
		Relocation:  pm.Relocation,
		Contacts:    contacts,
	}
	if err := r.loadOptionalProfile(ctx, db, pm.ID, &params); err != nil {
		return nil, err
	}
	if err := r.loadSettings(ctx, db, pm.ID, &params); err != nil {
		return nil, err
	}
	if err := r.loadContent(ctx, db, pm.ID, &params); err != nil {
		return nil, err
	}
	return domain.NewResume(params)
}

func (r *ResumeRepository) loadOptionalProfile(ctx context.Context, db *gorm.DB, id uuid.UUID, p *domain.ResumeParams) error {
	var m gormmodel.Availability
	if err := db.WithContext(ctx).First(&m, "resume_id = ?", id).Error; err == nil {
		status := enums.AvailabilityStatus(m.Status)
		var from *valueobjects.Date
		if m.From != nil {
			value, e := valueobjects.NewDate(*m.From)
			if e != nil {
				return e
			}
			from = &value
		}
		var note *valueobjects.LocalizedText
		if m.NoteRU != nil || m.NoteEN != nil {
			value, e := valueobjects.NewLocalizedText(
				valueobjects.LocalizedTextParams{
					RU: stringValue(m.NoteRU),
					EN: stringValue(m.NoteEN),
				})
			if e != nil {
				return e
			}
			note = &value
		}
		value, e := valueobjects.NewAvailability(
			valueobjects.AvailabilityParams{
				Status: status,
				From:   from,
				Note:   note,
			})
		if e != nil {
			return e
		}
		p.Availability = &value
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	var theme gormmodel.Theme
	if err := db.WithContext(ctx).First(&theme, "resume_id = ?", id).Error; err == nil {
		var primary, secondary *valueobjects.Color
		if theme.Primary != nil {
			value, e := valueobjects.NewColor(*theme.Primary)
			if e != nil {
				return e
			}
			primary = &value
		}
		if theme.Secondary != nil {
			value, e := valueobjects.NewColor(*theme.Secondary)
			if e != nil {
				return e
			}
			secondary = &value
		}
		value, e := valueobjects.NewThemeSettings(
			valueobjects.ThemeSettingsParams{
				Primary:   primary,
				Secondary: secondary,
			})
		if e != nil {
			return e
		}
		p.Theme = &value
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	var text gormmodel.ContactsText
	if err := db.WithContext(ctx).First(&text, "resume_id = ?", id).Error; err == nil {
		var title, subtitle *valueobjects.LocalizedText
		if text.TitleRU != nil || text.TitleEN != nil {
			value, e := valueobjects.NewLocalizedText(
				valueobjects.LocalizedTextParams{
					RU: stringValue(text.TitleRU),
					EN: stringValue(text.TitleEN),
				},
			)
			if e != nil {
				return e
			}
			title = &value
		}
		if text.SubtitleRU != nil || text.SubtitleEN != nil {
			value, e := valueobjects.NewLocalizedText(
				valueobjects.LocalizedTextParams{
					RU: stringValue(text.SubtitleRU),
					EN: stringValue(text.SubtitleEN),
				},
			)
			if e != nil {
				return e
			}
			subtitle = &value
		}
		value, e := valueobjects.NewContactsText(
			valueobjects.ContactsTextParams{
				Title:    title,
				Subtitle: subtitle,
			},
		)
		if e != nil {
			return e
		}
		p.ContactsText = &value
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (r *ResumeRepository) loadContacts(
	ctx context.Context,
	db *gorm.DB,
	id uuid.UUID,
) (entities.Contacts, error) {
	var m gormmodel.Contact
	if err := db.WithContext(ctx).First(&m, "resume_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.EmptyContacts(), nil
		}
		return entities.Contacts{}, err
	}
	var telegram *valueobjects.Telegram
	if m.Telegram != nil {
		value, e := valueobjects.NewTelegram(*m.Telegram)
		if e != nil {
			return entities.Contacts{}, e
		}
		telegram = &value
	}
	var email *valueobjects.Email
	if m.Email != nil {
		value, e := valueobjects.NewEmail(*m.Email)
		if e != nil {
			return entities.Contacts{}, e
		}
		email = &value
	}
	return entities.NewContacts(
		entities.ContactsParams{
			Phone:    stringValue(m.Phone),
			Telegram: telegram,
			Email:    email,
		},
	)
}

func (r *ResumeRepository) loadSettings(
	ctx context.Context,
	db *gorm.DB,
	id uuid.UUID,
	p *domain.ResumeParams,
) error {
	var sections []gormmodel.Section
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&sections).Error; err != nil {
		return err
	}
	p.Sections = make([]valueobjects.SectionSettings, 0, len(sections))
	for _, m := range sections {
		v, e := valueobjects.NewSectionSettings(
			valueobjects.SectionSettingsParams{
				ID:        enums.SectionID(m.SectionID),
				Enabled:   m.Enabled,
				Collapsed: m.Collapsed,
			},
		)
		if e != nil {
			return e
		}
		p.Sections = append(p.Sections, v)
	}
	var formats []gormmodel.ResumeWorkFormat
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&formats).Error; err != nil {
		return err
	}
	p.WorkFormats = make([]enums.WorkFormat, 0, len(formats))
	for _, m := range formats {
		p.WorkFormats = append(p.WorkFormats, enums.WorkFormat(m.Value))
	}
	var hidden []gormmodel.ResumeHiddenContact
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&hidden).Error; err != nil {
		return err
	}
	p.HiddenContacts = make([]enums.ContactKey, 0, len(hidden))
	for _, m := range hidden {
		p.HiddenContacts = append(p.HiddenContacts, enums.ContactKey(m.Value))
	}
	var hard []gormmodel.ResumeHardSkill
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&hard).Error; err != nil {
		return err
	}
	p.HardSkills = make([]enums.HardSkill, 0, len(hard))
	for _, m := range hard {
		p.HardSkills = append(p.HardSkills, enums.HardSkill(m.Value))
	}
	var soft []gormmodel.ResumeSoftSkill
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&soft).Error; err != nil {
		return err
	}
	p.SoftSkills = make([]enums.SoftSkill, 0, len(soft))
	for _, m := range soft {
		p.SoftSkills = append(p.SoftSkills, enums.SoftSkill(m.Value))
	}
	var langs []gormmodel.ResumeLanguage
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&langs).Error; err != nil {
		return err
	}
	p.Languages = make([]entities.ResumeLanguage, 0, len(langs))
	for _, m := range langs {
		l := enums.LanguageType(m.Language)
		var g *enums.LanguageGrade
		if m.Grade != nil {
			v := enums.LanguageGrade(*m.Grade)
			g = &v
		}
		v, e := entities.NewResumeLanguage(entities.ResumeLanguageParams{Language: l, Grade: g})
		if e != nil {
			return e
		}
		p.Languages = append(p.Languages, v)
	}
	var tech []gormmodel.ResumeTechnology
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&tech).Error; err != nil {
		return err
	}
	p.Technologies = make([]entities.Technology, 0, len(tech))
	for _, m := range tech {
		var g *enums.TechnologyGrade
		if m.Grade != nil {
			v := enums.TechnologyGrade(*m.Grade)
			g = &v
		}
		v, e := entities.NewTechnology(entities.TechnologyParams{Name: m.Name, Grade: g})
		if e != nil {
			return e
		}
		p.Technologies = append(p.Technologies, v)
	}
	return nil
}

func (r *ResumeRepository) loadContent(
	ctx context.Context,
	db *gorm.DB,
	id uuid.UUID,
	p *domain.ResumeParams,
) error {
	var experiences []gormmodel.Experience
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&experiences).Error; err != nil {
		return err
	}
	p.Experience = make([]entities.Experience, 0, len(experiences))
	for _, m := range experiences {
		eid, e := valueobjects.IDFromUUID(m.ID)
		if e != nil {
			return e
		}
		start, e := parseYearMonthDate(m.StartDate)
		if e != nil {
			return e
		}
		var end *valueobjects.YearMonth
		if m.EndDate != nil {
			v, e := parseYearMonthDate(*m.EndDate)
			if e != nil {
				return e
			}
			end = &v
		}
		var ach []gormmodel.ExperienceAchievement
		if e := db.WithContext(ctx).Where("experience_id = ?", m.ID).Order("position").Find(&ach).Error; e != nil {
			return e
		}
		achValues := make([]string, 0, len(ach))
		for _, a := range ach {
			achValues = append(achValues, a.Value)
		}
		var stack []gormmodel.ExperienceTechnology
		if e := db.WithContext(ctx).Where("experience_id = ?", m.ID).Order("position").Find(&stack).Error; e != nil {
			return e
		}
		stackValues := make([]entities.Technology, 0, len(stack))
		for _, t := range stack {
			var g *enums.TechnologyGrade
			if t.Grade != nil {
				v := enums.TechnologyGrade(*t.Grade)
				g = &v
			}
			v, e := entities.NewTechnology(entities.TechnologyParams{Name: t.Name, Grade: g})
			if e != nil {
				return e
			}
			stackValues = append(stackValues, v)
		}
		v, e := entities.NewExperience(
			entities.ExperienceParams{
				ID:           eid,
				Company:      m.Company,
				Position:     m.JobPosition,
				Location:     m.Location,
				StartDate:    start,
				EndDate:      end,
				Description:  m.Description,
				Achievements: achValues,
				Stack:        stackValues,
			},
		)
		if e != nil {
			return e
		}
		p.Experience = append(p.Experience, v)
	}
	var education []gormmodel.Education
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&education).Error; err != nil {
		return err
	}
	p.Education = make([]entities.Education, 0, len(education))
	for _, m := range education {
		eid, e := valueobjects.IDFromUUID(m.ID)
		if e != nil {
			return e
		}
		var end *valueobjects.Year
		if m.EndDate != nil {
			v, e := parseYearDate(*m.EndDate)
			if e != nil {
				return e
			}
			end = &v
		}
		v, e := entities.NewEducation(
			entities.EducationParams{
				ID:          eid,
				Location:    m.Location,
				University:  m.University,
				Grade:       enums.EducationLevel(m.Grade),
				ProgramName: m.ProgramName,
				ProfileName: m.ProfileName,
				Faculty:     m.Faculty,
				Course:      m.Course,
				EndDate:     end,
			},
		)
		if e != nil {
			return e
		}
		p.Education = append(p.Education, v)
	}
	var courses []gormmodel.Course
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&courses).Error; err != nil {
		return err
	}
	p.Courses = make([]entities.Course, 0, len(courses))
	for _, m := range courses {
		cid, e := valueobjects.IDFromUUID(m.ID)
		if e != nil {
			return e
		}
		end, e := parseYearMonthDate(m.EndDate)
		if e != nil {
			return e
		}
		var link *valueobjects.ExternalURL
		if m.CertificateLink != nil {
			v, e := valueobjects.NewExternalURL(*m.CertificateLink)
			if e != nil {
				return e
			}
			link = &v
		}
		v, e := entities.NewCourse(
			entities.CourseParams{
				ID:              cid,
				Name:            m.Name,
				Organization:    m.Organization,
				EndDate:         end,
				Hours:           m.Hours,
				CertificateLink: link,
			},
		)
		if e != nil {
			return e
		}
		p.Courses = append(p.Courses, v)
	}
	var activities []gormmodel.Activity
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&activities).Error; err != nil {
		return err
	}
	p.Activities = make([]entities.Activity, 0, len(activities))
	for _, m := range activities {
		aid, e := valueobjects.IDFromUUID(m.ID)
		if e != nil {
			return e
		}
		var date *valueobjects.YearMonth
		if m.Date != nil {
			v, e := parseYearMonthDate(*m.Date)
			if e != nil {
				return e
			}
			date = &v
		}
		var link *valueobjects.ExternalURL
		if m.Link != nil {
			v, e := valueobjects.NewExternalURL(*m.Link)
			if e != nil {
				return e
			}
			link = &v
		}
		v, e := entities.NewActivity(
			entities.ActivityParams{
				ID:          aid,
				Type:        enums.ActivityType(m.Type),
				Title:       m.Title,
				Description: m.Description,
				Date:        date,
				Link:        link,
			},
		)
		if e != nil {
			return e
		}
		p.Activities = append(p.Activities, v)
	}

	var projects []gormmodel.Project
	if err := db.WithContext(ctx).Where("resume_id = ?", id).Order("position").Find(&projects).Error; err != nil {
		return err
	}
	p.Projects = make([]entities.Project, 0, len(projects))
	for _, m := range projects {
		projectID, e := valueobjects.IDFromUUID(m.ID)
		if e != nil {
			return e
		}

		var description, role *string
		description = cloneString(m.Description)
		role = cloneString(m.Role)

		var link *valueobjects.ExternalURL
		if m.Link != nil {
			v, e := valueobjects.NewExternalURL(*m.Link)
			if e != nil {
				return e
			}
			link = &v
		}
		var gitLink *valueobjects.ExternalURL
		if m.GitLink != nil {
			v, e := valueobjects.NewExternalURL(*m.GitLink)
			if e != nil {
				return e
			}
			gitLink = &v
		}

		var screenshots []gormmodel.ProjectScreenshot
		if err := db.WithContext(ctx).Where("project_id = ?", m.ID).Order("position").Find(&screenshots).Error; err != nil {
			return err
		}
		screenshotValues := make([]valueobjects.AssetPath, 0, len(screenshots))
		for _, screenshot := range screenshots {
			v, e := valueobjects.NewAssetPath(screenshot.Path)
			if e != nil {
				return e
			}
			screenshotValues = append(screenshotValues, v)
		}

		var stack []gormmodel.ProjectTechnology
		if err := db.WithContext(ctx).Where("project_id = ?", m.ID).Order("position").Find(&stack).Error; err != nil {
			return err
		}
		stackValues := make([]entities.Technology, 0, len(stack))
		for _, technology := range stack {
			var grade *enums.TechnologyGrade
			if technology.Grade != nil {
				v := enums.TechnologyGrade(*technology.Grade)
				grade = &v
			}
			v, e := entities.NewTechnology(entities.TechnologyParams{Name: technology.Name, Grade: grade})
			if e != nil {
				return e
			}
			stackValues = append(stackValues, v)
		}

		var contributors []gormmodel.ProjectContributor
		if err := db.WithContext(ctx).Where("project_id = ?", m.ID).Order("position").Find(&contributors).Error; err != nil {
			return err
		}
		contributorValues := make([]entities.Contributor, 0, len(contributors))
		for _, contributor := range contributors {
			var telegram *valueobjects.Telegram
			if contributor.Telegram != nil {
				v, e := valueobjects.NewTelegram(*contributor.Telegram)
				if e != nil {
					return e
				}
				telegram = &v
			}
			var email *valueobjects.Email
			if contributor.Email != nil {
				v, e := valueobjects.NewEmail(*contributor.Email)
				if e != nil {
					return e
				}
				email = &v
			}
			contacts := entities.EmptyContacts()
			if contributor.Phone != nil || telegram != nil || email != nil {
				phone := ""
				if contributor.Phone != nil {
					phone = *contributor.Phone
				}
				contacts, e = entities.NewContacts(entities.ContactsParams{Phone: phone, Telegram: telegram, Email: email})
				if e != nil {
					return e
				}
			}
			v, e := entities.NewContributor(entities.ContributorParams{
				Name: contributor.Name, Surname: contributor.Surname, Role: contributor.Role, Contacts: contacts,
			})
			if e != nil {
				return e
			}
			contributorValues = append(contributorValues, v)
		}

		v, e := entities.NewProject(entities.ProjectParams{
			ID: projectID, Name: m.Name, Description: description, Role: role,
			Screenshots: screenshotValues, Link: link, GitLink: gitLink,
			Stack: stackValues, Contributors: contributorValues,
		})
		if e != nil {
			return e
		}
		p.Projects = append(p.Projects, v)
	}
	return nil
}
