package queries

import (
	"context"
	"errors"
	"fmt"

	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	resumelinkdomain "developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type AnalyticsResumeOverviewReader interface {
	GetOwnedOverview(context.Context, resumequeries.GetResumeQuery) (resumequeries.ResumeOverviewView, error)
}

type analyticsAccessReader struct {
	links        ResumeLinkPublicReader
	resumes      ports.ResumeOwnershipReader
	overview     AnalyticsResumeOverviewReader
	publicResume ResumePublicQuery
}

func NewAnalyticsAccessReader(
	links ResumeLinkPublicReader,
	resumes ports.ResumeOwnershipReader,
	overview AnalyticsResumeOverviewReader,
	publicResume ResumePublicQuery,
) analyticsports.AnalyticsAccessReader {
	return &analyticsAccessReader{links: links, resumes: resumes, overview: overview, publicResume: publicResume}
}

func (r *analyticsAccessReader) ResolveLink(
	ctx context.Context,
	linkID sharedvalueobjects.ID,
) (sharedvalueobjects.ID, error) {
	if r == nil || r.links == nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("analytics link reader is nil")
	}
	if linkID.IsZero() {
		return sharedvalueobjects.ID{}, fmt.Errorf("invalid analytics link id")
	}
	link, err := r.links.GetByID(ctx, PublicResumeLinkQuery{LinkID: linkID})
	if err != nil {
		return sharedvalueobjects.ID{}, normalizeAnalyticsAccessError(err)
	}
	resumeID, err := sharedvalueobjects.ParseID(link.ResumeID)
	if err != nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("resolve analytics link: %w", err)
	}
	return resumeID, nil
}

func (r *analyticsAccessReader) EnsureResumeOwned(
	ctx context.Context,
	userID, resumeID sharedvalueobjects.ID,
) error {
	if r == nil || r.resumes == nil {
		return fmt.Errorf("analytics ownership reader is nil")
	}
	return normalizeAnalyticsAccessError(EnsureResumeOwned(ctx, r.resumes, resumequeries.GetResumeQuery{ResumeID: resumeID, UserID: userID}))
}

func (r *analyticsAccessReader) ResumeSectionOrder(
	ctx context.Context,
	userID, resumeID sharedvalueobjects.ID,
) ([]string, error) {
	if r == nil || r.overview == nil {
		return nil, fmt.Errorf("analytics resume overview reader is nil")
	}
	overview, err := r.overview.GetOwnedOverview(ctx, resumequeries.GetResumeQuery{ResumeID: resumeID, UserID: userID})
	if err != nil {
		return nil, normalizeAnalyticsAccessError(err)
	}
	order := []string{"about"}
	seen := map[string]bool{"about": true}
	for _, section := range overview.Sections {
		if section.Enabled && section.ID != "techStack" && !seen[section.ID] {
			order = append(order, section.ID)
			seen[section.ID] = true
		}
	}
	return order, nil
}

func (r *analyticsAccessReader) ResumeAnalyticsLayout(
	ctx context.Context,
	userID, resumeID sharedvalueobjects.ID,
) (analyticsports.ResumeAnalyticsLayout, error) {
	if r == nil || r.publicResume == nil {
		return analyticsports.ResumeAnalyticsLayout{}, fmt.Errorf("analytics resume reader is nil")
	}
	if err := r.EnsureResumeOwned(ctx, userID, resumeID); err != nil {
		return analyticsports.ResumeAnalyticsLayout{}, err
	}
	view, err := r.publicResume.Handle(ctx, resumequeries.PublicResumeQuery{ID: resumeID})
	if err != nil {
		return analyticsports.ResumeAnalyticsLayout{}, normalizeAnalyticsAccessError(err)
	}
	layout, err := resumequeries.ToResumeAnalyticsLayout(view)
	if err != nil {
		return analyticsports.ResumeAnalyticsLayout{}, fmt.Errorf("build analytics resume layout: %w", err)
	}
	result := analyticsports.ResumeAnalyticsLayout{
		Sections:     make([]analyticsports.ResumeAnalyticsSection, 0, len(layout.Sections)),
		Technologies: append([]string{}, layout.Technologies...),
	}
	for _, section := range layout.Sections {
		result.Sections = append(result.Sections, analyticsports.ResumeAnalyticsSection{ID: section.ID, Hash: section.Hash})
	}
	return result, nil
}

func normalizeAnalyticsAccessError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrResumeLinkNotFound) ||
		errors.Is(err, resumelinkdomain.ErrEntityNotFound) ||
		errors.Is(err, resumequeries.ErrResumeNotFound) {
		return fmt.Errorf("%w: %w", analyticsdomain.ErrAnalyticsSubjectNotFound, err)
	}
	return err
}
