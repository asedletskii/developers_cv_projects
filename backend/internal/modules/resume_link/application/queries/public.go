package queries

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume_link/application/commands"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkQuery struct {
	LinkID sharedvalueobjects.ID
	Visit  *commands.CreatePublicVisitInput
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

type PublicResumeLinkQueryHandler struct {
	links  ResumeLinkPublicReader
	resume ResumePublicQuery
	owner  ports.ResumeOwnershipReader
	visit  *commands.CreatePublicVisitHandler
	logger *zap.Logger
}

func NewPublicResumeLinkQueryHandler(links ResumeLinkPublicReader, resume ResumePublicQuery) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{links: links, resume: resume}
}

func NewPublicResumeLinkQueryHandlerWithVisit(links ResumeLinkPublicReader, resume ResumePublicQuery, owner ports.ResumeOwnershipReader, visit *commands.CreatePublicVisitHandler, logger *zap.Logger) *PublicResumeLinkQueryHandler {
	return &PublicResumeLinkQueryHandler{links: links, resume: resume, owner: owner, visit: visit, logger: logger}
}

func (h *PublicResumeLinkQueryHandler) Handle(ctx context.Context, query PublicResumeLinkQuery) (resumequeries.PublicResumeView, error) {
	if h == nil || h.links == nil || h.resume == nil {
		return resumequeries.PublicResumeView{}, fmt.Errorf("public resume link query handler dependencies are nil")
	}
	if query.LinkID.IsZero() {
		return resumequeries.PublicResumeView{}, fmt.Errorf("invalid resume link id")
	}

	link, err := h.links.GetByID(ctx, query)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return resumequeries.PublicResumeView{}, ErrResumeLinkNotFound
		}
		return resumequeries.PublicResumeView{}, err
	}
	if domain.ResumeLinkStatus(link.Status) == domain.ResumeLinkStatusInactive {
		return resumequeries.PublicResumeView{}, domain.ErrResumeLinkInactive
	}
	if domain.ResumeLinkStatus(link.Status) != domain.ResumeLinkStatusActive {
		return resumequeries.PublicResumeView{}, ErrResumeLinkNotFound
	}

	resumeUUID, err := uuid.Parse(link.ResumeID)
	if err != nil {
		return resumequeries.PublicResumeView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(resumeUUID)
	if err != nil {
		return resumequeries.PublicResumeView{}, fmt.Errorf("resolve resume link: invalid resume id: %w", err)
	}

	view, err := h.resume.Handle(ctx, resumequeries.PublicResumeQuery{ID: resumeID})
	if query.Visit != nil {
		query.Visit.ResumeID = resumeID
	}
	if err != nil {
		return resumequeries.PublicResumeView{}, err
	}
	return view, nil
}

func (h *PublicResumeLinkQueryHandler) HandleWithVisit(ctx context.Context, query PublicResumeLinkQuery) (resumequeries.PublicResumeView, ports.VisitorIdentityResult, error) {
	view, err := h.Handle(ctx, query)
	if err != nil {
		return resumequeries.PublicResumeView{}, ports.VisitorIdentityResult{}, err
	}
	if h == nil || h.visit == nil || query.Visit == nil {
		return view, ports.VisitorIdentityResult{}, nil
	}
	query.Visit.Resume = view
	if !query.Visit.Visitor.OwnerID.IsZero() && h.owner != nil {
		owned, ownerErr := h.owner.GetOwned(ctx, resumequeries.GetResumeQuery{ResumeID: query.Visit.ResumeID, UserID: query.Visit.Visitor.OwnerID})
		if ownerErr != nil {
			if !errors.Is(ownerErr, resumequeries.ErrResumeNotFound) && h.logger != nil {
				h.logger.Error("analytics owner resolution failed",
					zap.Error(ownerErr),
					zap.String("link_id", query.LinkID.String()),
					zap.String("resume_id", query.Visit.ResumeID.String()),
					zap.String("analytics_stage", "owner_resolution"),
				)
			}
			query.Visit.Visitor.OwnerID = sharedvalueobjects.ID{}
		} else if owned.ID == "" {
			query.Visit.Visitor.OwnerID = sharedvalueobjects.ID{}
		}
	}
	_, identity, err := h.visit.Handle(ctx, *query.Visit)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("analytics visit persistence failed",
				zap.Error(err),
				zap.String("link_id", query.LinkID.String()),
				zap.String("resume_id", query.Visit.ResumeID.String()),
				zap.String("analytics_stage", "visit_creation"),
			)
		}
		return view, ports.VisitorIdentityResult{}, nil
	}
	return view, identity, nil
}
