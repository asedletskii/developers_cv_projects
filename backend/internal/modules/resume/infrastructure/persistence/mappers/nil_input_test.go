package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestRepositoryModelMappersRejectNilEntities(t *testing.T) {
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "resume",
			call: func() error {
				_, err := ResumeToDB(nil, time.Time{}, time.Time{})
				return err
			},
		},
		{
			name: "activity",
			call: func() error {
				_, err := ActivityToDB(nil)
				return err
			},
		},
		{
			name: "course",
			call: func() error {
				_, err := CourseToDB(nil)
				return err
			},
		},
		{
			name: "education",
			call: func() error {
				_, err := EducationToDB(nil)
				return err
			},
		},
		{
			name: "experience",
			call: func() error {
				_, err := ExperienceToDB(nil)
				return err
			},
		},
		{
			name: "project",
			call: func() error {
				_, err := ProjectToDB(nil)
				return err
			},
		},
		{
			name: "contacts",
			call: func() error {
				_, _, _, err := ContactsToDBModels(nil)
				return err
			},
		},
		{
			name: "profile",
			call: func() error {
				_, _, err := ProfileToDBModels(nil)
				return err
			},
		},
		{
			name: "overview",
			call: func() error {
				_, _, err := OverviewToDBModels(identity.ID{}, nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			require.Error(t, err)
		})
	}
}
