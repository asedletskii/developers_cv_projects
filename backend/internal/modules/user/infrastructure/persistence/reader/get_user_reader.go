package reader

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/user/application/queries"
)

type GetUserReader struct {
	pool *pgxpool.Pool
}

func NewGetUserReader(pool *pgxpool.Pool) *GetUserReader {
	return &GetUserReader{pool: pool}
}

func (r *GetUserReader) GetUser(ctx context.Context, query queries.GetUserQuery) (queries.User, error) {
	if r == nil || r.pool == nil {
		return queries.User{}, fmt.Errorf("get user: reader dependencies are nil")
	}
	if query.UserID.IsZero() {
		return queries.User{}, fmt.Errorf("get user: invalid user id")
	}

	var user queries.User
	var patronymic, position, location, photoPath *string
	err := r.pool.QueryRow(ctx, getUserSQL, query.UserID.UUID()).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&patronymic,
		&position,
		&location,
		&photoPath,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.User{}, queries.ErrUserNotFound
		}
		return queries.User{}, fmt.Errorf("get user: %w", err)
	}

	user.Patronymic = patronymic
	user.Position = position
	user.Location = location
	user.PhotoPath = photoPath
	return user, nil
}

const getUserSQL = `
SELECT
    u.id::text,
    p.name,
    p.surname,
    p.patronymic,
    p.position,
    p.location,
    p.photo_path
FROM users u
JOIN user_profiles p ON p.user_id = u.id
WHERE u.id = $1
`
