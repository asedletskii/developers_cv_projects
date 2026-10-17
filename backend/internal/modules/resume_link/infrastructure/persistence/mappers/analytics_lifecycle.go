package mappers

import (
	"fmt"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
)

type AnalyticsLifecycle struct {
	LinkStatus   string `gorm:"column:status"`
	ResumeStatus string `gorm:"column:resume_status"`
}

func AnalyticsLifecycleToDomain(
	model AnalyticsLifecycle,
) (domain.ResumeLinkStatus, resumedomain.ResumeStatus, error) {
	linkStatus := domain.ResumeLinkStatus(model.LinkStatus)
	if !linkStatus.IsValid() {
		return "", "", fmt.Errorf("map analytics lifecycle: invalid resume link status %q", model.LinkStatus)
	}

	resumeStatus, err := resumedomain.ParseResumeStatus(model.ResumeStatus)
	if err != nil {
		return "", "", fmt.Errorf("map analytics lifecycle: invalid resume status %q", model.ResumeStatus)
	}

	return linkStatus, resumeStatus, nil
}
