package mappers

import (
	"fmt"

	"developers-cv/backend/internal/modules/analytics/domain"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
)

func ClickToDB(click *domain.AnalyticsClick) (gormmodel.AnalyticsClick, error) {
	if click == nil {
		return gormmodel.AnalyticsClick{}, fmt.Errorf("click is nil")
	}
	var target *string
	if click.Target() != "" {
		value := click.Target()
		target = &value
	}
	return gormmodel.AnalyticsClick{
		ViewID: click.ViewID().UUID(), LinkID: click.LinkID().UUID(), Section: string(click.Section()),
		XPct: click.XPct(), YPct: click.YPct(), ViewportW: click.ViewportW(), LayoutHash: click.LayoutHash(),
		Target: target, OccurredAt: click.OccurredAt(),
	}, nil
}
