package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *ResumeRepository) Get(ctx context.Context, params ports.GetResumeParams) (*domain.Resume, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume: repository dependencies are nil")
	}
	var resume *domain.Resume
	err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		var err error
		resume, err = r.loadResume(snapshotCtx, params.ResumeID, params.UserID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return resume, nil
}

func (r *ResumeRepository) Create(ctx context.Context, resume *domain.Resume) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("create resume: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume: %w", err)
	}
	if resume == nil {
		return fmt.Errorf("create resume: resume is nil")
	}

	now := time.Now().UTC()
	model := mappers.ResumeToDB(resume, now, now)
	tx = tx.WithContext(ctx)
	if err := tx.Create(&model).Error; err != nil {
		return fmt.Errorf("create resume: %w", mapDBError(err))
	}
	if err := persistAllSections(tx, resume); err != nil {
		return fmt.Errorf("create resume sections: %w", err)
	}
	return nil
}

func (r *ResumeRepository) Save(ctx context.Context, resume *domain.Resume) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("save resume: repository dependencies are nil")
	}
	if resume == nil {
		return fmt.Errorf("save resume: resume is nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("save resume: %w", err)
	}
	result := tx.WithContext(ctx).
		Model(&gormmodel.Resume{}).
		Where("id = ? AND user_id = ? AND version = ?", resume.ID().UUID(), resume.UserID().UUID(), resume.Version()).
		Updates(map[string]any{
			"title":      resume.Title(),
			"language":   string(resume.Language()),
			"status":     string(resume.Status()),
			"version":    gorm.Expr("version + 1"),
			"updated_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("save resume: %w", mapDBError(result.Error))
	}
	if result.RowsAffected != 1 {
		return domain.ErrConcurrencyConflict
	}
	if err := persistAllSections(tx.WithContext(ctx), resume); err != nil {
		return fmt.Errorf("save resume sections: %w", err)
	}
	return nil
}

type resumeSectionPersister struct {
	section ports.ResumeMutationSection
	name    string
	save    func(*gorm.DB, *domain.Resume) error
}

var resumeSectionPersisters = [...]resumeSectionPersister{
	{section: ports.ResumeMutationOverview, name: "overview", save: persistOverview},
	{section: ports.ResumeMutationProfile, name: "profile", save: persistProfile},
	{section: ports.ResumeMutationContacts, name: "contacts", save: persistContacts},
	{section: ports.ResumeMutationSkills, name: "skills", save: persistSkills},
	{section: ports.ResumeMutationLanguages, name: "spoken languages", save: persistLanguages},
	{section: ports.ResumeMutationTechnologies, name: "technologies", save: persistTechnologies},
	{section: ports.ResumeMutationExperiences, name: "experiences", save: persistExperiences},
	{section: ports.ResumeMutationEducations, name: "educations", save: persistEducations},
	{section: ports.ResumeMutationCourses, name: "courses", save: persistCourses},
	{section: ports.ResumeMutationActivities, name: "activities", save: persistActivities},
	{section: ports.ResumeMutationProjects, name: "projects", save: persistProjects},
}

func persistAllSections(tx *gorm.DB, resume *domain.Resume) error {
	for _, section := range resumeSectionPersisters {
		if err := section.save(tx, resume); err != nil {
			return fmt.Errorf("persist %s: %w", section.name, err)
		}
	}
	return nil
}

func (r *ResumeRepository) SaveSection(ctx context.Context, resume *domain.Resume, section ports.ResumeMutationSection) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("save resume section: repository dependencies are nil")
	}
	if resume == nil {
		return fmt.Errorf("save resume section: resume is nil")
	}
	saveSection, supported := mutationSectionPersister(section)
	if !supported {
		return fmt.Errorf("save resume section: unsupported mutation section %q", section)
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("save resume section: %w", err)
	}

	updates := map[string]any{
		"version":    gorm.Expr("version + 1"),
		"updated_at": time.Now().UTC(),
	}
	switch section {
	case ports.ResumeMutationTitle:
		updates["title"] = resume.Title()
	case ports.ResumeMutationStatus:
		updates["status"] = string(resume.Status())
	}
	result := tx.WithContext(ctx).Model(&gormmodel.Resume{}).
		Where("id = ? AND user_id = ? AND version = ?", resume.ID().UUID(), resume.UserID().UUID(), resume.Version()).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("save resume section: %w", mapDBError(result.Error))
	}
	if result.RowsAffected != 1 {
		return domain.ErrConcurrencyConflict
	}

	if saveSection != nil {
		if err := saveSection(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("save resume section %q: %w", section, err)
		}
	}
	return nil
}

func mutationSectionPersister(section ports.ResumeMutationSection) (func(*gorm.DB, *domain.Resume) error, bool) {
	switch section {
	case ports.ResumeMutationTitle, ports.ResumeMutationStatus:
		return nil, true
	}
	for _, persister := range resumeSectionPersisters {
		if persister.section == section {
			return persister.save, true
		}
	}
	return nil, false
}

func (r *ResumeRepository) LockForDelete(ctx context.Context, params ports.DeleteResumeParams) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("lock resume for delete: repository dependencies are nil")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("lock resume for delete: %w", err)
	}
	var model gormmodel.Resume
	if err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id", "version").
		Where("id = ? AND user_id = ?", params.ResumeID.UUID(), params.UserID.UUID()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrConcurrencyConflict
		}
		return fmt.Errorf("lock resume for delete: %w", mapDBError(err))
	}
	if model.Version != params.Version {
		return domain.ErrConcurrencyConflict
	}
	return nil
}

func (r *ResumeRepository) Delete(ctx context.Context, params ports.DeleteResumeParams) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("delete resume: repository dependencies are nil")
	}
	resumeID, userID, version := params.ResumeID, params.UserID, params.Version
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume: %w", err)
	}
	res := tx.WithContext(ctx).Where("id = ? AND user_id = ? AND version = ?", resumeID.UUID(), userID.UUID(), version).Delete(&gormmodel.Resume{})
	if res.Error != nil {
		return fmt.Errorf("delete resume: %w", mapDBError(res.Error))
	}
	if res.RowsAffected != 1 {
		return domain.ErrConcurrencyConflict
	}
	return nil
}

func (r *ResumeRepository) loadResume(ctx context.Context, resumeID sharedvalueobjects.ID, userID sharedvalueobjects.ID) (*domain.Resume, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.Resume
	if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", resumeID.UUID(), userID.UUID()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume: %w", mapDBError(err))
	}
	profile, err := loadResumeProfile(ctx, tx, m.ID)
	if err != nil {
		return nil, err
	}
	contacts, err := loadResumeContacts(ctx, tx, m.ID)
	if err != nil {
		return nil, err
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
	overview, err := loadResumeOverview(ctx, tx, m.ID)
	if err != nil {
		return nil, err
	}

	snapshot, err := mappers.ResumeToDomainSnapshot(
		m,
		overview,
		profile,
		contacts,
		skills,
		languages,
		technologies,
		experiences,
		educations,
		courses,
		activities,
		projects,
	)
	if err != nil {
		return nil, fmt.Errorf("map resume: %w", err)
	}
	return domain.RestoreResume(snapshot), nil
}
