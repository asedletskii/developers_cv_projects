package repository

import (
	"reflect"
	"testing"
)

func TestResumeLinkRepositoryWriteSideAPI(t *testing.T) {
	typ := reflect.TypeOf(&ResumeLinkRepository{})
	for _, method := range []string{"Create", "GetByID", "Delete", "ReplaceStatus", "CreateView", "GetView", "UpdateViewState", "RecordEvents", "ApplyBatch"} {
		if _, ok := typ.MethodByName(method); !ok {
			t.Fatalf("repository must expose %s", method)
		}
	}
	if _, ok := typ.MethodByName("ListByResumeID"); ok {
		t.Fatal("repository must not expose read-side ListByResumeID")
	}
}
