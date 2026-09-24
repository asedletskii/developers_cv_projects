package reader

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/user/application/queries"
)

type GetProfileReader struct {
	pool *pgxpool.Pool
}

func NewGetProfileReader(pool *pgxpool.Pool) *GetProfileReader {
	return &GetProfileReader{pool: pool}
}

func (r *GetProfileReader) GetProfile(ctx context.Context, query queries.GetProfileQuery) (queries.Profile, error) {
	if r == nil || r.pool == nil {
		return queries.Profile{}, fmt.Errorf("get profile: reader dependencies are nil")
	}
	if query.UserID.IsZero() {
		return queries.Profile{}, fmt.Errorf("get profile: invalid user id")
	}

	var profile queries.Profile
	var patronymic, position, location, photoPath *string
	err := r.pool.QueryRow(ctx, getProfileSQL, query.UserID.UUID()).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Surname,
		&patronymic,
		&position,
		&location,
		&photoPath,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.Profile{}, queries.ErrProfileNotFound
		}
		return queries.Profile{}, fmt.Errorf("get profile: %w", err)
	}

	profile.Patronymic = patronymic
	profile.Position = position
	profile.Location = location
	profile.PhotoPath = photoPath
	return profile, nil
}

const getProfileSQL = `
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
