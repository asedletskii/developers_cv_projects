package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/application/queries"
	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/platform/database"
	"developers-cv/backend/internal/shared/domain/valueobjects"
)

type UserReader struct {
	db *gorm.DB
}

func NewUserReader(db *gorm.DB) *UserReader {
	return &UserReader{db: db}
}

func (r *UserReader) GetUser(ctx context.Context, query queries.GetUserQuery) (*domain.User, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get user: reader dependencies are nil")
	}
	if query.ID.IsZero() {
		return nil, fmt.Errorf("get user: %w", domain.ErrInvalidUser)
	}

	var model gormmodel.User
	db := database.DBFromContext(ctx, r.db)
	if err := db.WithContext(ctx).First(&model, "id = ?", query.ID.UUID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return fromModel(model)
}

func fromModel(model gormmodel.User) (*domain.User, error) {
	id, err := valueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map user id: %w", err)
	}
	return domain.NewUser(domain.UserParams{
		ID:        id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	})
}
