package repository

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"developers-cv/backend/internal/platform/database"

	"gorm.io/gorm"
)

var _ commands.CreateResumeRootRepository = (*ResumeRootRepository)(nil)
var _ commands.GetResumeRootRepository = (*ResumeRootRepository)(nil)
var _ commands.UpdateResumeRootRepository = (*ResumeRootRepository)(nil)
var _ commands.DeleteResumeRootRepository = (*ResumeRootRepository)(nil)

func (r *ResumeRootRepository) CreateRoot(ctx context.Context, c commands.CreateResumeRootRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume root: %w", err)
	}
	m, p, sections, err := rootModels(c.Resume)
	if err != nil {
		return err
	}
	if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("create resume root: %w", err)
	}
	if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
		return fmt.Errorf("create resume presentation: %w", err)
	}
	if len(sections) > 0 {
		if err := tx.WithContext(ctx).Create(&sections).Error; err != nil {
			return fmt.Errorf("create resume presentation sections: %w", err)
		}
	}
	return nil
}
func (r *ResumeRootRepository) GetRoot(ctx context.Context, q commands.GetResumeRootRepositoryQuery) (*domain.ResumeRoot, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume root: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.ResumeRoot
	if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume root: %w", err)
	}
	var p gormmodel.ResumePresentation
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("get resume presentation: %w", err)
		}
		return nil, fmt.Errorf("get resume presentation: %w", err)
	}
	var ss []gormmodel.ResumePresentationSection
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).Order("position").Find(&ss).Error; err != nil {
		return nil, fmt.Errorf("get resume presentation sections: %w", err)
	}
	return mapRoot(m, p, ss)
}
func (r *ResumeRootRepository) UpdateRoot(ctx context.Context, c commands.UpdateResumeRootRepositoryCommand) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("update resume root: %w", err)
	}
	m, p, sections, err := rootModels(c.Resume)
	if err != nil {
		return err
	}
	result := tx.WithContext(ctx).Model(&gormmodel.ResumeRoot{}).Where("id = ? AND user_id = ?", m.ID, m.UserID).Updates(map[string]any{"title": m.Title, "updated_at": m.UpdatedAt})
	if result.Error != nil {
		return fmt.Errorf("update resume root: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	presentationResult := tx.WithContext(ctx).Model(&gormmodel.ResumePresentation{}).Where("resume_id = ?", m.ID).Updates(map[string]any{
		"greeting_first_ru": p.GreetingFirstRU, "greeting_first_en": p.GreetingFirstEN, "greeting_second_ru": p.GreetingSecondRU, "greeting_second_en": p.GreetingSecondEN, "greeting_highlight_ru": p.GreetingHighlightRU, "greeting_highlight_en": p.GreetingHighlightEN, "greeting_third_ru": p.GreetingThirdRU, "greeting_third_en": p.GreetingThirdEN, "availability_status": p.AvailabilityStatus, "availability_from": p.AvailabilityFrom, "availability_note_ru": p.AvailabilityNoteRU, "availability_note_en": p.AvailabilityNoteEN, "resume_link": p.ResumeLink, "theme_primary": p.ThemePrimary, "theme_secondary": p.ThemeSecondary,
	})
	if presentationResult.Error != nil {
		return fmt.Errorf("update resume presentation: %w", presentationResult.Error)
	}
	if presentationResult.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).Delete(&gormmodel.ResumePresentationSection{}).Error; err != nil {
		return fmt.Errorf("replace resume presentation sections: %w", err)
	}
	if len(sections) > 0 {
		if err := tx.WithContext(ctx).Create(&sections).Error; err != nil {
			return fmt.Errorf("create resume presentation sections: %w", err)
		}
	}
	return nil
}
func (r *ResumeRootRepository) DeleteRoot(ctx context.Context, q commands.GetResumeRootRepositoryQuery) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume root: %w", err)
	}
	result := tx.WithContext(ctx).Where("id = ? AND user_id = ?", q.ResumeID.UUID(), q.UserID.UUID()).Delete(&gormmodel.ResumeRoot{})
	if result.Error != nil {
		return fmt.Errorf("delete resume root: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	return nil
}
func rootModels(r *domain.ResumeRoot) (gormmodel.ResumeRoot, gormmodel.ResumePresentation, []gormmodel.ResumePresentationSection, error) {
	if r == nil {
		return gormmodel.ResumeRoot{}, gormmodel.ResumePresentation{}, nil, fmt.Errorf("resume root is nil")
	}
	p := r.Presentation()
	m := gormmodel.ResumeRoot{ID: r.ID().UUID(), UserID: r.UserID().UUID(), Title: r.Title(), CreatedAt: r.CreatedAt(), UpdatedAt: r.UpdatedAt()}
	pm := gormmodel.ResumePresentation{ResumeID: r.ID().UUID(), GreetingFirstRU: p.Greeting.First().RU(), GreetingFirstEN: p.Greeting.First().EN(), GreetingSecondRU: p.Greeting.Second().RU(), GreetingSecondEN: p.Greeting.Second().EN(), GreetingHighlightRU: p.Greeting.SecondHighlight().RU(), GreetingHighlightEN: p.Greeting.SecondHighlight().EN(), GreetingThirdRU: p.Greeting.Third().RU(), GreetingThirdEN: p.Greeting.Third().EN()}
	if p.Availability != nil {
		s := string(p.Availability.Status())
		pm.AvailabilityStatus = &s
		if p.Availability.From() != nil {
			v := p.Availability.From().Time()
			pm.AvailabilityFrom = &v
		}
		if p.Availability.Note() != nil {
			ru := p.Availability.Note().RU()
			en := p.Availability.Note().EN()
			pm.AvailabilityNoteRU = &ru
			pm.AvailabilityNoteEN = &en
		}
	}
	if p.ResumeLink != nil {
		v := p.ResumeLink.String()
		pm.ResumeLink = &v
	}
	if p.Theme != nil {
		if p.Theme.Primary != nil {
			v := p.Theme.Primary.String()
			pm.ThemePrimary = &v
		}
		if p.Theme.Secondary != nil {
			v := p.Theme.Secondary.String()
			pm.ThemeSecondary = &v
		}
	}
	ss := make([]gormmodel.ResumePresentationSection, 0, len(p.Sections))
	for i, s := range p.Sections {
		ss = append(ss, gormmodel.ResumePresentationSection{ResumeID: r.ID().UUID(), SectionID: string(s.ID()), Position: i, Enabled: s.Enabled(), Collapsed: s.Collapsed()})
	}
	return m, pm, ss, nil
}
func mapRoot(m gormmodel.ResumeRoot, p gormmodel.ResumePresentation, ss []gormmodel.ResumePresentationSection) (*domain.ResumeRoot, error) {
	id, e := valueobjects.IDFromUUID(m.ID)
	if e != nil {
		return nil, e
	}
	uid, e := useridentity.IDFromUUID(m.UserID)
	if e != nil {
		return nil, e
	}
	f, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingFirstRU, EN: p.GreetingFirstEN})
	if e != nil {
		return nil, e
	}
	s, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingSecondRU, EN: p.GreetingSecondEN})
	if e != nil {
		return nil, e
	}
	h, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingHighlightRU, EN: p.GreetingHighlightEN})
	if e != nil {
		return nil, e
	}
	t, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: p.GreetingThirdRU, EN: p.GreetingThirdEN})
	if e != nil {
		return nil, e
	}
	g := valueobjects.NewGreeting(valueobjects.GreetingParams{First: f, Second: s, SecondHighlight: h, Third: t})
	var av *valueobjects.Availability
	if p.AvailabilityStatus != nil {
		var from *valueobjects.Date
		if p.AvailabilityFrom != nil {
			d, e := valueobjects.NewDate(*p.AvailabilityFrom)
			if e != nil {
				return nil, e
			}
			from = &d
		}
		var note *valueobjects.LocalizedText
		if p.AvailabilityNoteRU != nil || p.AvailabilityNoteEN != nil {
			ru, en := "", ""
			if p.AvailabilityNoteRU != nil {
				ru = *p.AvailabilityNoteRU
			}
			if p.AvailabilityNoteEN != nil {
				en = *p.AvailabilityNoteEN
			}
			v, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: ru, EN: en})
			if e != nil {
				return nil, e
			}
			note = &v
		}
		v, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: enums.AvailabilityStatus(*p.AvailabilityStatus), From: from, Note: note})
		if e != nil {
			return nil, e
		}
		av = &v
	}
	var theme *valueobjects.ThemeSettings
	if p.ThemePrimary != nil || p.ThemeSecondary != nil {
		var a, b *valueobjects.Color
		if p.ThemePrimary != nil {
			v, e := valueobjects.NewColor(*p.ThemePrimary)
			if e != nil {
				return nil, e
			}
			a = &v
		}
		if p.ThemeSecondary != nil {
			v, e := valueobjects.NewColor(*p.ThemeSecondary)
			if e != nil {
				return nil, e
			}
			b = &v
		}
		v, e := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: a, Secondary: b})
		if e != nil {
			return nil, e
		}
		theme = &v
	}
	var link *valueobjects.AssetPath
	if p.ResumeLink != nil {
		v, e := valueobjects.NewAssetPath(*p.ResumeLink)
		if e != nil {
			return nil, e
		}
		link = &v
	}
	sections := make([]valueobjects.SectionSettings, 0, len(ss))
	for _, x := range ss {
		v, e := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionID(x.SectionID), Enabled: x.Enabled, Collapsed: x.Collapsed})
		if e != nil {
			return nil, e
		}
		sections = append(sections, v)
	}
	return domain.NewResumeRoot(domain.ResumeRootParams{ID: id, UserID: uid, Title: m.Title, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, Presentation: domain.ResumePresentation{Theme: theme, Greeting: g, Availability: av, ResumeLink: link, Sections: sections}})
}
