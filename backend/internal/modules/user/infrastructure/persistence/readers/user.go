package readers

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/user/application/queries"
	usercontracts "developers-cv/backend/internal/modules/user/contracts"
	"developers-cv/backend/internal/modules/user/domain"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/user/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
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
	var model gormmodel.User
	db := database.DBFromContext(ctx, r.db)
	if err := db.WithContext(ctx).First(&model, "id = ?", query.ID.UUID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, queries.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	user, err := mappers.UserToDomain(model)
	if err != nil {
		return nil, fmt.Errorf("get user: map user: %w", err)
	}
	return user, nil
}

func (r *UserReader) EnsureUserExists(ctx context.Context, userID sharedvalueobjects.ID) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("ensure user exists: reader dependencies are nil")
	}
	if userID.IsZero() {
		return usercontracts.ErrUserNotFound
	}
	if _, err := r.GetUser(ctx, queries.GetUserQuery{ID: userID}); err != nil {
		if errors.Is(err, queries.ErrUserNotFound) {
			return usercontracts.ErrUserNotFound
		}
		return fmt.Errorf("ensure user exists: %w", err)
	}
	return nil
}
