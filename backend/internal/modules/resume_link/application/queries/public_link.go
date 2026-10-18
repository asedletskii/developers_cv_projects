package queries

import (
	"context"
	"errors"
	"fmt"

	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicResumeLinkResolver interface {
	Resolve(context.Context, sharedvalueobjects.ID) (sharedvalueobjects.ID, error)
}

type ActivePublicResumeLinkResolver struct {
	links ResumeLinkPublicReader
}

func NewActivePublicResumeLinkResolver(links ResumeLinkPublicReader) *ActivePublicResumeLinkResolver {
	return &ActivePublicResumeLinkResolver{links: links}
}

func (r *ActivePublicResumeLinkResolver) Resolve(
	ctx context.Context,
	linkID sharedvalueobjects.ID,
) (sharedvalueobjects.ID, error) {
	if r == nil || r.links == nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("public resume link resolver dependencies are nil")
	}
	if linkID.IsZero() {
		return sharedvalueobjects.ID{}, fmt.Errorf("invalid resume link id")
	}

	link, err := r.links.GetByID(ctx, PublicResumeLinkQuery{LinkID: linkID})
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return sharedvalueobjects.ID{}, ErrResumeLinkNotFound
		}
		return sharedvalueobjects.ID{}, err
	}
	switch domain.ResumeLinkStatus(link.Status) {
	case domain.ResumeLinkStatusActive:
	case domain.ResumeLinkStatusInactive:
		return sharedvalueobjects.ID{}, domain.ErrResumeLinkInactive
	default:
		return sharedvalueobjects.ID{}, ErrResumeLinkNotFound
	}

	resumeID, err := sharedvalueobjects.ParseID(link.ResumeID)
	if err != nil {
		return sharedvalueobjects.ID{}, fmt.Errorf("resolve resume link: %w", err)
	}
	return resumeID, nil
}
