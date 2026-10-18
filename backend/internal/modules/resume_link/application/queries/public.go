package queries

import (
	"context"
	"errors"
	"fmt"

	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkQuery struct {
	LinkID sharedvalueobjects.ID
	Visit  *analyticsports.CreatePublicVisitInput
}

type ResumeLinkPublicView struct {
	ID       string
	ResumeID string
	Status   string
}

type ResumeLinkPublicReader interface {
	GetByID(context.Context, PublicResumeLinkQuery) (ResumeLinkPublicView, error)
}

type ResumePublicQuery interface {
	Handle(context.Context, resumequeries.PublicResumeQuery) (resumequeries.PublicResumeView, error)
}

type PublicAnalyticsFailureStage string

const (
	PublicAnalyticsStageOwnerResolution PublicAnalyticsFailureStage = "owner_resolution"
	PublicAnalyticsStageLayoutBuild     PublicAnalyticsFailureStage = "layout_build"
	PublicAnalyticsStageVisitCreation   PublicAnalyticsFailureStage = "visit_creation"
)

type PublicAnalyticsFailure struct {
	Stage    PublicAnalyticsFailureStage
	LinkID   sharedvalueobjects.ID
	ResumeID sharedvalueobjects.ID
	Err      error
}

type PublicAnalyticsErrorReporter func(PublicAnalyticsFailure)

type PublicResumeLinkQueryHandler struct {
	linkResolver         PublicResumeLinkResolver
	resume               ResumePublicQuery
	owner                ports.ResumeOwnershipReader
	visit                analyticsports.PublicVisitCreator
	reportAnalyticsError PublicAnalyticsErrorReporter
}

func NewPublicResumeLinkQueryHandler(linkResolver PublicResumeLinkResolver, resume ResumePublicQuery) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{linkResolver: linkResolver, resume: resume}
}

func NewPublicResumeLinkQueryHandlerWithVisit(
	linkResolver PublicResumeLinkResolver,
	resume ResumePublicQuery,
	owner ports.ResumeOwnershipReader,
	visit analyticsports.PublicVisitCreator,
	reportAnalyticsError PublicAnalyticsErrorReporter,
) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{
		linkResolver: linkResolver, resume: resume, owner: owner, visit: visit,
		reportAnalyticsError: reportAnalyticsError,
	}
}

func (h *PublicResumeLinkQueryHandler) Handle(ctx context.Context, query PublicResumeLinkQuery) (resumequeries.PublicResumeView, error) {
	view, _, err := h.handle(ctx, query)
	return view, err
}

func (h *PublicResumeLinkQueryHandler) handle(
	ctx context.Context,
	query PublicResumeLinkQuery,
) (resumequeries.PublicResumeView, sharedvalueobjects.ID, error) {
	if h == nil || h.linkResolver == nil || h.resume == nil {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("public resume link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("invalid resume link id")
	}

	resumeID, err := h.linkResolver.Resolve(ctx, query.LinkID)
	if err != nil {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, err
	}

	view, err := h.resume.Handle(ctx, resumequeries.PublicResumeQuery{ID: resumeID})
	if query.Visit != nil {
		visit := *query.Visit
		visit.ResumeID = resumeID
		query.Visit = &visit
	}
	if err != nil {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, err
	}
	return view, resumeID, nil
}

func (h *PublicResumeLinkQueryHandler) HandleWithVisit(
	ctx context.Context,
	query PublicResumeLinkQuery,
) (resumequeries.PublicResumeView, analyticsports.VisitorIdentityResult, error) {
	view, resumeID, err := h.handle(ctx, query)
	if err != nil {
		return resumequeries.PublicResumeView{}, analyticsports.VisitorIdentityResult{}, err
	}
	if h == nil || h.visit == nil || query.Visit == nil {
		return view, analyticsports.VisitorIdentityResult{}, nil
	}
	visitInput := *query.Visit
	visitInput.ResumeID = resumeID
	layout, err := resumequeries.ToResumeAnalyticsLayout(view)
	if err != nil {
		h.reportPublicAnalyticsFailure(PublicAnalyticsFailure{
			Stage:    PublicAnalyticsStageLayoutBuild,
			LinkID:   query.LinkID,
			ResumeID: resumeID,
			Err:      err,
		})
		return view, analyticsports.VisitorIdentityResult{}, nil
	}
	visitInput.ResumeLayout = analyticsports.ResumeAnalyticsLayout{
		Sections:     make([]analyticsports.ResumeAnalyticsSection, 0, len(layout.Sections)),
		Technologies: append([]string{}, layout.Technologies...),
	}
	for _, section := range layout.Sections {
		visitInput.ResumeLayout.Sections = append(visitInput.ResumeLayout.Sections, analyticsports.ResumeAnalyticsSection{
			ID: section.ID, Hash: section.Hash,
		})
	}
	if !visitInput.Visitor.OwnerID.IsZero() {
		if h.owner == nil {
			visitInput.Visitor.OwnerID = sharedvalueobjects.ID{}
		} else {
			owned, ownerErr := h.owner.OwnsResume(ctx, resumequeries.GetResumeQuery{ResumeID: visitInput.ResumeID, UserID: visitInput.Visitor.OwnerID})
			if ownerErr != nil {
				if !errors.Is(ownerErr, resumequeries.ErrResumeNotFound) {
					h.reportPublicAnalyticsFailure(PublicAnalyticsFailure{
						Stage:    PublicAnalyticsStageOwnerResolution,
						LinkID:   query.LinkID,
						ResumeID: visitInput.ResumeID,
						Err:      ownerErr,
					})
				}
				visitInput.Visitor.OwnerID = sharedvalueobjects.ID{}
			} else if !owned {
				visitInput.Visitor.OwnerID = sharedvalueobjects.ID{}
			}
		}
	}
	_, identity, err := h.visit.Handle(ctx, visitInput)
	if err != nil {
		h.reportPublicAnalyticsFailure(PublicAnalyticsFailure{
			Stage:    PublicAnalyticsStageVisitCreation,
			LinkID:   query.LinkID,
			ResumeID: visitInput.ResumeID,
			Err:      err,
		})
		return view, analyticsports.VisitorIdentityResult{}, nil
	}
	return view, identity, nil
}

func (h *PublicResumeLinkQueryHandler) reportPublicAnalyticsFailure(failure PublicAnalyticsFailure) {
	if h.reportAnalyticsError != nil {
		h.reportAnalyticsError(failure)
	}
}
