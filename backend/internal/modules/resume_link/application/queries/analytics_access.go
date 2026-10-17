package queries

import (
	"context"
	"fmt"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func resolveAnalyticsLink(
	ctx context.Context,
	links ResumeLinkPublicReader,
	userID, linkID sharedvalueobjects.ID,
) (ResumeLinkPublicView, sharedvalueobjects.ID, error) {
	if links == nil {
		return ResumeLinkPublicView{}, sharedvalueobjects.ID{}, fmt.Errorf("analytics link reader is nil")
	}
	if userID.IsZero() || linkID.IsZero() {
		return ResumeLinkPublicView{}, sharedvalueobjects.ID{}, fmt.Errorf("invalid analytics query")
	}
	link, err := links.GetByID(ctx, PublicResumeLinkQuery{LinkID: linkID})
	if err != nil {
		return ResumeLinkPublicView{}, sharedvalueobjects.ID{}, err
	}
	resumeID, err := sharedvalueobjects.ParseID(link.ResumeID)
	if err != nil {
		return ResumeLinkPublicView{}, sharedvalueobjects.ID{}, fmt.Errorf("resolve resume link: %w", err)
	}
	return link, resumeID, nil
}

func authorizeAnalyticsLink(
	ctx context.Context,
	links ResumeLinkPublicReader,
	resumes ports.ResumeOwnershipReader,
	userID, linkID sharedvalueobjects.ID,
) (ResumeLinkPublicView, sharedvalueobjects.ID, error) {
	if resumes == nil {
		return ResumeLinkPublicView{}, sharedvalueobjects.ID{}, fmt.Errorf("analytics ownership reader is nil")
	}
	link, resumeID, err := resolveAnalyticsLink(ctx, links, userID, linkID)
	if err != nil {
		return ResumeLinkPublicView{}, sharedvalueobjects.ID{}, err
	}
	if err := EnsureResumeOwned(ctx, resumes, resumequeries.GetResumeQuery{ResumeID: resumeID, UserID: userID}); err != nil {
		return ResumeLinkPublicView{}, sharedvalueobjects.ID{}, err
	}
	return link, resumeID, nil
}

func authorizeAnalyticsResume(
	ctx context.Context,
	resumes ports.ResumeOwnershipReader,
	userID, resumeID sharedvalueobjects.ID,
) error {
	if resumes == nil {
		return fmt.Errorf("analytics access dependencies are nil")
	}
	if userID.IsZero() || resumeID.IsZero() {
		return fmt.Errorf("invalid analytics query")
	}
	err := EnsureResumeOwned(ctx, resumes, resumequeries.GetResumeQuery{ResumeID: resumeID, UserID: userID})
	return err
}
