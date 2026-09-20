package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"developers-cv/backend/internal/modules/resume/application/queries"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
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
	var result queries.ResumeView
	err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		db, err := resolveDB(snapshotCtx, r.db, "get owned resume")
		if err != nil {
			return err
		}
		userID := query.UserID.UUID()
		metadata, err := r.loadMetadata(snapshotCtx, db, query.ResumeID.UUID(), &userID)
		if err != nil {
			return err
		}
		result, err = r.loadView(snapshotCtx, db, query.ResumeID.UUID(), metadata, false)
		return err
	})
	if err != nil {
		return queries.ResumeView{}, err
	}
	return result, nil
}

func (r *ResumeReader) GetOwnedSections(ctx context.Context, query resumecontracts.OwnedResumeSectionsQuery) ([]queries.ResumeSectionSettings, error) {
	if r == nil {
		return nil, fmt.Errorf("get owned resume sections: reader is nil")
	}
	var result []queries.ResumeSectionSettings
	err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		db, err := resolveDB(snapshotCtx, r.db, "get owned resume sections")
		if err != nil {
			return err
		}
		userID := query.UserID.UUID()
		if _, err := r.loadMetadata(snapshotCtx, db, query.ResumeID.UUID(), &userID); err != nil {
			return err
		}
		var sections []gormmodel.ResumeOverviewSection
		if err := db.WithContext(snapshotCtx).Where("resume_id = ?", query.ResumeID.UUID()).Order("position").Find(&sections).Error; err != nil {
			return fmt.Errorf("get owned resume sections: %w", err)
		}
		result = make([]queries.ResumeSectionSettings, 0, len(sections))
		for _, section := range sections {
			result = append(result, queries.ResumeSectionSettings{ID: section.SectionID, Enabled: section.Enabled, Collapsed: section.Collapsed, Style: section.Style})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ResumeReader) LockOwnedResumeForUpdate(ctx context.Context, query resumecontracts.OwnedResumeLock) (bool, error) {
	if r == nil {
		return false, fmt.Errorf("lock owned resume for update: reader is nil")
	}
	if query.ResumeID.IsZero() || query.UserID.IsZero() {
		return false, fmt.Errorf("lock owned resume for update: invalid identifiers")
	}
	db, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return false, fmt.Errorf("lock owned resume for update: %w", err)
	}
	var model gormmodel.Resume
	err = db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").
		Where("id = ? AND user_id = ?", query.ResumeID.UUID(), query.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("lock owned resume for update: %w", err)
	}
	return true, nil
}

func (r *ResumeReader) OwnsResume(ctx context.Context, query resumecontracts.ResumeOwnershipQuery) (bool, error) {
	if r == nil {
		return false, fmt.Errorf("check resume ownership: reader is nil")
	}
	db, err := resolveDB(ctx, r.db, "check resume ownership")
	if err != nil {
		return false, err
	}
	var model gormmodel.Resume
	err = db.Select("id").Where("id = ? AND user_id = ?", query.ResumeID.UUID(), query.UserID.UUID()).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check resume ownership: %w", err)
	}
	return true, nil
}

func (r *ResumeReader) GetPublicationState(
	ctx context.Context,
	query resumecontracts.ResumePublicationQuery,
) (resumecontracts.ResumePublicationState, error) {
	if r == nil {
		return resumecontracts.ResumePublicationState{}, fmt.Errorf("get resume publication state: reader is nil")
	}
	if query.ResumeID.IsZero() {
		return resumecontracts.ResumePublicationState{}, fmt.Errorf("get resume publication state: invalid resume id")
	}
	db, err := resolveDB(ctx, r.db, "get resume publication state")
	if err != nil {
		return resumecontracts.ResumePublicationState{}, err
	}
	return r.readPublicationState(ctx, db, query.ResumeID)
}

func (r *ResumeReader) LockPublicationState(
	ctx context.Context,
	query resumecontracts.ResumePublicationQuery,
) (resumecontracts.ResumePublicationState, error) {
	if r == nil {
		return resumecontracts.ResumePublicationState{}, fmt.Errorf("lock resume publication state: reader is nil")
	}
	if query.ResumeID.IsZero() {
		return resumecontracts.ResumePublicationState{}, fmt.Errorf("lock resume publication state: invalid resume id")
	}
	db, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return resumecontracts.ResumePublicationState{}, fmt.Errorf("lock resume publication state: %w", err)
	}
	return r.readPublicationState(ctx, db.WithContext(ctx).Clauses(clause.Locking{Strength: "SHARE"}), query.ResumeID)
}

func (r *ResumeReader) readPublicationState(
	ctx context.Context,
	db *gorm.DB,
	resumeID sharedvalueobjects.ID,
) (resumecontracts.ResumePublicationState, error) {
	var model gormmodel.Resume
	if err := db.WithContext(ctx).Select("id, status").Where("id = ?", resumeID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resumecontracts.ResumePublicationState{}, resumecontracts.ErrResumeNotFound
		}
		return resumecontracts.ResumePublicationState{}, fmt.Errorf("read resume publication state: %w", err)
	}
	return resumecontracts.ResumePublicationState{
		Exists:    true,
		Published: model.Status == string(domain.ResumeStatusOpened),
		Status:    model.Status,
	}, nil
}

func (r *ResumeReader) GetPublic(ctx context.Context, query resumecontracts.PublicResumeQuery) (queries.ResumeView, error) {
	if r == nil {
		return queries.ResumeView{}, fmt.Errorf("get public resume: reader is nil")
	}
	var result queries.ResumeView
	err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		db, err := resolveDB(snapshotCtx, r.db, "get public resume")
		if err != nil {
			return err
		}
		metadata, err := r.loadMetadata(snapshotCtx, db, query.ID.UUID(), nil)
		if err != nil {
			return err
		}
		view, err := r.loadView(snapshotCtx, db, query.ID.UUID(), metadata, true)
		if err != nil {
			return err
		}
		result = view
		return nil
	})
	if err != nil {
		return queries.ResumeView{}, err
	}
	return result, nil
}

func (r *ResumeReader) loadView(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, metadata queries.ResumeMetadataView, public bool) (queries.ResumeView, error) {
	view := queries.ResumeView{
		ID: metadata.ID, Title: metadata.Title, Language: metadata.Language, Status: metadata.Status,
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
	if public {
		if err := r.loadPublicContacts(ctx, db, resumeID, &view); err != nil {
			return queries.ResumeView{}, err
		}
	} else if err := r.loadContacts(ctx, db, resumeID, &view); err != nil {
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
			return queries.ResumeMetadataView{}, resumecontracts.ErrResumeNotFound
		}
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: %w", err)
	}
	if userID == nil && model.Status != string(domain.ResumeStatusOpened) {
		return queries.ResumeMetadataView{}, resumecontracts.ErrResumeNotFound
	}
	return queries.ResumeMetadataView{
		ID: model.ID.String(), Title: model.Title, Language: model.Language, Status: model.Status,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt, Version: model.Version,
	}, nil
}
