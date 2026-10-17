package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/resume_link/domain"
	"developers-cv/backend/internal/modules/resume_link/infrastructure/persistence/gormmodel"
)

func AnalyticsClickToDB(click *domain.AnalyticsClick) (gormmodel.ResumeLinkClick, error) {
	if click == nil {
		return gormmodel.ResumeLinkClick{}, fmt.Errorf("analytics click is nil")
	}
	var target *string
	if click.Target() != "" {
		value := click.Target()
		target = &value
	}
	return gormmodel.ResumeLinkClick{
		ViewID: click.ViewID().UUID(), LinkID: click.LinkID().UUID(), Section: string(click.Section()),
		XPct: click.XPct(), YPct: click.YPct(), ViewportW: click.ViewportW(), LayoutHash: click.LayoutHash(),
		Target: target, OccurredAt: click.OccurredAt(),
	}, nil
}
