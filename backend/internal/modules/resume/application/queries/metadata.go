package queries

import (
	"context"
	"time"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetResumeMetadataQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type ResumeMetadataReader interface {
	GetResumeMetadata(context.Context, GetResumeMetadataQuery) (ResumeMetadataView, error)
}

type ResumeMetadataView struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
