package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	analyticscontracts "developers-cv/backend/internal/modules/analytics/contracts"
	"developers-cv/backend/internal/modules/analytics/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestCreatePublicVisitPersistsCompleteView(t *testing.T) {
	views := &viewRepositoryStub{}
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	subjects := &analyticsSubjectReaderStub{lifecycle: &analyticscontracts.LinkLifecycle{
		LinkID: linkID, ResumeID: resumeID, LinkActive: true, ResumePublished: true,
	}}
	handler := NewCreatePublicVisitHandler(views, &unitOfWorkStub{}, visitorIdentityStub{}, trafficClassifierStub{}, layoutSnapshotBuilderStub{}, subjects)
	identity, err := handler.Handle(context.Background(), analyticscontracts.CreatePublicVisitInput{LinkID: linkID, ResumeID: resumeID})
	require.NoError(t, err)
	require.NotNil(t, views.created)
	require.Equal(t, views.created.ID(), identity.ViewID)
	require.Equal(t, linkID, views.created.LinkID())
	require.Equal(t, resumeID, views.created.ResumeID())
	require.Equal(t, domain.TrafficHuman, views.created.Traffic())
	require.Equal(t, domain.DeviceDesktop, views.created.Device())
	require.Equal(t, []domain.AnalyticsSectionID{domain.AnalyticsSectionAbout}, views.created.Layout().Order())
}
