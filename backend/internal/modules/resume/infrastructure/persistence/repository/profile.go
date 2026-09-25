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

func (r *ProfileRepository) CreateProfile(ctx context.Context, c ports.CreateResumeProfileRepositoryCommand) error {
	tx, e := database.DBFromTransactionContext(ctx)
	if e != nil {
		return fmt.Errorf("create profile: %w", e)
	}
	m, fs, e := profileModels(c.Profile)
	if e != nil {
		return e
	}
	var rootCount int64
	if e = tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id = ? AND user_id = ?", m.ResumeID, c.UserID.UUID()).Count(&rootCount).Error; e != nil {
		return fmt.Errorf("create profile ownership: %w", e)
	}
	if rootCount == 0 {
		return domain.ErrEntityNotFound
	}
	if e = tx.WithContext(ctx).Create(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntity
		}
		return fmt.Errorf("create profile: %w", e)
	}
	if len(fs) > 0 {
		if e = tx.WithContext(ctx).Create(&fs).Error; e != nil {
			return fmt.Errorf("create profile work formats: %w", e)
		}
	}
	return nil
}
func (r *ProfileRepository) GetProfile(ctx context.Context, q ports.GetResumeProfileRepositoryCommand) (*domain.Profile, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get profile: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.ResumeProfile
	if e := tx.WithContext(ctx).Where("resume_id = ? AND resume_id IN (SELECT id FROM resumes WHERE id = ? AND user_id = ?)", q.ResumeID.UUID(), q.ResumeID.UUID(), q.UserID.UUID()).First(&m).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get profile: %w", e)
	}
	var fs []gormmodel.ResumeProfileWorkFormat
	if e := tx.WithContext(ctx).Where("profile_id = ?", m.ID).Order("position").Find(&fs).Error; e != nil {
		return nil, fmt.Errorf("get profile work formats: %w", e)
	}
	return mapProfile(m, fs)
}
func (r *ProfileRepository) UpdateProfile(ctx context.Context, c ports.UpdateResumeProfileRepositoryCommand) error {
	tx, e := database.DBFromTransactionContext(ctx)
	if e != nil {
		return fmt.Errorf("update profile: %w", e)
	}
	var current gormmodel.ResumeProfile
	if e = tx.WithContext(ctx).Where("resume_id = ? AND resume_id IN (SELECT id FROM resumes WHERE id = ? AND user_id = ?)", c.Profile.ResumeID().UUID(), c.Profile.ResumeID().UUID(), c.UserID.UUID()).First(&current).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return domain.ErrEntityNotFound
		}
		return fmt.Errorf("update profile lookup: %w", e)
	}
	m, fs, e := profileModels(c.Profile)
	if e != nil {
		return e
	}
	m.ID = current.ID
	m.ResumeID = current.ResumeID
	result := tx.WithContext(ctx).Model(&gormmodel.ResumeProfile{}).Where("id = ? AND resume_id IN (SELECT id FROM resumes WHERE id = ? AND user_id = ?)", m.ID, m.ResumeID, c.UserID.UUID()).Updates(map[string]any{"name": m.Name, "surname": m.Surname, "patronymic": m.Patronymic, "position": m.Position, "photo_path": m.PhotoPath, "date_of_birth": m.DateOfBirth, "location": m.Location, "relocation": m.Relocation})
	if result.Error != nil {
		return fmt.Errorf("update profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEntityNotFound
	}
	if e = tx.WithContext(ctx).Where("profile_id = ?", m.ID).Delete(&gormmodel.ResumeProfileWorkFormat{}).Error; e != nil {
		return fmt.Errorf("replace profile work formats: %w", e)
	}
	if len(fs) > 0 {
		if e = tx.WithContext(ctx).Create(&fs).Error; e != nil {
			return fmt.Errorf("create profile work formats: %w", e)
		}
	}
	return nil
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
	photo, e := valueobjects.NewAssetPath(m.PhotoPath)
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
