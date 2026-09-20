package repository

import "gorm.io/gorm"

type baseRepository struct {
	db *gorm.DB
}

func newBaseRepository(db *gorm.DB) *baseRepository { return &baseRepository{db: db} }

type ResumeRepository struct{ *baseRepository }

func NewResumeRepository(db *gorm.DB) *ResumeRepository {
	return &ResumeRepository{newBaseRepository(db)}
}
