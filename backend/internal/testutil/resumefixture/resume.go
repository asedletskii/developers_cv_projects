package resumefixture

import (
	"context"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
)

type Repository struct {
	Resume *domain.Resume
	Saved  bool
}

func (r *Repository) Get(context.Context, ports.GetResumeParams) (*domain.Resume, error) {
	return r.Resume, nil
}

func (r *Repository) SaveSection(_ context.Context, resume *domain.Resume, _ ports.ResumeMutationSection) error {
	r.Resume = resume
	r.Saved = true
	return nil
}

func (*Repository) Create(context.Context, *domain.Resume) error { return nil }

func (*Repository) LockForDelete(context.Context, ports.DeleteResumeParams) error { return nil }

func (*Repository) Delete(context.Context, ports.DeleteResumeParams) error { return nil }

type UnitOfWork struct{}

func (UnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
