package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkQuery struct {
	LinkID sharedvalueobjects.ID
	Visit  *ports.CreatePublicVisitInput
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
	links                ResumeLinkPublicReader
	resume               ResumePublicQuery
	owner                ports.ResumeOwnershipReader
	visit                ports.PublicVisitCreator
	reportAnalyticsError PublicAnalyticsErrorReporter
}

func NewPublicResumeLinkQueryHandler(links ResumeLinkPublicReader, resume ResumePublicQuery) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{links: links, resume: resume}
}

func NewPublicResumeLinkQueryHandlerWithVisit(
	links ResumeLinkPublicReader,
	resume ResumePublicQuery,
	owner ports.ResumeOwnershipReader,
	visit ports.PublicVisitCreator,
	reportAnalyticsError PublicAnalyticsErrorReporter,
) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{
		links: links, resume: resume, owner: owner, visit: visit,
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
	if h == nil || h.links == nil || h.resume == nil {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("public resume link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("invalid resume link id")
	}

	link, err := h.links.GetByID(ctx, query)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, ErrResumeLinkNotFound
		}
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, err
	}
	if domain.ResumeLinkStatus(link.Status) == domain.ResumeLinkStatusInactive {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, domain.ErrResumeLinkInactive
	}
	if domain.ResumeLinkStatus(link.Status) != domain.ResumeLinkStatusActive {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, ErrResumeLinkNotFound
	}

	resumeUUID, err := uuid.Parse(link.ResumeID)
	if err != nil {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(resumeUUID)
	if err != nil {
		return resumequeries.PublicResumeView{}, sharedvalueobjects.ID{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
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
) (resumequeries.PublicResumeView, ports.VisitorIdentityResult, error) {
	view, resumeID, err := h.handle(ctx, query)
	if err != nil {
		return resumequeries.PublicResumeView{}, ports.VisitorIdentityResult{}, err
	}
	if h == nil || h.visit == nil || query.Visit == nil {
		return view, ports.VisitorIdentityResult{}, nil
	}
	visitInput := *query.Visit
	visitInput.ResumeID = resumeID
	visitInput.Resume = view
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
		return view, ports.VisitorIdentityResult{}, nil
	}
	return view, identity, nil
}

func (h *PublicResumeLinkQueryHandler) reportPublicAnalyticsFailure(failure PublicAnalyticsFailure) {
	if h.reportAnalyticsError != nil {
		h.reportAnalyticsError(failure)
	}
}
