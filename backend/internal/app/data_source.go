package app

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/google/uuid"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumereaders "developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	userqueries "developers-cv/backend/internal/modules/user/application/queries"
	userauth "developers-cv/backend/internal/modules/user/infrastructure/auth"
	userreaders "developers-cv/backend/internal/modules/user/infrastructure/persistence/readers"
	"developers-cv/backend/internal/platform/health"
	httpmiddleware "developers-cv/backend/internal/platform/http/middleware"
	"developers-cv/backend/internal/platform/postgres"
)

func newPublicResumeReader(db *gorm.DB) resumequeries.PublicResumeReader {
	return resumereaders.NewPublicResumeReader(db)
}

func newListResumesReader(db *gorm.DB) resumequeries.ListResumesReader {
	return resumereaders.NewListResumesReader(db)
}

func newResumeMetadataReader(db *gorm.DB) resumequeries.ResumeMetadataReader {
	return resumereaders.NewResumeMetadataReader(db)
}

func newResumeOverviewReader(db *gorm.DB) resumequeries.OverviewReader {
	return resumereaders.NewOverviewReader(db)
}

func newResumeReader(deps resumereaders.ResumeReaderDependencies) resumequeries.ResumeReader {
	return resumereaders.NewResumeReader(deps)
}

func newGetProfileReader(db *gorm.DB) userqueries.GetProfileReader {
	return userreaders.NewGetProfileReader(db)
}

func newAuthenticate(authenticator *userauth.MockAuthenticator) httpmiddleware.AuthenticateFunc {
	return func(ctx context.Context, request *http.Request) (uuid.UUID, error) {
		userID, err := authenticator.Authenticate(ctx, request)
		if err != nil {
			if errors.Is(err, userauth.ErrUnauthenticated) {
				return uuid.Nil, httpmiddleware.ErrUnauthenticated
			}
			return uuid.Nil, err
		}
		return userID, nil
	}
}

func newHealthService(db *gorm.DB) *health.Service {
	return health.NewService(postgres.NewHealthChecker(db))
}
