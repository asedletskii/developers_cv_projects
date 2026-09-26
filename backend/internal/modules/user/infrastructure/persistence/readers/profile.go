package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
)

type GetProfileReader struct{ db *gorm.DB }

func NewGetProfileReader(db *gorm.DB) *GetProfileReader { return &GetProfileReader{db: db} }

func (r *GetProfileReader) GetProfile(ctx context.Context, query queries.GetProfileQuery) (queries.Profile, error) {
	if r == nil || r.db == nil {
		return queries.Profile{}, fmt.Errorf("get profile: reader dependencies are nil")
	}
	if query.UserID.IsZero() {
		return queries.Profile{}, fmt.Errorf("get profile: invalid user id")
	}

	db := database.DBFromContext(ctx, r.db).WithContext(ctx)
	var model gormmodel.UserProfile
	if err := db.Where("user_id = ?", query.UserID.UUID()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queries.Profile{}, queries.ErrProfileNotFound
		}
		return queries.Profile{}, fmt.Errorf("get profile: %w", err)
	}
	return queries.Profile{ID: model.UserID.String(), Name: model.Name, Surname: model.Surname, Patronymic: model.Patronymic, Position: model.Position, Location: model.Location, PhotoPath: model.PhotoPath}, nil
}
