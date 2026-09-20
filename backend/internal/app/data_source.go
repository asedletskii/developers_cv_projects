package app

import (
	"gorm.io/gorm"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userreaders "developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	"developers-cv/backend/internal/platform/health"
	"developers-cv/backend/internal/platform/postgres"
)

func newListResumesReader(db *gorm.DB) resumequeries.ListResumesReader {
	return resumereaders.NewListResumesReader(db)
}
func newGetProfileReader(db *gorm.DB) userqueries.GetProfileReader {
	return userreaders.NewGetProfileReader(db)
}
func newHealthService(db *gorm.DB) *health.Service {
	return health.NewService(postgres.NewHealthChecker(db))
}
