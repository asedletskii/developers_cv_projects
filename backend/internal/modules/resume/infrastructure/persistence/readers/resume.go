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
	if query.ResumeID.IsZero() || query.UserID.IsZero() {
		return queries.ResumeView{}, fmt.Errorf("get owned resume: invalid identifiers")
	}
	db, err := resolveDB(ctx, r.db, "get owned resume")
	if err != nil {
		return queries.ResumeView{}, err
	}

	userID := query.UserID.UUID()
	metadata, err := r.loadMetadata(ctx, db, query.ResumeID.UUID(), &userID)
	if err != nil {
		return queries.ResumeView{}, err
	}
	return r.loadView(ctx, db, query.ResumeID.UUID(), metadata)
}

func (r *ResumeReader) GetPublic(ctx context.Context, query queries.PublicResumeQuery) (queries.ResumeView, error) {
	if r == nil {
		return queries.ResumeView{}, fmt.Errorf("get public resume: reader is nil")
	}
	if query.ID.IsZero() {
		return queries.ResumeView{}, fmt.Errorf("get public resume: invalid resume id")
	}
	db, err := resolveDB(ctx, r.db, "get public resume")
	if err != nil {
		return queries.ResumeView{}, err
	}

	metadata, err := r.loadMetadata(ctx, db, query.ID.UUID(), nil)
	if err != nil {
		return queries.ResumeView{}, err
	}
	return r.loadView(ctx, db, query.ID.UUID(), metadata)
}

func (r *ResumeReader) loadView(ctx context.Context, db *gorm.DB, resumeID uuid.UUID, metadata queries.ResumeMetadataView) (queries.ResumeView, error) {
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
	if err := r.loadContacts(ctx, db, resumeID, &view); err != nil {
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
