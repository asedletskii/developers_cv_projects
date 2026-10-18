package repository

import (
	"reflect"
	"testing"
)

func TestAnalyticsRepositoryWriteSideAPI(t *testing.T) {
	typ := reflect.TypeOf(&AnalyticsRepository{})
	methods := []string{
		"CreateView", "GetView", "UpdateViewState", "RecordEvents", "LoadAnalyticsBatchForUpdate",
		"SaveAnalyticsBatch", "DeleteExpiredClicks", "DeleteExpiredViews",
	}
	for _, method := range methods {
		if _, ok := typ.MethodByName(method); !ok {
			t.Fatalf("analytics repository must expose %s", method)
		}
	}
}
