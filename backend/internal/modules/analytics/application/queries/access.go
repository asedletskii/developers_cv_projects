package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/analytics/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func authorizeAnalyticsLink(
	ctx context.Context,
	access ports.AnalyticsAccessReader,
	userID, linkID sharedvalueobjects.ID,
) error {
	if access == nil {
		return fmt.Errorf("analytics access dependencies are nil")
	}
	if userID.IsZero() || linkID.IsZero() {
		return fmt.Errorf("invalid analytics query")
	}
	resumeID, err := access.ResolveLink(ctx, linkID)
	if err != nil {
		return err
	}
	return access.EnsureResumeOwned(ctx, userID, resumeID)
}

func authorizeAnalyticsResume(
	ctx context.Context,
	access ports.AnalyticsAccessReader,
	userID, resumeID sharedvalueobjects.ID,
) error {
	if access == nil {
		return fmt.Errorf("analytics access dependencies are nil")
	}
	if userID.IsZero() || resumeID.IsZero() {
		return fmt.Errorf("invalid analytics query")
	}
	return access.EnsureResumeOwned(ctx, userID, resumeID)
}
