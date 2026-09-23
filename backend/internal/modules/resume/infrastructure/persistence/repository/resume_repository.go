package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type ResumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) *ResumeRepository { return &ResumeRepository{db: db} }

func (r *ResumeRepository) GetByID(ctx context.Context, userID useridentity.ID, id valueobjects.ID) (*domain.Resume, error) {
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

func (r *ResumeRepository) Save(ctx context.Context, resume *domain.Resume) error {
	if resume == nil {
		return fmt.Errorf("save resume: %w", domain.ErrInvalidResume)
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("save resume: %w", err)
	}
	if err := r.replaceAggregate(ctx, tx, resume); err != nil {
		return fmt.Errorf("save resume: %w", err)
	}
	return nil
}

func (r *ResumeRepository) replaceAggregate(ctx context.Context, tx *gorm.DB, resume *domain.Resume) error {
	m := toResumeModel(resume)
	if err := tx.WithContext(ctx).Save(&m).Error; err != nil {
		return err
	}

	resumeID := resume.ID().UUID()
	if err := deleteChildren(ctx, tx, resumeID); err != nil {
		return err
	}
	if err := createProfile(ctx, tx, resume); err != nil {
		return err
	}
	if err := createSettings(ctx, tx, resume); err != nil {
		return err
	}
	if err := createContent(ctx, tx, resume); err != nil {
		return err
	}
	return createProjects(ctx, tx, resume)
}

func deleteChildren(ctx context.Context, tx *gorm.DB, resumeID uuid.UUID) error {
	models := []any{
		&gormmodel.Greeting{}, &gormmodel.Availability{}, &gormmodel.Theme{}, &gormmodel.ContactsText{},
		&gormmodel.Contact{}, &gormmodel.Section{}, &gormmodel.ResumeWorkFormat{},
		&gormmodel.ResumeHiddenContact{}, &gormmodel.ResumeHardSkill{}, &gormmodel.ResumeSoftSkill{},
		&gormmodel.ResumeLanguage{}, &gormmodel.ResumeTechnology{}, &gormmodel.Experience{},
		&gormmodel.Education{}, &gormmodel.Course{}, &gormmodel.Activity{}, &gormmodel.Project{},
	}
	for _, model := range models {
		if err := tx.WithContext(ctx).Delete(model, "resume_id = ?", resumeID).Error; err != nil {
			return err
		}
	}
	return nil
}

func createProfile(ctx context.Context, tx *gorm.DB, p *domain.Resume) error {
	g := p.Greeting()
	gm := gormmodel.Greeting{
		ResumeID:          p.ID().UUID(),
		FirstRU:           g.First().RU(),
		FirstEN:           g.First().EN(),
		SecondRU:          g.Second().RU(),
		SecondEN:          g.Second().EN(),
		SecondHighlightRU: g.SecondHighlight().RU(),
		SecondHighlightEN: g.SecondHighlight().EN(),
		ThirdRU:           g.Third().RU(),
		ThirdEN:           g.Third().EN(),
	}
	if err := tx.WithContext(ctx).Create(&gm).Error; err != nil {
		return err
	}
	if v := p.Availability(); v != nil {
		m := gormmodel.Availability{
			ResumeID: p.ID().UUID(),
			Status:   string(v.Status()),
			NoteRU:   localizedPart(v.Note(), "ru"),
			NoteEN:   localizedPart(v.Note(), "en"),
		}
		if v.From() != nil {
			value := v.From().String()
			m.From = &value
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	if v := p.Theme(); v != nil {
		m := gormmodel.Theme{
			ResumeID:  p.ID().UUID(),
			Primary:   colorPtr(v.Primary),
			Secondary: colorPtr(v.Secondary),
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	if v := p.ContactsText(); v != nil {
		m := gormmodel.ContactsText{
			ResumeID:   p.ID().UUID(),
			TitleRU:    localizedPtr(v.Title(), "ru"),
			TitleEN:    localizedPtr(v.Title(), "en"),
			SubtitleRU: localizedPtr(v.Subtitle(), "ru"),
			SubtitleEN: localizedPtr(v.Subtitle(), "en"),
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	c := p.Contacts()
	cm := gormmodel.Contact{ResumeID: p.ID().UUID(), Phone: c.Phone()}
	if v := c.Telegram(); v != nil {
		value := v.String()
		cm.Telegram = &value
	}
	if v := c.Email(); v != nil {
		value := v.String()
		cm.Email = &value
	}
	if cm.Phone == nil && cm.Telegram == nil && cm.Email == nil {
		return nil
	}
	return tx.WithContext(ctx).Create(&cm).Error
}

func createSettings(ctx context.Context, tx *gorm.DB, p *domain.Resume) error {
	pid := p.ID().UUID()
	for i, v := range p.Sections() {
		if err := tx.WithContext(ctx).Create(
			&gormmodel.Section{
				ResumeID:  pid,
				SectionID: string(v.ID()),
				Position:  i,
				Enabled:   v.Enabled(),
				Collapsed: v.Collapsed(),
			}).Error; err != nil {
			return err
		}
	}
	for i, v := range p.WorkFormats() {
		if err := tx.WithContext(ctx).Create(
			&gormmodel.ResumeWorkFormat{
				ResumeID: pid,
				Value:    string(v),
				Position: i,
			}).Error; err != nil {
			return err
		}
	}
	for i, v := range p.HiddenContacts() {
		if err := tx.WithContext(ctx).Create(
			&gormmodel.ResumeHiddenContact{
				ResumeID: pid,
				Value:    string(v),
				Position: i,
			}).Error; err != nil {
			return err
		}
	}
	for i, v := range p.HardSkills() {
		if err := tx.WithContext(ctx).Create(
			&gormmodel.ResumeHardSkill{
				ResumeID: pid,
				Value:    string(v),
				Position: i,
			}).Error; err != nil {
			return err
		}
	}
	for i, v := range p.SoftSkills() {
		if err := tx.WithContext(ctx).Create(
			&gormmodel.ResumeSoftSkill{
				ResumeID: pid,
				Value:    string(v),
				Position: i,
			}).Error; err != nil {
			return err
		}
	}
	for i, v := range p.Languages() {
		m := gormmodel.ResumeLanguage{
			ResumeID: pid,
			Language: string(v.Language()),
			Position: i,
			Grade:    enumPtr(v.Grade()),
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	for i, v := range p.Technologies() {
		m := gormmodel.ResumeTechnology{
			ResumeID: pid,
			Name:     v.Name(),
			Position: i,
			Grade:    enumPtr(v.Grade()),
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	return nil
}

func createContent(ctx context.Context, tx *gorm.DB, p *domain.Resume) error {
	pid := p.ID().UUID()
	for i, v := range p.Experience() {
		m := gormmodel.Experience{
			ID:          v.ID().UUID(),
			ResumeID:    pid,
			Position:    i,
			Company:     v.Company(),
			JobPosition: v.Position(),
			Location:    v.Location(),
			StartDate:   v.StartDate().String() + "-01",
			Description: v.Description(),
		}
		if v.EndDate() != nil {
			value := v.EndDate().String() + "-01"
			m.EndDate = &value
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
		for j, a := range v.Achievements() {
			if err := tx.WithContext(ctx).Create(
				&gormmodel.ExperienceAchievement{
					ExperienceID: m.ID,
					Position:     j,
					Value:        a,
				}).Error; err != nil {
				return err
			}
		}
		for j, t := range v.Stack() {
			if err := tx.WithContext(ctx).Create(
				&gormmodel.ExperienceTechnology{
					ExperienceID: m.ID,
					Position:     j,
					Name:         t.Name(),
					Grade:        enumPtr(t.Grade()),
				}).Error; err != nil {
				return err
			}
		}
	}
	for i, v := range p.Education() {
		m := gormmodel.Education{
			ID:          v.ID().UUID(),
			ResumeID:    pid,
			Position:    i,
			Location:    v.Location(),
			University:  v.University(),
			Grade:       int(v.Grade()),
			ProgramName: v.ProgramName(),
			ProfileName: v.ProfileName(),
			Faculty:     v.Faculty(),
			Course:      v.Course(),
		}
		if v.EndDate() != nil {
			value := v.EndDate().String() + "-01"
			m.EndDate = &value
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	for i, v := range p.Courses() {
		m := gormmodel.Course{
			ID:           v.ID().UUID(),
			ResumeID:     pid,
			Position:     i,
			Name:         v.Name(),
			Organization: v.Organization(),
			EndDate:      v.EndDate().String() + "-01",
			Hours:        v.Hours(),
		}
		if v.CertificateLink() != nil {
			value := v.CertificateLink().String()
			m.CertificateLink = &value
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	for i, v := range p.Activities() {
		m := gormmodel.Activity{
			ID:          v.ID().UUID(),
			ResumeID:    pid,
			Position:    i,
			Type:        string(v.Type()),
			Title:       v.Title(),
			Description: v.Description(),
			Link:        externalURLPtr(v.Link()),
		}
		if v.Date() != nil {
			value := v.Date().String() + "-01"
			m.Date = &value
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
	}
	return nil
}

func createProjects(ctx context.Context, tx *gorm.DB, p *domain.Resume) error {
	pid := p.ID().UUID()
	for i, v := range p.Projects() {
		m := gormmodel.Project{
			ID:          v.ID().UUID(),
			ResumeID:    pid,
			Position:    i,
			Name:        v.Name(),
			Description: v.Description(),
			Role:        v.Role(),
			Link:        externalURLPtr(v.Link()),
			GitLink:     externalURLPtr(v.GitLink()),
		}
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			return err
		}
		for j, s := range v.Screenshots() {
			if err := tx.WithContext(ctx).Create(
				&gormmodel.ProjectScreenshot{
					ProjectID: m.ID,
					Position:  j,
					Path:      s.String(),
				}).Error; err != nil {
				return err
			}
		}
		for j, t := range v.Stack() {
			if err := tx.WithContext(ctx).Create(
				&gormmodel.ProjectTechnology{
					ProjectID: m.ID,
					Position:  j,
					Name:      t.Name(),
					Grade:     enumPtr(t.Grade()),
				}).Error; err != nil {
				return err
			}
		}
		for j, c := range v.Contributors() {
			cm := gormmodel.ProjectContributor{
				ID:        uuid.New(),
				ProjectID: m.ID,
				Position:  j,
				Name:      c.Name(),
				Surname:   c.Surname(),
				Role:      c.Role(),
				Phone:     c.Contacts().Phone(),
			}
			if t := c.Contacts().Telegram(); t != nil {
				value := t.String()
				cm.Telegram = &value
			}
			if e := c.Contacts().Email(); e != nil {
				value := e.String()
				cm.Email = &value
			}
			if err := tx.WithContext(ctx).Create(&cm).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func toResumeModel(p *domain.Resume) gormmodel.Resume {
	now := time.Now().UTC()
	return gormmodel.Resume{
		ID:          p.ID().UUID(),
		UserID:      p.UserID().UUID(),
		Name:        p.Name(),
		Surname:     p.Surname(),
		Patronymic:  p.Patronymic(),
		Position:    p.Position(),
		ResumeLink:  assetPathPtr(p.ResumeLink()),
		PhotoPath:   p.PhotoPath().String(),
		DateOfBirth: p.DateOfBirth().String(),
		Location:    p.Location(),
		Relocation:  p.Relocation(),
		UpdatedAt:   now,
	}
}

func localizedPart(v *valueobjects.LocalizedText, language string) *string {
	if v == nil {
		return nil
	}
	value := v.Value(language)
	if value == "" {
		return nil
	}
	return &value
}
func localizedPtr(v *valueobjects.LocalizedText, language string) *string {
	return localizedPart(v, language)
}
func colorPtr(v *valueobjects.Color) *string {
	if v == nil {
		return nil
	}
	value := v.String()
	return &value
}
func assetPathPtr(v *valueobjects.AssetPath) *string {
	if v == nil {
		return nil
	}
	value := v.String()
	return &value
}
func externalURLPtr(v *valueobjects.ExternalURL) *string {
	if v == nil {
		return nil
	}
	value := v.String()
	return &value
}
func enumPtr[T ~string](v *T) *string {
	if v == nil {
		return nil
	}
	value := string(*v)
	return &value
}

func (r *ResumeRepository) loadAggregate(
	ctx context.Context,
	db *gorm.DB, pm gormmodel.Resume,
) (*domain.Resume, error) {
	id, err := valueobjects.IDFromUUID(pm.ID)
	if err != nil {
		return nil, err
	}
	date, err := valueobjects.ParseDate(pm.DateOfBirth)
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

func newGreeting(m gormmodel.Greeting) (valueobjects.Greeting, error) {
	f, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.FirstRU, EN: m.FirstEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	s, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.SecondRU, EN: m.SecondEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	h, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.SecondHighlightRU, EN: m.SecondHighlightEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	t, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: m.ThirdRU, EN: m.ThirdEN})
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	return valueobjects.NewGreeting(valueobjects.GreetingParams{First: f, Second: s, SecondHighlight: h, Third: t}), nil
}

func (r *ResumeRepository) loadOptionalProfile(ctx context.Context, db *gorm.DB, id uuid.UUID, p *domain.ResumeParams) error {
	var m gormmodel.Availability
	if err := db.WithContext(ctx).First(&m, "resume_id = ?", id).Error; err == nil {
		status := enums.AvailabilityStatus(m.Status)
		var from *valueobjects.Date
		if m.From != nil {
			value, e := valueobjects.ParseDate(*m.From)
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
	return nil
}

func parseYearMonthDate(value string) (valueobjects.YearMonth, error) {
	if len(value) < 7 {
		return valueobjects.YearMonth{}, fmt.Errorf("invalid year-month database value %q", value)
	}
	return valueobjects.ParseYearMonth(value[:7])
}

func parseYearDate(value string) (valueobjects.Year, error) {
	if len(value) < 4 {
		return valueobjects.Year{}, fmt.Errorf("invalid year database value %q", value)
	}
	return valueobjects.ParseYear(value[:4])
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}
