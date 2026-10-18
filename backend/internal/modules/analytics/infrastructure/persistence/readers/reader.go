package readers

import (
	"gorm.io/gorm"
)

type AnalyticsReader struct {
	db *gorm.DB
}

func NewAnalyticsReader(db *gorm.DB) *AnalyticsReader {
	return &AnalyticsReader{db: db}
}
