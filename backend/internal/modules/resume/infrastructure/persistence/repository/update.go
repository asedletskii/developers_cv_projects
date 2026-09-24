package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

func (r *ResumeRepository) Update(ctx context.Context, command commands.UpdateResumeRepositoryCommand) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("update resume: repository dependencies are nil")
	}
	resume := command.Resume
	if resume == nil {
		return fmt.Errorf("update resume: %w", domain.ErrInvalidResume)
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("update resume: %w", err)
	}
	if err := replaceAggregate(ctx, tx, resume); err != nil {
		return fmt.Errorf("update resume: %w", err)
	}
	return nil
}

func replaceAggregate(ctx context.Context, tx *gorm.DB, resume *domain.Resume) error {
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

func (r *ResumeRepository) UpdateTitle(ctx context.Context, command commands.UpdateResumeTitleRepositoryCommand) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("update resume title: repository dependencies are nil")
	}
	userID := command.UserID
	resumeID := command.ResumeID
	title := command.Title
	if userID.IsZero() {
		return fmt.Errorf("update resume title: invalid user id")
	}
	if resumeID.IsZero() {
		return fmt.Errorf("update resume title: invalid resume id")
	}
	if err := domain.ValidateResumeTitle(title); err != nil {
		return err
	}

	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("update resume title: %w", err)
	}

	result := tx.WithContext(ctx).
		Model(&gormmodel.Resume{}).
		Where("id = ? AND user_id = ?", resumeID.UUID(), userID.UUID()).
		Updates(map[string]any{
			"title":      title,
			"updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("update resume title: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
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
			value := v.From().Time()
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
