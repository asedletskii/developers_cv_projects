package mappers

import (
	"fmt"
	"time"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func ResumeLinkToDB(link *domain.ResumeLink, createdAt, updatedAt time.Time) (gormmodel.ResumeLink, error) {
	if link == nil {
		return gormmodel.ResumeLink{}, fmt.Errorf("resume link is nil")
	}
	return gormmodel.ResumeLink{
		ID: link.ID().UUID(), ResumeID: link.ResumeID().UUID(), Name: link.Name(),
		Status: string(link.Status()), Version: link.Version(), CreatedAt: createdAt, UpdatedAt: updatedAt,
	}, nil
}

func ResumeLinkToDomain(model gormmodel.ResumeLink) (*domain.ResumeLink, error) {
	id, err := sharedvalueobjects.IDFromUUID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("map resume link id: %w", err)
	}
	resumeID, err := sharedvalueobjects.IDFromUUID(model.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("map resume link resume id: %w", err)
	}
	link, err := domain.RestoreResumeLink(domain.ResumeLinkSnapshot{
		ID: id, ResumeID: resumeID, Name: model.Name,
		Version: model.Version, Status: domain.ResumeLinkStatus(model.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("map resume link: %w", err)
	}
	return link, nil
}
