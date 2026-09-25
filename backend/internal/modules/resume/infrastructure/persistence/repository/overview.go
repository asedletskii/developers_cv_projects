package repository

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var _ ports.CreateOverviewRepository = (*OverviewRepository)(nil)
var _ ports.GetOverviewRepository = (*OverviewRepository)(nil)
var _ ports.UpdateOverviewRepository = (*OverviewRepository)(nil)

type OverviewRepository struct{ *baseRepository }

func NewOverviewRepository(db *gorm.DB) *OverviewRepository {
	return &OverviewRepository{newBaseRepository(db)}
}
func (r *OverviewRepository) CreateOverview(ctx context.Context, c ports.CreateOverviewRepositoryCommand) error {
	tx, e := database.DBFromTransactionContext(ctx)
	if e != nil {
		return fmt.Errorf("create overview: %w", e)
	}
	if c.Overview == nil {
		return fmt.Errorf("overview is nil")
	}
	p, s, e := overviewModels(c.ResumeID, c.Overview)
	if e != nil {
		return e
	}
	var count int64
	if e = tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id = ? AND user_id = ?", c.ResumeID.UUID(), c.UserID.UUID()).Count(&count).Error; e != nil {
		return fmt.Errorf("create overview ownership: %w", e)
	}
	if count == 0 {
		return domain.ErrEntityNotFound
	}
	if e = tx.WithContext(ctx).Create(&p).Error; e != nil {
		if errors.Is(e, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
		return fmt.Errorf("create overview: %w", e)
	}
	if len(s) > 0 {
		if e = tx.WithContext(ctx).Create(&s).Error; e != nil {
			return fmt.Errorf("create overview sections: %w", e)
		}
	}
	return nil
}
func (r *OverviewRepository) GetOverview(ctx context.Context, c ports.GetOverviewRepositoryCommand) (*domain.ResumeOverview, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get overview: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var p gormmodel.ResumeOverview
	if e := tx.WithContext(ctx).Raw("SELECT p.* FROM resume_overviews p JOIN resumes r ON r.id=p.resume_id WHERE p.resume_id=? AND r.user_id=?", c.ResumeID.UUID(), c.UserID.UUID()).Scan(&p).Error; e != nil {
		return nil, fmt.Errorf("get overview: %w", e)
	}
	if p.ResumeID == [16]byte{} {
		return nil, domain.ErrEntityNotFound
	}
	var s []gormmodel.ResumeOverviewSection
	if e := tx.WithContext(ctx).Where("resume_id=?", c.ResumeID.UUID()).Order("position").Find(&s).Error; e != nil {
		return nil, fmt.Errorf("get overview sections: %w", e)
	}
	return mapOverview(p, s)
}
func (r *OverviewRepository) UpdateOverview(ctx context.Context, c ports.UpdateOverviewRepositoryCommand) error {
	tx, e := database.DBFromTransactionContext(ctx)
	if e != nil {
		return fmt.Errorf("update overview: %w", e)
	}
	p, s, e := overviewModels(c.ResumeID, c.Overview)
	if e != nil {
		return e
	}
	res := tx.WithContext(ctx).Model(&gormmodel.ResumeOverview{}).
		Where("resume_id = ? AND resume_id IN (SELECT id FROM resumes WHERE id = ? AND user_id = ?)", c.ResumeID.UUID(), c.ResumeID.UUID(), c.UserID.UUID()).
		Updates(map[string]any{"greeting_first": p.GreetingFirst, "greeting_second": p.GreetingSecond, "greeting_highlight": p.GreetingHighlight, "greeting_third": p.GreetingThird, "availability_status": p.AvailabilityStatus, "availability_from": p.AvailabilityFrom, "availability_note": p.AvailabilityNote, "resume_link": p.ResumeLink, "theme_primary": p.ThemePrimary, "theme_secondary": p.ThemeSecondary})
	if res.Error != nil {
		return fmt.Errorf("update overview: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	if e = tx.WithContext(ctx).Where("resume_id=?", c.ResumeID.UUID()).Delete(&gormmodel.ResumeOverviewSection{}).Error; e != nil {
		return fmt.Errorf("replace overview sections: %w", e)
	}
	if len(s) > 0 {
		if e = tx.WithContext(ctx).Create(&s).Error; e != nil {
			return fmt.Errorf("create overview sections: %w", e)
		}
	}
	return nil
}
func overviewModels(id valueobjects.ID, o *domain.ResumeOverview) (gormmodel.ResumeOverview, []gormmodel.ResumeOverviewSection, error) {
	p := o.Presentation()
	m := gormmodel.ResumeOverview{ResumeID: id.UUID(), GreetingFirst: p.Greeting.First(), GreetingSecond: p.Greeting.Second(), GreetingHighlight: p.Greeting.SecondHighlight(), GreetingThird: p.Greeting.Third()}
	if p.Availability != nil {
		v := string(p.Availability.Status())
		m.AvailabilityStatus = &v
		if d := p.Availability.From(); d != nil {
			x := d.Time()
			m.AvailabilityFrom = &x
		}
		if n := p.Availability.Note(); n != nil {
			v := *n
			m.AvailabilityNote = &v
		}
	}
	if p.ResumeLink != nil {
		v := p.ResumeLink.String()
		m.ResumeLink = &v
	}
	if p.Theme != nil {
		if p.Theme.Primary != nil {
			v := p.Theme.Primary.String()
			m.ThemePrimary = &v
		}
		if p.Theme.Secondary != nil {
			v := p.Theme.Secondary.String()
			m.ThemeSecondary = &v
		}
	}
	s := make([]gormmodel.ResumeOverviewSection, 0, len(p.Sections))
	for i, x := range p.Sections {
		s = append(s, gormmodel.ResumeOverviewSection{ResumeID: id.UUID(), SectionID: string(x.ID()), Position: i, Enabled: x.Enabled(), Collapsed: x.Collapsed()})
	}
	return m, s, nil
}
func mapOverview(m gormmodel.ResumeOverview, s []gormmodel.ResumeOverviewSection) (*domain.ResumeOverview, error) {
	g := valueobjects.NewGreeting(valueobjects.GreetingParams{First: m.GreetingFirst, Second: m.GreetingSecond, SecondHighlight: m.GreetingHighlight, Third: m.GreetingThird})
	var av *valueobjects.Availability
	if m.AvailabilityStatus != nil {
		st := enums.AvailabilityStatus(*m.AvailabilityStatus)
		v, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: st, Note: m.AvailabilityNote})
		if e != nil {
			return nil, e
		}
		if m.AvailabilityFrom != nil {
			d, e2d := valueobjects.NewDate(*m.AvailabilityFrom)
			if e2d != nil {
				return nil, e2d
			}
			v2, e2 := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: st, From: &d, Note: m.AvailabilityNote})
			if e2 != nil {
				return nil, e2
			}
			v = v2
		}
		av = &v
	}
	var link *valueobjects.AssetPath
	if m.ResumeLink != nil {
		v, e := valueobjects.NewAssetPath(*m.ResumeLink)
		if e != nil {
			return nil, e
		}
		link = &v
	}
	var theme *valueobjects.ThemeSettings
	if m.ThemePrimary != nil || m.ThemeSecondary != nil {
		var a, b *valueobjects.Color
		if m.ThemePrimary != nil {
			v, e := valueobjects.NewColor(*m.ThemePrimary)
			if e != nil {
				return nil, e
			}
			a = &v
		}
		if m.ThemeSecondary != nil {
			v, e := valueobjects.NewColor(*m.ThemeSecondary)
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
	ss := make([]valueobjects.SectionSettings, 0, len(s))
	for _, x := range s {
		v, e := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionID(x.SectionID), Enabled: x.Enabled, Collapsed: x.Collapsed})
		if e != nil {
			return nil, e
		}
		ss = append(ss, v)
	}
	return domain.NewResumeOverview(domain.ResumeOverviewParams{Presentation: domain.OverviewPresentation{Theme: theme, Greeting: g, Availability: av, ResumeLink: link, Sections: ss}})
}
