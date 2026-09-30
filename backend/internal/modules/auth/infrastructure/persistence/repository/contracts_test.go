package repository_test

import (
	"testing"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/infrastructure/persistence/repository"
)

func TestRepositoryContracts(t *testing.T) {
	var _ ports.IdentityRepository = (*repository.IdentityRepository)(nil)
	var _ ports.SessionRepository = (*repository.SessionRepository)(nil)
	var _ ports.AuthTransactionRepository = (*repository.AuthTransactionRepository)(nil)
}
