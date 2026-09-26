package readers

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
)

type ProfileReader struct{ db *gorm.DB }

func NewProfileReader(db *gorm.DB) *ProfileReader { return &ProfileReader{db: db} }

func (r *ProfileReader) GetResumeProfile(ctx context.Context, q queries.GetResumeProfileQuery) (queries.ResumeProfileView, error) {
	db, err := resolveDB(ctx, r.db, "get resume profile")
	if err != nil {
		return queries.ResumeProfileView{}, err
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return queries.ResumeProfileView{}, fmt.Errorf("get resume profile: invalid query")
	}

	var model gormmodel.ResumeProfile
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_profiles.resume_id").Where("resume_profiles.resume_id = ?", q.ResumeID.UUID())
	readDB = readDB.Where("resumes.user_id = ?", q.UserID.UUID())
	err = readDB.First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ResumeProfileView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeProfileView{}, fmt.Errorf("get resume profile: %w", err)
	}

	var formats []gormmodel.ResumeProfileWorkFormat
	if err := db.Where("profile_id = ?", model.ID).Order("position").Find(&formats).Error; err != nil {
		return queries.ResumeProfileView{}, fmt.Errorf("get profile work formats: %w", err)
	}

	workFormat := make([]string, 0, len(formats))
	for _, format := range formats {
		workFormat = append(workFormat, format.Value)
	}

	return queries.ResumeProfileView{
		ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Surname: model.Surname,
		Patronymic: model.Patronymic, Position: model.PositionValue(), PhotoPath: model.PhotoPathValue(),
		DateOfBirth: model.DateOfBirth.Format("2006-01-02"), Location: model.LocationValue(), Relocation: model.Relocation,
		WorkFormat: workFormat,
	}, nil
}

func (r *ProfileReader) ReadResumeProfile(ctx context.Context, resumeID valueobjects.ID) (queries.ResumeProfileView, error) {
	db, err := resolveDB(ctx, r.db, "get resume profile")
	if err != nil {
		return queries.ResumeProfileView{}, err
	}
	if resumeID.IsZero() {
		return queries.ResumeProfileView{}, fmt.Errorf("get resume profile: invalid query")
	}

	var model gormmodel.ResumeProfile
	readDB := db.Joins("JOIN resumes ON resumes.id = resume_profiles.resume_id").Where("resume_profiles.resume_id = ?", resumeID.UUID())
	err = readDB.First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queries.ResumeProfileView{}, queries.ErrResumeNotFound
	}
	if err != nil {
		return queries.ResumeProfileView{}, fmt.Errorf("get resume profile: %w", err)
	}

	var formats []gormmodel.ResumeProfileWorkFormat
	if err := db.Where("profile_id = ?", model.ID).Order("position").Find(&formats).Error; err != nil {
		return queries.ResumeProfileView{}, fmt.Errorf("get profile work formats: %w", err)
	}

	workFormat := make([]string, 0, len(formats))
	for _, format := range formats {
		workFormat = append(workFormat, format.Value)
	}

	return queries.ResumeProfileView{
		ID: model.ID.String(), ResumeID: model.ResumeID.String(), Name: model.Name, Surname: model.Surname,
		Patronymic: model.Patronymic, Position: model.PositionValue(), PhotoPath: model.PhotoPathValue(),
		DateOfBirth: model.DateOfBirth.Format("2006-01-02"), Location: model.LocationValue(), Relocation: model.Relocation,
		WorkFormat: workFormat,
	}, nil
}
