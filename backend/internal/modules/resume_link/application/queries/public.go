package queries

import (
	"context"
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	application "developers-cv/backend/internal/modules/resume_link/application"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkQuery struct {
	LinkID sharedvalueobjects.ID
	Visit  *analyticscontracts.CreatePublicVisitInput
}

type ResumeLinkPublicView struct {
	ID       string
	ResumeID string
	Status   string
}

type ResumeLinkPublicReader interface {
	GetByID(context.Context, PublicResumeLinkQuery) (ResumeLinkPublicView, error)
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
	resume               resumecontracts.PublicResumeQueryHandler
	owner                resumecontracts.ResumeOwnershipChecker
	visit                analyticscontracts.PublicVisitCreator
	reportAnalyticsError PublicAnalyticsErrorReporter
}

func NewPublicResumeLinkQueryHandler(linkResolver PublicResumeLinkResolver, resume resumecontracts.PublicResumeQueryHandler) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{linkResolver: linkResolver, resume: resume}
}

func NewPublicResumeLinkQueryHandlerWithVisit(
	linkResolver PublicResumeLinkResolver,
	resume resumecontracts.PublicResumeQueryHandler,
	owner resumecontracts.ResumeOwnershipChecker,
	visit analyticscontracts.PublicVisitCreator,
	reportAnalyticsError PublicAnalyticsErrorReporter,
) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{
		linkResolver: linkResolver, resume: resume, owner: owner, visit: visit,
		reportAnalyticsError: reportAnalyticsError,
	}
}

func (h *PublicResumeLinkQueryHandler) Handle(ctx context.Context, query PublicResumeLinkQuery) (PublicResumeResult, error) {
	view, _, err := h.load(ctx, query)
	if err != nil {
		return PublicResumeResult{}, err
	}
	return mapPublicResume(view), nil
}

func (h *PublicResumeLinkQueryHandler) load(
	ctx context.Context,
	query PublicResumeLinkQuery,
) (resumecontracts.PublicResumeView, sharedvalueobjects.ID, error) {
	if h == nil || h.linkResolver == nil || h.resume == nil {
		return resumecontracts.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("public resume link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return resumecontracts.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("invalid resume link id")
	}

	resumeID, err := h.linkResolver.Resolve(ctx, query.LinkID)
	if err != nil {
		return resumecontracts.PublicResumeView{}, sharedvalueobjects.ID{}, err
	}

	view, err := h.resume.Handle(ctx, resumecontracts.PublicResumeQuery{ID: resumeID})
	if err != nil {
		return resumecontracts.PublicResumeView{}, sharedvalueobjects.ID{}, application.NormalizeProviderError(err)
	}
	return view, resumeID, nil
}

func (h *PublicResumeLinkQueryHandler) HandleWithVisit(
	ctx context.Context,
	query PublicResumeLinkQuery,
) (PublicResumeResult, analyticscontracts.VisitorIdentityResult, error) {
	view, resumeID, err := h.load(ctx, query)
	if err != nil {
		return PublicResumeResult{}, analyticscontracts.VisitorIdentityResult{}, err
	}
	result := mapPublicResume(view)
	if h == nil || h.visit == nil || query.Visit == nil {
		return result, analyticscontracts.VisitorIdentityResult{}, nil
	}
	visitInput := *query.Visit
	visitInput.ResumeID = resumeID
	layout, err := BuildAnalyticsLayout(view)
	if err != nil {
		h.reportPublicAnalyticsFailure(PublicAnalyticsFailure{
			Stage:    PublicAnalyticsStageLayoutBuild,
			LinkID:   query.LinkID,
			ResumeID: resumeID,
			Err:      err,
		})
		return result, analyticscontracts.VisitorIdentityResult{}, nil
	}
	visitInput.ResumeLayout = layout
	if !visitInput.Visitor.OwnerID.IsZero() {
		if h.owner == nil {
			visitInput.Visitor.OwnerID = sharedvalueobjects.ID{}
		} else {
			owned, ownerErr := h.owner.OwnsResume(ctx, resumecontracts.ResumeOwnershipQuery{ResumeID: visitInput.ResumeID, UserID: visitInput.Visitor.OwnerID})
			if ownerErr != nil {
				h.reportPublicAnalyticsFailure(PublicAnalyticsFailure{
					Stage:    PublicAnalyticsStageOwnerResolution,
					LinkID:   query.LinkID,
					ResumeID: visitInput.ResumeID,
					Err:      ownerErr,
				})
				visitInput.Visitor.OwnerID = sharedvalueobjects.ID{}
			} else if !owned {
				visitInput.Visitor.OwnerID = sharedvalueobjects.ID{}
			}
		}
	}
	identity, err := h.visit.Handle(ctx, visitInput)
	if err != nil {
		h.reportPublicAnalyticsFailure(PublicAnalyticsFailure{
			Stage:    PublicAnalyticsStageVisitCreation,
			LinkID:   query.LinkID,
			ResumeID: visitInput.ResumeID,
			Err:      err,
		})
		return result, analyticscontracts.VisitorIdentityResult{}, nil
	}
	return result, identity, nil
}

func (h *PublicResumeLinkQueryHandler) reportPublicAnalyticsFailure(failure PublicAnalyticsFailure) {
	if h.reportAnalyticsError != nil {
		h.reportAnalyticsError(failure)
	}
}
