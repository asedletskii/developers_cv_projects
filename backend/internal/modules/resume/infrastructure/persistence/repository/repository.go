package repository

import (
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/commands"
)

type ResumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) *ResumeRepository { return &ResumeRepository{db: db} }

var _ commands.CreateResumeRepository = (*ResumeRepository)(nil)
var _ commands.UpdateResumeRepository = (*ResumeRepository)(nil)
var _ commands.GetResumeRepository = (*ResumeRepository)(nil)
var _ commands.UpdateResumeTitleRepository = (*ResumeRepository)(nil)
