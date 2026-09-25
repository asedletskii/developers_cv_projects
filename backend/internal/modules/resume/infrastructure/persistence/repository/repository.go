package repository

import "gorm.io/gorm"

type baseRepository struct {
	db *gorm.DB
}

func newBaseRepository(db *gorm.DB) *baseRepository { return &baseRepository{db: db} }

type ResumeRootRepository struct{ *baseRepository }
type ProfileRepository struct{ *baseRepository }
type ContactsRepository struct{ *baseRepository }
type SkillsRepository struct{ *baseRepository }
type LanguagesRepository struct{ *baseRepository }
type TechnologiesRepository struct{ *baseRepository }
type ExperiencesRepository struct{ *baseRepository }
type EducationRepository struct{ *baseRepository }
type CourseRepository struct{ *baseRepository }
type ActivityRepository struct{ *baseRepository }
type ProjectRepository struct{ *baseRepository }

func NewResumeRootRepository(db *gorm.DB) *ResumeRootRepository {
	return &ResumeRootRepository{newBaseRepository(db)}
}
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{newBaseRepository(db)}
}
func NewContactsRepository(db *gorm.DB) *ContactsRepository {
	return &ContactsRepository{newBaseRepository(db)}
}
func NewSkillsRepository(db *gorm.DB) *SkillsRepository {
	return &SkillsRepository{newBaseRepository(db)}
}
func NewLanguagesRepository(db *gorm.DB) *LanguagesRepository {
	return &LanguagesRepository{newBaseRepository(db)}
}
func NewTechnologiesRepository(db *gorm.DB) *TechnologiesRepository {
	return &TechnologiesRepository{newBaseRepository(db)}
}
func NewExperiencesRepository(db *gorm.DB) *ExperiencesRepository {
	return &ExperiencesRepository{newBaseRepository(db)}
}
func NewEducationRepository(db *gorm.DB) *EducationRepository {
	return &EducationRepository{newBaseRepository(db)}
}
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{newBaseRepository(db)}
}
func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{newBaseRepository(db)}
}
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{newBaseRepository(db)}
}
