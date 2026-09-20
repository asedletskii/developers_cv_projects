package analytics

import (
	"fmt"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
)

type LayoutSnapshotBuilder struct{}

func NewLayoutSnapshotBuilder() *LayoutSnapshotBuilder { return &LayoutSnapshotBuilder{} }

func (b *LayoutSnapshotBuilder) Build(resume analyticscontracts.ResumeAnalyticsLayout) (analyticsdomain.LayoutSnapshot, error) {
	order := make([]analyticsdomain.AnalyticsSectionID, 0, len(resume.Sections))
	hashes := make(map[analyticsdomain.AnalyticsSectionID]string, len(resume.Sections))
	for _, section := range resume.Sections {
		id, err := analyticsdomain.ParseAnalyticsSectionID(section.ID)
		if err != nil {
			return analyticsdomain.LayoutSnapshot{}, fmt.Errorf("build analytics layout: %w", err)
		}
		order = append(order, id)
		hashes[id] = section.Hash
	}
	if len(order) == 0 {
		return analyticsdomain.LayoutSnapshot{}, fmt.Errorf("build analytics layout: no visible sections")
	}
	layout, err := analyticsdomain.NewLayoutSnapshot(order, hashes, resume.Technologies)
	if err != nil {
		return analyticsdomain.LayoutSnapshot{}, fmt.Errorf("build analytics layout: %w", err)
	}
	return layout, nil
}
