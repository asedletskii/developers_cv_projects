package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestCreatePublicVisitPersistsCompleteView(t *testing.T) {
	views := &viewRepositoryStub{}
	handler := NewCreatePublicVisitHandler(views, &unitOfWorkStub{}, visitorIdentityStub{}, trafficClassifierStub{}, layoutSnapshotBuilderStub{})
	linkID := sharedvalueobjects.NewID()
	view, identity, err := handler.Handle(context.Background(), CreatePublicVisitInput{LinkID: linkID})
	require.NoError(t, err)
	require.NotNil(t, view)
	require.Same(t, view, views.created)
	require.Equal(t, view.ID(), identity.ViewID)
	require.Equal(t, linkID, view.LinkID())
	require.Equal(t, domain.TrafficHuman, view.Traffic())
	require.Equal(t, domain.DeviceDesktop, view.Device())
	require.Equal(t, []domain.AnalyticsSectionID{domain.AnalyticsSectionAbout}, view.Layout().Order())
}
