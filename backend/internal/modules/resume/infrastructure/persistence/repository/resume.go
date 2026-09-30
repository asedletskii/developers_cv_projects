package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *ResumeRepository) Get(ctx context.Context, params ports.GetResumeParams) (*domain.Resume, error) {
	return r.loadResume(ctx, params.ResumeID, params.UserID)
}

func (r *ResumeRepository) Create(ctx context.Context, resume *domain.Resume) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume: %w", err)
	}
	if resume == nil {
		return fmt.Errorf("resume is nil")
	}
	m := gormmodel.Resume{ID: resume.ID().UUID(), UserID: resume.UserID().UUID(), Title: resume.Title(), Language: string(resume.Language()), Status: string(resume.Status()), Version: resume.Version(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("create resume: %w", mapDBError(err))
	}
	if resume.Profile() != nil {
		pm, fs, err := profileModels(resume.Profile())
		if err != nil {
			return fmt.Errorf("create resume profile: %w", err)
		}
		if err := tx.WithContext(ctx).Create(&pm).Error; err != nil {
			return fmt.Errorf("create resume profile: %w", mapDBError(err))
		}
		if len(fs) > 0 {
			if err := tx.WithContext(ctx).Create(&fs).Error; err != nil {
				return fmt.Errorf("create resume profile work formats: %w", mapDBError(err))
			}
		}
	}
	if resume.Overview() != nil {
		p, sections, err := overviewModels(resume.ID(), resume.Overview())
		if err != nil {
			return fmt.Errorf("create resume overview: %w", err)
		}
		if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
			return fmt.Errorf("create resume overview: %w", mapDBError(err))
		}
		if len(sections) > 0 {
			if err := tx.WithContext(ctx).Create(&sections).Error; err != nil {
				return fmt.Errorf("create resume overview sections: %w", mapDBError(err))
			}
		}
	}
	if resume.Contacts() != nil {
		if err := persistContacts(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume contacts: %w", err)
		}
	}
	if resume.Skills() != nil {
		if err := persistSkills(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume skills: %w", err)
		}
	}
	if resume.Languages() != nil {
		if err := persistLanguages(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume languages: %w", err)
		}
	}
	if resume.Technologies() != nil {
		if err := persistTechnologies(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume technologies: %w", err)
		}
	}
	if len(resume.Experiences()) > 0 {
		if err := persistExperiences(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume experiences: %w", err)
		}
	}
	if len(resume.Educations()) > 0 {
		if err := persistEducations(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume educations: %w", err)
		}
	}
	if len(resume.Courses()) > 0 {
		if err := persistCourses(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume courses: %w", err)
		}
	}
	if len(resume.Activities()) > 0 {
		if err := persistActivities(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume activities: %w", err)
		}
	}
	if len(resume.Projects()) > 0 {
		if err := persistProjects(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume projects: %w", err)
		}
	}
	return nil
}

func (r *ResumeRepository) Save(ctx context.Context, resume *domain.Resume) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("save resume: repository dependencies are nil")
	}
	if resume == nil || resume.ID().IsZero() || resume.UserID().IsZero() || resume.Version() < 1 {
		return fmt.Errorf("save resume: invalid aggregate or expected version")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("save resume: %w", err)
	}
	result := tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id = ? AND user_id = ? AND version = ?", resume.ID().UUID(), resume.UserID().UUID(), resume.Version()).Updates(map[string]any{"title": resume.Title(), "language": string(resume.Language()), "status": string(resume.Status()), "version": gorm.Expr("version + 1"), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return fmt.Errorf("save resume: %w", mapDBError(result.Error))
	}
	if result.RowsAffected != 1 {
		return domain.ErrConcurrencyConflict
	}
	if err := persistOverview(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistProfile(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistContacts(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistSkills(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistLanguages(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistTechnologies(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistExperiences(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistEducations(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistCourses(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistActivities(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistProjects(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	return nil
}

func persistOverview(tx *gorm.DB, resume *domain.Resume) error {
	o := resume.Overview()
	if o == nil {
		return fmt.Errorf("save resume overview: overview is required")
	}
	p, sections, err := overviewModels(resume.ID(), o)
	if err != nil {
		return fmt.Errorf("save resume overview: %w", err)
	}
	result := tx.Model(&gormmodel.ResumeOverview{}).Where("resume_id = ?", resume.ID().UUID()).Updates(map[string]any{
		"greeting_first": p.GreetingFirst, "greeting_second": p.GreetingSecond, "greeting_highlight": p.GreetingHighlight, "greeting_third": p.GreetingThird,
		"availability_status": p.AvailabilityStatus, "availability_from": p.AvailabilityFrom, "availability_note": p.AvailabilityNote, "resume_link": p.ResumeLink, "theme_primary": p.ThemePrimary, "theme_secondary": p.ThemeSecondary,
	})
	if result.Error != nil {
		return fmt.Errorf("update resume overview: %w", mapDBError(result.Error))
	}
	if result.RowsAffected == 0 {
		if err := tx.Create(&p).Error; err != nil {
			return fmt.Errorf("create resume overview: %w", mapDBError(err))
		}
	}
	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeOverviewSection{}).Error; err != nil {
		return fmt.Errorf("replace resume overview sections: %w", mapDBError(err))
	}
	if len(sections) > 0 {
		if err := tx.Create(&sections).Error; err != nil {
			return fmt.Errorf("create resume overview sections: %w", mapDBError(err))
		}
	}
	return nil
}

func persistProfile(tx *gorm.DB, resume *domain.Resume) error {
	p := resume.Profile()
	if p == nil {
		var ids []string
		if err := tx.Model(&gormmodel.ResumeProfile{}).Where("resume_id = ?", resume.ID().UUID()).Pluck("id", &ids).Error; err != nil {
			return fmt.Errorf("load resume profile ids: %w", mapDBError(err))
		}
		if len(ids) > 0 {
			if err := tx.Where("profile_id IN ?", ids).Delete(&gormmodel.ResumeProfileWorkFormat{}).Error; err != nil {
				return fmt.Errorf("delete resume profile work formats: %w", mapDBError(err))
			}
		}
		if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeProfile{}).Error; err != nil {
			return fmt.Errorf("delete resume profile: %w", mapDBError(err))
		}
		return nil
	}
	m, formats, err := profileModels(p)
	if err != nil {
		return fmt.Errorf("map resume profile: %w", err)
	}
	var current gormmodel.ResumeProfile
	err = tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume profile: %w", mapDBError(err))
		}
	} else if err != nil {
		return fmt.Errorf("load resume profile: %w", mapDBError(err))
	} else {
		if err := tx.Model(&gormmodel.ResumeProfile{}).Where("id = ?", current.ID).Updates(map[string]any{"name": m.Name, "surname": m.Surname, "patronymic": m.Patronymic, "position": m.Position, "photo_path": m.PhotoPath, "date_of_birth": m.DateOfBirth, "location": m.Location, "relocation": m.Relocation}).Error; err != nil {
			return fmt.Errorf("update resume profile: %w", mapDBError(err))
		}
		m.ID = current.ID
	}
	if err := tx.Where("profile_id = ?", m.ID).Delete(&gormmodel.ResumeProfileWorkFormat{}).Error; err != nil {
		return fmt.Errorf("replace resume profile work formats: %w", mapDBError(err))
	}
	for i := range formats {
		formats[i].ProfileID = m.ID
	}
	if len(formats) > 0 {
		if err := tx.Create(&formats).Error; err != nil {
			return fmt.Errorf("create resume profile work formats: %w", mapDBError(err))
		}
	}
	return nil
}

func persistContacts(tx *gorm.DB, resume *domain.Resume) error {
	contacts := resume.Contacts()
	if contacts == nil {
		var current gormmodel.ResumeContacts
		err := tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("load resume contacts: %w", mapDBError(err))
		}
		if err := tx.Where("contacts_id = ?", current.ID).Delete(&gormmodel.ResumeContactHidden{}).Error; err != nil {
			return fmt.Errorf("delete resume contact hidden: %w", mapDBError(err))
		}
		if err := tx.Where("contacts_id = ?", current.ID).Delete(&gormmodel.ResumeContactsText{}).Error; err != nil {
			return fmt.Errorf("delete resume contacts text: %w", mapDBError(err))
		}
		if err := tx.Delete(&current).Error; err != nil {
			return fmt.Errorf("delete resume contacts: %w", mapDBError(err))
		}
		return nil
	}

	m, hidden, text, err := contactsModels(contacts)
	if err != nil {
		return fmt.Errorf("map resume contacts: %w", err)
	}
	var current gormmodel.ResumeContacts
	err = tx.Where("resume_id = ?", resume.ID().UUID()).First(&current).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&m).Error; err != nil {
			return fmt.Errorf("create resume contacts: %w", mapDBError(err))
		}
	} else if err != nil {
		return fmt.Errorf("load resume contacts: %w", mapDBError(err))
	} else {
		if current.ID != m.ID {
			return fmt.Errorf("%w: contacts identity mismatch", domain.ErrResumeInvariant)
		}
		if err := tx.Model(&gormmodel.ResumeContacts{}).Where("id = ? AND resume_id = ?", current.ID, resume.ID().UUID()).Updates(map[string]any{"phone": m.Phone, "telegram": m.Telegram, "email": m.Email}).Error; err != nil {
			return fmt.Errorf("update resume contacts: %w", mapDBError(err))
		}
	}
	if err := tx.Where("contacts_id = ?", m.ID).Delete(&gormmodel.ResumeContactHidden{}).Error; err != nil {
		return fmt.Errorf("replace resume contact hidden: %w", mapDBError(err))
	}
	if len(hidden) > 0 {
		if err := tx.Create(&hidden).Error; err != nil {
			return fmt.Errorf("create resume contact hidden: %w", mapDBError(err))
		}
	}
	if err := tx.Where("contacts_id = ?", m.ID).Delete(&gormmodel.ResumeContactsText{}).Error; err != nil {
		return fmt.Errorf("replace resume contacts text: %w", mapDBError(err))
	}
	if text != nil {
		if err := tx.Create(text).Error; err != nil {
			return fmt.Errorf("create resume contacts text: %w", mapDBError(err))
		}
	}
	return nil
}

func (r *ResumeRepository) Delete(ctx context.Context, params ports.DeleteResumeParams) error {
	resumeID, userID := params.ResumeID, params.UserID
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume: %w", err)
	}
	res := tx.WithContext(ctx).Where("id = ? AND user_id = ?", resumeID.UUID(), userID.UUID()).Delete(&gormmodel.Resume{})
	if res.Error != nil {
		return fmt.Errorf("delete resume: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}

func (r *ResumeRepository) loadResume(ctx context.Context, resumeID valueobjects.ID, userID useridentity.ID) (*domain.Resume, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.Resume
	if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", resumeID.UUID(), userID.UUID()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume: %w", err)
	}
	id, err := valueobjects.IDFromUUID(m.ID)
	if err != nil {
		return nil, err
	}
	uid, err := useridentity.IDFromUUID(m.UserID)
	if err != nil {
		return nil, err
	}
	var profileModel gormmodel.ResumeProfile
	var profile *domain.Profile
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&profileModel).Error; err == nil {
		var formats []gormmodel.ResumeProfileWorkFormat
		if err := tx.WithContext(ctx).Where("profile_id = ?", profileModel.ID).Order("position").Find(&formats).Error; err != nil {
			return nil, fmt.Errorf("get resume profile work formats: %w", err)
		}
		profile, err = mapProfile(profileModel, formats)
		if err != nil {
			return nil, err
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get resume profile: %w", err)
	}
	var contactsModel gormmodel.ResumeContacts
	var contacts *domain.Contacts
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&contactsModel).Error; err == nil {
		var hidden []gormmodel.ResumeContactHidden
		if err := tx.WithContext(ctx).Where("contacts_id = ?", contactsModel.ID).Order("position").Find(&hidden).Error; err != nil {
			return nil, fmt.Errorf("get resume contact hidden: %w", err)
		}
		var textModel gormmodel.ResumeContactsText
		textErr := tx.WithContext(ctx).Where("contacts_id = ?", contactsModel.ID).First(&textModel).Error
		var text *gormmodel.ResumeContactsText
		if textErr == nil {
			text = &textModel
		} else if !errors.Is(textErr, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get resume contacts text: %w", textErr)
		}
		contacts, err = mapContacts(contactsModel, hidden, text, m.Language)
		if err != nil {
			return nil, err
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get resume contacts: %w", err)
	}
	skills, err := loadSkills(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	languages, err := loadLanguages(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	technologies, err := loadTechnologies(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	experiences, err := loadExperiences(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	educations, err := loadEducations(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	courses, err := loadCourses(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	activities, err := loadActivities(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	projects, err := loadProjects(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	var overviewModel gormmodel.ResumeOverview
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&overviewModel).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get resume overview: %w", err)
		}
	} else {
		var sections []gormmodel.ResumeOverviewSection
		if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).Order("position").Find(&sections).Error; err != nil {
			return nil, fmt.Errorf("get resume overview sections: %w", err)
		}
		o, err := mapOverview(overviewModel, sections)
		if err != nil {
			return nil, err
		}
		resume, err := domain.ReconstructResume(domain.ResumeSnapshot{ID: id, UserID: uid, Title: m.Title, Language: enums.ResumeLanguage(m.Language), Status: enums.ResumeStatus(m.Status), Version: m.Version, Overview: o, Profile: profile, Contacts: contacts, Skills: skills, Languages: languages, Technologies: technologies, Experiences: experiences, Educations: educations, Courses: courses, Activities: activities, Projects: projects})
		if err != nil {
			return nil, err
		}
		return resume, nil
	}
	return domain.ReconstructResume(domain.ResumeSnapshot{ID: id, UserID: uid, Title: m.Title, Language: enums.ResumeLanguage(m.Language), Status: enums.ResumeStatus(m.Status), Version: m.Version, Profile: profile, Contacts: contacts, Skills: skills, Languages: languages, Technologies: technologies, Experiences: experiences, Educations: educations, Courses: courses, Activities: activities, Projects: projects})
}

func contactsModels(c *domain.Contacts) (gormmodel.ResumeContacts, []gormmodel.ResumeContactHidden, *gormmodel.ResumeContactsText, error) {
	if c == nil {
		return gormmodel.ResumeContacts{}, nil, nil, fmt.Errorf("contacts is nil")
	}
	d := c.Contacts()
	m := gormmodel.ResumeContacts{ID: c.ID().UUID(), ResumeID: c.ResumeID().UUID()}
	if v := d.Phone(); v != nil {
		m.Phone = v
	}
	if v := d.Telegram(); v != nil {
		x := v.String()
		m.Telegram = &x
	}
	if v := d.Email(); v != nil {
		x := v.String()
		m.Email = &x
	}
	hiddenValues := c.HiddenContacts()
	h := make([]gormmodel.ResumeContactHidden, 0, len(hiddenValues))
	for i, v := range hiddenValues {
		h = append(h, gormmodel.ResumeContactHidden{ContactsID: m.ID, Value: string(v), Position: i})
	}
	var text *gormmodel.ResumeContactsText
	if x := c.Text(); x != nil {
		text = &gormmodel.ResumeContactsText{ContactsID: m.ID}
		if v := x.Title(); v != nil {
			value := *v
			text.Title = &value
		}
		if v := x.Subtitle(); v != nil {
			value := *v
			text.Subtitle = &value
		}
	}
	return m, h, text, nil
}

func mapContacts(m gormmodel.ResumeContacts, h []gormmodel.ResumeContactHidden, t *gormmodel.ResumeContactsText, _ string) (*domain.Contacts, error) {
	var tg *valueobjects.Telegram
	if m.Telegram != nil {
		x, err := valueobjects.NewTelegram(*m.Telegram)
		if err != nil {
			return nil, err
		}
		tg = &x
	}
	var em *valueobjects.Email
	if m.Email != nil {
		x, err := valueobjects.NewEmail(*m.Email)
		if err != nil {
			return nil, err
		}
		em = &x
	}
	contactValues := entities.ContactsParams{Phone: pointerStringValue(m.Phone), Telegram: tg, Email: em}
	contactEntity, err := entities.NewContacts(contactValues)
	if err != nil {
		return nil, err
	}
	hidden := make([]enums.ContactKey, len(h))
	for i, v := range h {
		hidden[i] = enums.ContactKey(v.Value)
	}
	var text *valueobjects.ContactsText
	if t != nil {
		x, err := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: t.Title, Subtitle: t.Subtitle})
		if err != nil {
			return nil, err
		}
		text = &x
	}
	id, err := valueobjects.IDFromUUID(m.ID)
	if err != nil {
		return nil, err
	}
	rid, err := valueobjects.IDFromUUID(m.ResumeID)
	if err != nil {
		return nil, err
	}
	return domain.NewContacts(domain.ContactsParams{ID: id, ResumeID: rid, Contacts: contactEntity, HiddenContacts: hidden, Text: text})
}

func pointerStringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func overviewModels(id valueobjects.ID, o *domain.ResumeOverview) (gormmodel.ResumeOverview, []gormmodel.ResumeOverviewSection, error) {
	if o == nil {
		return gormmodel.ResumeOverview{}, nil, fmt.Errorf("overview is nil")
	}
	theme := o.Theme()
	greeting := o.Greeting()
	availability := o.Availability()
	resumeLink := o.ResumeLink()
	sections := o.Sections()
	m := gormmodel.ResumeOverview{ResumeID: id.UUID(), GreetingFirst: greeting.First(), GreetingSecond: greeting.Second(), GreetingHighlight: greeting.SecondHighlight(), GreetingThird: greeting.Third()}
	if availability != nil {
		v := string(availability.Status())
		m.AvailabilityStatus = &v
		if d := availability.From(); d != nil {
			x := d.Time()
			m.AvailabilityFrom = &x
		}
		if n := availability.Note(); n != nil {
			v := *n
			m.AvailabilityNote = &v
		}
	}
	if resumeLink != nil {
		v := resumeLink.String()
		m.ResumeLink = &v
	}
	if theme != nil {
		if theme.Primary != nil {
			v := theme.Primary.String()
			m.ThemePrimary = &v
		}
		if theme.Secondary != nil {
			v := theme.Secondary.String()
			m.ThemeSecondary = &v
		}
	}
	mapped := make([]gormmodel.ResumeOverviewSection, 0, len(sections))
	for i, x := range sections {
		mapped = append(mapped, gormmodel.ResumeOverviewSection{ResumeID: id.UUID(), SectionID: string(x.ID()), Position: i, Enabled: x.Enabled(), Collapsed: x.Collapsed()})
	}
	return m, mapped, nil
}

func mapOverview(m gormmodel.ResumeOverview, s []gormmodel.ResumeOverviewSection) (*domain.ResumeOverview, error) {
	g := valueobjects.NewGreeting(valueobjects.GreetingParams{First: m.GreetingFirst, Second: m.GreetingSecond, SecondHighlight: m.GreetingHighlight, Third: m.GreetingThird})
	var av *valueobjects.Availability
	if m.AvailabilityStatus != nil {
		st := enums.AvailabilityStatus(*m.AvailabilityStatus)
		params := valueobjects.AvailabilityParams{Status: st, Note: m.AvailabilityNote}
		if m.AvailabilityFrom != nil {
			d, err := valueobjects.NewDate(*m.AvailabilityFrom)
			if err != nil {
				return nil, err
			}
			params.From = &d
		}
		v, err := valueobjects.NewAvailability(params)
		if err != nil {
			return nil, err
		}
		av = &v
	}
	var link *valueobjects.ExternalURL
	if m.ResumeLink != nil {
		v, err := valueobjects.NewExternalURL(*m.ResumeLink)
		if err != nil {
			return nil, err
		}
		link = &v
	}
	var theme *valueobjects.ThemeSettings
	if m.ThemePrimary != nil || m.ThemeSecondary != nil {
		var a, b *valueobjects.Color
		if m.ThemePrimary != nil {
			v, err := valueobjects.NewColor(*m.ThemePrimary)
			if err != nil {
				return nil, err
			}
			a = &v
		}
		if m.ThemeSecondary != nil {
			v, err := valueobjects.NewColor(*m.ThemeSecondary)
			if err != nil {
				return nil, err
			}
			b = &v
		}
		v, err := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: a, Secondary: b})
		if err != nil {
			return nil, err
		}
		theme = &v
	}
	ss := make([]valueobjects.SectionSettings, 0, len(s))
	for _, x := range s {
		v, err := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionID(x.SectionID), Enabled: x.Enabled, Collapsed: x.Collapsed})
		if err != nil {
			return nil, err
		}
		ss = append(ss, v)
	}
	return domain.NewResumeOverview(domain.ResumeOverviewParams{Theme: theme, Greeting: g, Availability: av, ResumeLink: link, Sections: ss})
}

func profileModels(p *domain.Profile) (gormmodel.ResumeProfile, []gormmodel.ResumeProfileWorkFormat, error) {
	if p == nil {
		return gormmodel.ResumeProfile{}, nil, fmt.Errorf("profile is nil")
	}
	m := gormmodel.ResumeProfile{ID: p.ID().UUID(), ResumeID: p.ResumeID().UUID(), Name: p.Name(), Surname: p.Surname(), Patronymic: p.Patronymic(), Position: p.Position(), PhotoPath: p.PhotoPath().String(), DateOfBirth: p.DateOfBirth().Time(), Location: p.Location(), Relocation: p.Relocation()}
	fs := make([]gormmodel.ResumeProfileWorkFormat, 0, len(p.WorkFormats()))
	for i, v := range p.WorkFormats() {
		fs = append(fs, gormmodel.ResumeProfileWorkFormat{ProfileID: m.ID, Value: string(v), Position: i})
	}
	return m, fs, nil
}
func mapProfile(m gormmodel.ResumeProfile, fs []gormmodel.ResumeProfileWorkFormat) (*domain.Profile, error) {
	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return nil, e
	}
	rid, e := valueobjects.IDFromUUID(m.ResumeID)
	if e != nil {
		return nil, e
	}
	photo, e := valueobjects.NewExternalURL(m.PhotoPath)
	if e != nil {
		return nil, e
	}
	dob, e := valueobjects.NewDate(m.DateOfBirth)
	if e != nil {
		return nil, e
	}
	formats := make([]enums.WorkFormat, len(fs))
	for i, v := range fs {
		formats[i] = enums.WorkFormat(v.Value)
	}
	return domain.NewProfile(domain.ProfileParams{ID: id, ResumeID: rid, Name: m.Name, Surname: m.Surname, Patronymic: m.Patronymic, Position: m.Position, PhotoPath: photo, DateOfBirth: dob, Location: m.Location, Relocation: m.Relocation, WorkFormats: formats})
}

func loadExperiences(tx *gorm.DB, resumeID valueobjects.ID) ([]*domain.Experience, error) {
	var models []gormmodel.ResumeExperience
	if err := tx.Where("resume_id = ?", resumeID.UUID()).Order("sort_order").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("load resume experiences: %w", mapDBError(err))
	}
	items := make([]*domain.Experience, 0, len(models))
	for _, model := range models {
		item, err := mapExperience(tx, model)
		if err != nil {
			return nil, fmt.Errorf("load resume experience: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func persistExperiences(tx *gorm.DB, resume *domain.Resume) error {
	models := make([]gormmodel.ResumeExperience, 0, len(resume.Experiences()))
	items := resume.Experiences()
	for position, item := range items {
		if item == nil {
			return fmt.Errorf("%w: nil experience", domain.ErrResumeInvariant)
		}
		if item.ResumeID() != resume.ID() {
			return fmt.Errorf("%w: experience does not belong to resume", domain.ErrResumeInvariant)
		}
		model, err := experienceModel(item)
		if err != nil {
			return fmt.Errorf("map resume experience: %w", err)
		}
		model.SortOrder = position
		models = append(models, model)
	}

	if err := tx.Where("resume_id = ?", resume.ID().UUID()).Delete(&gormmodel.ResumeExperience{}).Error; err != nil {
		return fmt.Errorf("replace resume experiences: %w", mapDBError(err))
	}
	for i, model := range models {
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("create resume experience: %w", mapDBError(err))
		}
		if err := replaceExperienceChildren(tx, model.ID, items[i]); err != nil {
			return fmt.Errorf("replace resume experience children: %w", err)
		}
	}
	return nil
}

func experienceModel(item *domain.Experience) (gormmodel.ResumeExperience, error) {
	if item == nil {
		return gormmodel.ResumeExperience{}, fmt.Errorf("experience is nil")
	}
	if item.ID().IsZero() || item.ResumeID().IsZero() {
		return gormmodel.ResumeExperience{}, fmt.Errorf("experience identity is required")
	}
	model := gormmodel.ResumeExperience{ID: item.ID().UUID(), ResumeID: item.ResumeID().UUID(), Company: item.Company(), Position: item.Position(), Location: item.Location(), StartDate: yearMonthDate(item.StartDate()), Description: item.Description()}
	if item.EndDate() != nil {
		value := yearMonthDate(*item.EndDate())
		model.EndDate = &value
	}
	return model, nil
}

func replaceExperienceChildren(tx *gorm.DB, experienceID uuid.UUID, item *domain.Experience) error {
	if item == nil {
		return fmt.Errorf("experience is required")
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceAchievement{}).Error; err != nil {
		return err
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceTask{}).Error; err != nil {
		return err
	}
	if err := tx.Where("experience_id = ?", experienceID).Delete(&gormmodel.ResumeExperienceStack{}).Error; err != nil {
		return err
	}
	for i, achievement := range item.Achievements() {
		if err := tx.Create(&gormmodel.ResumeExperienceAchievement{ExperienceID: experienceID, Position: i, Value: achievement}).Error; err != nil {
			return err
		}
	}
	for i, task := range item.Tasks() {
		if err := tx.Create(&gormmodel.ResumeExperienceTask{ExperienceID: experienceID, Position: i, Value: task}).Error; err != nil {
			return err
		}
	}
	for i, technology := range item.Stack() {
		var grade *string
		if value := technology.Grade(); value != nil {
			text := string(*value)
			grade = &text
		}
		if err := tx.Create(&gormmodel.ResumeExperienceStack{ExperienceID: experienceID, Position: i, Name: technology.Name(), Grade: grade}).Error; err != nil {
			return err
		}
	}
	return nil
}

func mapExperience(tx *gorm.DB, model gormmodel.ResumeExperience) (*domain.Experience, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, err
	}
	resumeID, err := valueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, err
	}
	start, err := valueobjects.ParseYearMonth(model.StartDate.Format("2006-01"))
	if err != nil {
		return nil, err
	}
	var end *valueobjects.YearMonth
	if model.EndDate != nil {
		value, err := valueobjects.ParseYearMonth(model.EndDate.Format("2006-01"))
		if err != nil {
			return nil, err
		}
		end = &value
	}
	var achievements []gormmodel.ResumeExperienceAchievement
	if err := tx.Where("experience_id = ?", model.ID).Order("position").Find(&achievements).Error; err != nil {
		return nil, err
	}
	var tasks []gormmodel.ResumeExperienceTask
	if err := tx.Where("experience_id = ?", model.ID).Order("position").Find(&tasks).Error; err != nil {
		return nil, err
	}
	var stack []gormmodel.ResumeExperienceStack
	if err := tx.Where("experience_id = ?", model.ID).Order("position").Find(&stack).Error; err != nil {
		return nil, err
	}
	stackValues := make([]domain.ExperienceStackParams, 0, len(stack))
	for _, technology := range stack {
		var grade *enums.TechnologyGrade
		if technology.Grade != nil {
			value := enums.TechnologyGrade(*technology.Grade)
			grade = &value
		}
		stackValues = append(stackValues, domain.ExperienceStackParams{Name: technology.Name, Grade: grade})
	}
	achievementValues := make([]string, len(achievements))
	for i, achievement := range achievements {
		achievementValues[i] = achievement.Value
	}
	taskValues := make([]string, len(tasks))
	for i, task := range tasks {
		taskValues[i] = task.Value
	}
	return domain.NewExperience(domain.ExperienceParams{
		ID:           id,
		ResumeID:     resumeID,
		Company:      model.Company,
		Position:     model.Position,
		Location:     model.Location,
		StartDate:    start,
		EndDate:      end,
		Description:  model.Description,
		Achievements: achievementValues,
		Tasks:        taskValues,
		Stack:        stackValues,
	})
}
