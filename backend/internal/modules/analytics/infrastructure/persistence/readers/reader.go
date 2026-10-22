package readers

import (
	"gorm.io/gorm"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
)

type AnalyticsReader struct {
	db       *gorm.DB
	subjects analyticscontracts.AnalyticsSubjectReader
}

func NewAnalyticsReader(db *gorm.DB, subjects analyticscontracts.AnalyticsSubjectReader) *AnalyticsReader {
	return &AnalyticsReader{db: db, subjects: subjects}
}
