package repository

import (
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

func TestRepositoryModelMappersRejectNilEntities(t *testing.T) {
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "activity",
			call: func() error {
				_, err := activityModel(nil)
				return err
			},
		},
		{
			name: "course",
			call: func() error {
				_, err := courseModel(nil)
				return err
			},
		},
		{
			name: "education",
			call: func() error {
				_, err := educationModel(nil)
				return err
			},
		},
		{
			name: "experience",
			call: func() error {
				_, err := experienceModel(nil)
				return err
			},
		},
		{
			name: "project",
			call: func() error {
				_, err := projectDomainModel(nil)
				return err
			},
		},
		{
			name: "contacts",
			call: func() error {
				_, _, _, err := contactsModels(nil)
				return err
			},
		},
		{
			name: "profile",
			call: func() error {
				_, _, err := profileModels(nil)
				return err
			},
		},
		{
			name: "overview",
			call: func() error {
				_, _, err := overviewModels(valueobjects.ID{}, nil)
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
