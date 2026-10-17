package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	resumequeries "developers-cv/backend/internal/modules/resume/application/queries"
	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type publicVisitViewsStub struct{ created *domain.View }

func (s *publicVisitViewsStub) CreateView(_ context.Context, view *domain.View) error {
	s.created = view
	return nil
}
func (s *publicVisitViewsStub) GetView(context.Context, sharedvalueobjects.ID) (*domain.View, error) {
	return nil, nil
}
func (s *publicVisitViewsStub) UpdateViewState(context.Context, sharedvalueobjects.ID, domain.ViewStateUpdate) error {
	return nil
}

type publicVisitUnitOfWorkStub struct{}

func (publicVisitUnitOfWorkStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type publicVisitIdentityStub struct{}

func (publicVisitIdentityStub) Identify(_ time.Time, _ ports.PublicVisitInput) (ports.VisitorIdentityResult, error) {
	return ports.VisitorIdentityResult{VisitorID: sharedvalueobjects.NewID(), IPHash: make([]byte, 32)}, nil
}

type publicVisitClassifierStub struct{}

func (publicVisitClassifierStub) Classify(context.Context, ports.TrafficClassificationInput) (ports.TrafficClassificationResult, error) {
	return ports.TrafficClassificationResult{Traffic: domain.TrafficHuman, Device: domain.DeviceDesktop, Browser: domain.BrowserChrome, OS: domain.OSMacOS}, nil
}

type publicVisitLayoutStub struct{}

func (publicVisitLayoutStub) Build(resumequeries.PublicResumeView) (domain.LayoutSnapshot, error) {
	return domain.NewLayoutSnapshot([]resumedomain.SectionID{resumedomain.SectionAbout}, map[resumedomain.SectionID]string{resumedomain.SectionAbout: "0123456789abcdef"})
}

func TestCreatePublicVisitPersistsCompleteView(t *testing.T) {
	views := &publicVisitViewsStub{}
	handler := NewCreatePublicVisitHandler(views, publicVisitUnitOfWorkStub{}, publicVisitIdentityStub{}, publicVisitClassifierStub{}, publicVisitLayoutStub{})
	linkID := sharedvalueobjects.NewID()
	view, identity, err := handler.Handle(context.Background(), CreatePublicVisitInput{LinkID: linkID})
	require.NoError(t, err)
	require.NotNil(t, view)
	require.Same(t, view, views.created)
	require.Equal(t, view.ID(), identity.ViewID)
	require.Equal(t, linkID, view.LinkID())
	require.Equal(t, domain.TrafficHuman, view.Traffic())
	require.Equal(t, domain.DeviceDesktop, view.Device())
	require.Equal(t, []resumedomain.SectionID{resumedomain.SectionAbout}, view.Layout().Order())
}
