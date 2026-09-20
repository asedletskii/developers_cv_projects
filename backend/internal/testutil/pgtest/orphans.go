package pgtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
)

type OrphanCounts struct {
	AnalyticsWithoutLink      int64
	LinksWithoutResume        int64
	ResumesWithoutUser        int64
	AuthIdentitiesWithoutUser int64
	AuthSessionsWithoutUser   int64
}

func (c OrphanCounts) Total() int64 {
	return c.AnalyticsWithoutLink + c.LinksWithoutResume + c.ResumesWithoutUser +
		c.AuthIdentitiesWithoutUser + c.AuthSessionsWithoutUser
}

func CountOrphans(ctx context.Context, conn *pgx.Conn) (OrphanCounts, error) {
	if conn == nil {
		return OrphanCounts{}, fmt.Errorf("count orphan rows: connection is nil")
	}
	var counts OrphanCounts
	err := conn.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM analytics.analytics_views v
			 LEFT JOIN resume_link.resume_links l ON l.id = v.link_id
			 WHERE l.id IS NULL),
			(SELECT count(*) FROM resume_link.resume_links l
			 LEFT JOIN resume.resumes r ON r.id = l.resume_id
			 WHERE r.id IS NULL),
			(SELECT count(*) FROM resume.resumes r
			 LEFT JOIN users.users u ON u.id = r.user_id
			 WHERE u.id IS NULL),
			(SELECT count(*) FROM auth.user_identities i
			 LEFT JOIN users.users u ON u.id = i.user_id
			 WHERE u.id IS NULL),
			(SELECT count(*) FROM auth.auth_sessions s
			 LEFT JOIN users.users u ON u.id = s.user_id
			 WHERE u.id IS NULL)
	`).Scan(
		&counts.AnalyticsWithoutLink,
		&counts.LinksWithoutResume,
		&counts.ResumesWithoutUser,
		&counts.AuthIdentitiesWithoutUser,
		&counts.AuthSessionsWithoutUser,
	)
	if err != nil {
		return OrphanCounts{}, fmt.Errorf("count orphan rows: %w", err)
	}
	return counts, nil
}

func RequireNoOrphans(t testing.TB, ctx context.Context, conn *pgx.Conn) {
	t.Helper()
	counts, err := CountOrphans(ctx, conn)
	if err != nil {
		t.Fatalf("check orphan rows: %v", err)
	}
	if counts.Total() != 0 {
		t.Fatalf("found orphan rows: %+v", counts)
	}
}
