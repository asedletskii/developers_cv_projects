package mappers

import (
	"fmt"

	"github.com/google/uuid"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	resumelinkdomain "developers-cv/backend/internal/modules/resume_link/domain"
)

type AnalyticsLifecycle struct {
	ResumeID     uuid.UUID `gorm:"column:resume_id"`
	LinkStatus   string    `gorm:"column:status"`
	ResumeStatus string    `gorm:"column:resume_status"`
}

func AnalyticsLifecycleToDomain(
	model AnalyticsLifecycle,
) (resumelinkdomain.ResumeLinkStatus, resumedomain.ResumeStatus, error) {
	linkStatus := resumelinkdomain.ResumeLinkStatus(model.LinkStatus)
	if !linkStatus.IsValid() {
		return "", "", fmt.Errorf("map analytics lifecycle: invalid resume link status %q", model.LinkStatus)
	}

	resumeStatus, err := resumedomain.ParseResumeStatus(model.ResumeStatus)
	if err != nil {
		return "", "", fmt.Errorf("map analytics lifecycle: invalid resume status %q", model.ResumeStatus)
	}

	return linkStatus, resumeStatus, nil
}
