package readers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
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

func (r *ResumeReader) GetOwnedOverview(ctx context.Context, query queries.GetResumeQuery) (queries.ResumeOverviewView, error) {
	if r == nil {
		return queries.ResumeOverviewView{}, fmt.Errorf("get owned resume overview: reader is nil")
	}
	var result queries.ResumeOverviewView
	err := database.WithReadSnapshot(ctx, r.db, func(snapshotCtx context.Context) error {
		db, err := resolveDB(snapshotCtx, r.db, "get owned resume overview")
		if err != nil {
			return err
		}
		userID := query.UserID.UUID()
		if _, err := r.loadMetadata(snapshotCtx, db, query.ResumeID.UUID(), &userID); err != nil {
			return err
		}
		view := queries.ResumeView{}
		if err := r.loadOverview(snapshotCtx, db, query.ResumeID.UUID(), &view); err != nil {
			return err
		}
		result = view.Overview
		return nil
	})
	if err != nil {
		return queries.ResumeOverviewView{}, err
	}
	return result, nil
}

func (r *ResumeReader) OwnsResume(ctx context.Context, query queries.GetResumeQuery) (bool, error) {
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

func (r *ResumeReader) GetPublic(ctx context.Context, query queries.PublicResumeQuery) (queries.PublicResumeView, error) {
	if r == nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume: reader is nil")
	}
	var result queries.PublicResumeView
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
		result = queries.ToPublicResumeView(view)
		return nil
	})
	if err != nil {
		return queries.PublicResumeView{}, err
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
			return queries.ResumeMetadataView{}, queries.ErrResumeNotFound
		}
		return queries.ResumeMetadataView{}, fmt.Errorf("get resume metadata: %w", err)
	}
	if userID == nil && model.Status != string(domain.ResumeStatusOpened) {
		return queries.ResumeMetadataView{}, queries.ErrResumeNotFound
	}
	return queries.ResumeMetadataView{
		ID: model.ID.String(), Title: model.Title, Language: model.Language, Status: model.Status,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt, Version: model.Version,
	}, nil
}
