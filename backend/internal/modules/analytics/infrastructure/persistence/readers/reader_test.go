package readers

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/analytics/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestGetStatisticsReturnsErrorWhenReaderDatabaseIsUnavailable(t *testing.T) {
	reader := NewAnalyticsReader(nil, nil)

	_, err := reader.GetLinkStatistics(context.Background(), queries.LinkAnalyticsQuery{
		LinkID: sharedvalueobjects.NewID(),
	})
	require.Error(t, err)
}

func TestAnalyticsReaderExposesAnalyticsQueries(t *testing.T) {
	typ := reflect.TypeOf(&AnalyticsReader{})
	for _, method := range []string{"GetLinkStatistics", "GetDaily", "GetVisits", "GetSections", "GetTechnologies", "GetResumeStatistics", "GetClicks"} {
		if _, ok := typ.MethodByName(method); !ok {
			t.Fatalf("analytics reader must expose query %s", method)
		}
	}
}
