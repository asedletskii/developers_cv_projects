package analyticsaccess

import (
	"context"
	"errors"
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	"developers-cv/backend/internal/modules/resume_link/application/ownership"
	resumelinkqueries "developers-cv/backend/internal/modules/resume_link/application/queries"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type LinkLifecycleReader interface {
	GetByID(context.Context, resumelinkqueries.PublicResumeLinkQuery) (resumelinkqueries.ResumeLinkPublicView, error)
	LockByID(context.Context, resumelinkqueries.PublicResumeLinkQuery) (resumelinkqueries.ResumeLinkPublicView, error)
	ListByResumeID(context.Context, resumelinkqueries.ListResumeLinksQuery) ([]resumelinkqueries.ResumeLinkListItem, error)
}

type Reader struct {
	links        LinkLifecycleReader
	resumes      resumecontracts.ResumeOwnershipChecker
	publication  resumecontracts.ResumePublicationLocker
	overview     resumecontracts.OwnedResumeSectionsQueryHandler
	publicResume resumecontracts.PublicResumeQueryHandler
}

func NewReader(
	links LinkLifecycleReader,
	resumes resumecontracts.ResumeOwnershipChecker,
	publication resumecontracts.ResumePublicationLocker,
	overview resumecontracts.OwnedResumeSectionsQueryHandler,
	publicResume resumecontracts.PublicResumeQueryHandler,
) *Reader {
	return &Reader{
		links: links, resumes: resumes, publication: publication,
		overview: overview, publicResume: publicResume,
	}
}

func (r *Reader) ResolveLink(
	ctx context.Context,
	linkID sharedvalueobjects.ID,
) (sharedvalueobjects.ID, error) {
	if r == nil || r.links == nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("analytics link reader is nil")
	}
	if linkID.IsZero() {
		return sharedvalueobjects.ID{}, fmt.Errorf("invalid analytics link id")
	}
	link, err := r.links.GetByID(ctx, resumelinkqueries.PublicResumeLinkQuery{LinkID: linkID})
	if err != nil {
		return sharedvalueobjects.ID{}, normalizeAnalyticsAccessError(err)
	}
	resumeID, err := sharedvalueobjects.ParseID(link.ResumeID)
	if err != nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("resolve analytics link: %w", err)
	}
	return resumeID, nil
}

func (r *Reader) LoadLinkLifecycle(
	ctx context.Context,
	linkID sharedvalueobjects.ID,
) (analyticscontracts.LinkLifecycle, error) {
	if r == nil || r.links == nil || r.publication == nil {
		return analyticscontracts.LinkLifecycle{}, fmt.Errorf("analytics lifecycle reader is nil")
	}
	if linkID.IsZero() {
		return analyticscontracts.LinkLifecycle{}, fmt.Errorf("invalid analytics link id")
	}
	query := resumelinkqueries.PublicResumeLinkQuery{LinkID: linkID}
	link, err := r.links.GetByID(ctx, query)
	if err != nil {
		return analyticscontracts.LinkLifecycle{}, normalizeAnalyticsAccessError(err)
	}
	resumeID, err := sharedvalueobjects.ParseID(link.ResumeID)
	if err != nil {
		return analyticscontracts.LinkLifecycle{}, fmt.Errorf("load analytics lifecycle: %w", err)
	}
	state, err := r.publication.LockPublicationState(ctx, resumecontracts.ResumePublicationQuery{ResumeID: resumeID})
	if err != nil {
		return analyticscontracts.LinkLifecycle{}, normalizeAnalyticsLifecycleError(err)
	}
	locked, err := r.links.LockByID(ctx, query)
	if err != nil {
		return analyticscontracts.LinkLifecycle{}, normalizeAnalyticsAccessError(err)
	}
	lockedResumeID, err := sharedvalueobjects.ParseID(locked.ResumeID)
	if err != nil {
		return analyticscontracts.LinkLifecycle{}, fmt.Errorf("load analytics lifecycle: %w", err)
	}
	if lockedResumeID != resumeID {
		return analyticscontracts.LinkLifecycle{}, fmt.Errorf("load analytics lifecycle: link resume changed while locking")
	}
	return analyticscontracts.LinkLifecycle{
		LinkID:          linkID,
		ResumeID:        resumeID,
		LinkActive:      locked.Status == string(domain.ResumeLinkStatusActive),
		ResumePublished: state.Published,
	}, nil
}

func (r *Reader) ListResumeLinks(
	ctx context.Context,
	resumeID sharedvalueobjects.ID,
) ([]analyticscontracts.ResumeLinkRef, error) {
	if r == nil || r.links == nil {
		return nil, fmt.Errorf("analytics link reader is nil")
	}
	if resumeID.IsZero() {
		return nil, fmt.Errorf("invalid analytics resume id")
	}
	items, err := r.links.ListByResumeID(ctx, resumelinkqueries.ListResumeLinksQuery{ResumeID: resumeID})
	if err != nil {
		return nil, normalizeAnalyticsAccessError(err)
	}
	result := make([]analyticscontracts.ResumeLinkRef, 0, len(items))
	for _, item := range items {
		id, err := sharedvalueobjects.ParseID(item.ID)
		if err != nil {
			return nil, fmt.Errorf("list analytics links: %w", err)
		}
		result = append(result, analyticscontracts.ResumeLinkRef{
			ID: id, Name: item.Name, Status: item.Status, CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (r *Reader) EnsureResumeOwned(
	ctx context.Context,
	userID, resumeID sharedvalueobjects.ID,
) error {
	if r == nil || r.resumes == nil {
		return fmt.Errorf("analytics ownership reader is nil")
	}
	return normalizeAnalyticsAccessError(ownership.EnsureResumeOwned(ctx, r.resumes, resumecontracts.ResumeOwnershipQuery{ResumeID: resumeID, UserID: userID}))
}

func (r *Reader) ResumeSectionOrder(
	ctx context.Context,
	userID, resumeID sharedvalueobjects.ID,
) ([]string, error) {
	if r == nil || r.overview == nil {
		return nil, fmt.Errorf("analytics resume sections reader is nil")
	}
	sections, err := r.overview.Handle(ctx, resumecontracts.OwnedResumeSectionsQuery{ResumeID: resumeID, UserID: userID})
	if err != nil {
		return nil, normalizeAnalyticsAccessError(application.NormalizeProviderError(err))
	}
	order := []string{"about"}
	seen := map[string]bool{"about": true}
	for _, section := range sections {
		if section.Enabled && section.ID != "techStack" && !seen[section.ID] {
			order = append(order, section.ID)
			seen[section.ID] = true
		}
	}
	return order, nil
}

func (r *Reader) ResumeAnalyticsLayout(
	ctx context.Context,
	userID, resumeID sharedvalueobjects.ID,
) (analyticscontracts.ResumeAnalyticsLayout, error) {
	if r == nil || r.publicResume == nil {
		return analyticscontracts.ResumeAnalyticsLayout{}, fmt.Errorf("analytics resume reader is nil")
	}
	if err := r.EnsureResumeOwned(ctx, userID, resumeID); err != nil {
		return analyticscontracts.ResumeAnalyticsLayout{}, err
	}
	view, err := r.publicResume.Handle(ctx, resumecontracts.PublicResumeQuery{ID: resumeID})
	if err != nil {
		return analyticscontracts.ResumeAnalyticsLayout{}, normalizeAnalyticsAccessError(err)
	}
	layout, err := resumelinkqueries.BuildAnalyticsLayout(view)
	if err != nil {
		return analyticscontracts.ResumeAnalyticsLayout{}, fmt.Errorf("build analytics resume layout: %w", err)
	}
	return layout, nil
}

func normalizeAnalyticsLifecycleError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, resumecontracts.ErrResumeNotFound) {
		return fmt.Errorf("%w: %w", analyticscontracts.ErrAnalyticsSubjectNotFound, err)
	}
	return normalizeAnalyticsAccessError(err)
}

func normalizeAnalyticsAccessError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, application.ErrResumeLinkNotFound) ||
		errors.Is(err, domain.ErrEntityNotFound) ||
		errors.Is(err, application.ErrResumeNotFound) {
		return fmt.Errorf("%w: %w", analyticscontracts.ErrAnalyticsSubjectNotFound, err)
	}
	return err
}
