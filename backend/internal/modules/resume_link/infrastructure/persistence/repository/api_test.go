package repository

import (
	"reflect"
	"testing"
)

func TestResumeLinkRepositoryWriteSideAPI(t *testing.T) {
	typ := reflect.TypeOf(&ResumeLinkRepository{})
	methods := []string{
		"Create", "GetByID", "Delete", "ReplaceStatus", "CreateView",
		"GetView", "UpdateViewState", "RecordEvents", "LoadAnalyticsBatchForUpdate", "SaveAnalyticsBatch",
	}
	for _, method := range methods {
		if _, ok := typ.MethodByName(method); !ok {
			t.Fatalf("repository must expose %s", method)
		}
	}
	if _, ok := typ.MethodByName("ListByResumeID"); ok {
		t.Fatal("repository must not expose read-side ListByResumeID")
	}
}
