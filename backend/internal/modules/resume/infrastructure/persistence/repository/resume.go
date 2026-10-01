package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/gormmodel"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/mappers"
	"developers-cv/backend/internal/platform/database"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func (r *ResumeRepository) Get(ctx context.Context, params ports.GetResumeParams) (*domain.Resume, error) {
	return r.loadResume(ctx, params.ResumeID, params.UserID)
}

func (r *ResumeRepository) Create(ctx context.Context, resume *domain.Resume) error {
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("create resume: %w", err)
	}
	if resume == nil {
		return fmt.Errorf("resume is nil")
	}
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()
	m, err := mappers.ResumeToDB(resume, createdAt, updatedAt)
	if err != nil {
		return fmt.Errorf("create resume: %w", err)
	}
	if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("create resume: %w", mapDBError(err))
	}
	if resume.Profile() != nil {
		pm, fs, err := mappers.ProfileToDBModels(resume.Profile())
		if err != nil {
			return fmt.Errorf("create resume profile: %w", err)
		}
		if err := tx.WithContext(ctx).Create(&pm).Error; err != nil {
			return fmt.Errorf("create resume profile: %w", mapDBError(err))
		}
		if len(fs) > 0 {
			if err := tx.WithContext(ctx).Create(&fs).Error; err != nil {
				return fmt.Errorf("create resume profile work formats: %w", mapDBError(err))
			}
		}
	}
	if resume.Overview() != nil {
		p, sections, err := mappers.OverviewToDBModels(resume.ID(), resume.Overview())
		if err != nil {
			return fmt.Errorf("create resume overview: %w", err)
		}
		if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
			return fmt.Errorf("create resume overview: %w", mapDBError(err))
		}
		if len(sections) > 0 {
			if err := tx.WithContext(ctx).Create(&sections).Error; err != nil {
				return fmt.Errorf("create resume overview sections: %w", mapDBError(err))
			}
		}
	}
	if resume.Contacts() != nil {
		if err := persistContacts(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume contacts: %w", err)
		}
	}
	if resume.Skills() != nil {
		if err := persistSkills(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume skills: %w", err)
		}
	}
	if resume.Languages() != nil {
		if err := persistLanguages(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume languages: %w", err)
		}
	}
	if resume.Technologies() != nil {
		if err := persistTechnologies(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume technologies: %w", err)
		}
	}
	if len(resume.Experiences()) > 0 {
		if err := persistExperiences(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume experiences: %w", err)
		}
	}
	if len(resume.Educations()) > 0 {
		if err := persistEducations(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume educations: %w", err)
		}
	}
	if len(resume.Courses()) > 0 {
		if err := persistCourses(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume courses: %w", err)
		}
	}
	if len(resume.Activities()) > 0 {
		if err := persistActivities(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume activities: %w", err)
		}
	}
	if len(resume.Projects()) > 0 {
		if err := persistProjects(tx.WithContext(ctx), resume); err != nil {
			return fmt.Errorf("create resume projects: %w", err)
		}
	}
	return nil
}

func (r *ResumeRepository) Save(ctx context.Context, resume *domain.Resume) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("save resume: repository dependencies are nil")
	}
	if resume == nil || resume.ID().IsZero() || resume.UserID().IsZero() || resume.Version() < 1 {
		return fmt.Errorf("save resume: invalid aggregate or expected version")
	}
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("save resume: %w", err)
	}
	result := tx.WithContext(ctx).Model(&gormmodel.Resume{}).Where("id = ? AND user_id = ? AND version = ?", resume.ID().UUID(), resume.UserID().UUID(), resume.Version()).Updates(map[string]any{"title": resume.Title(), "language": string(resume.Language()), "status": string(resume.Status()), "version": gorm.Expr("version + 1"), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return fmt.Errorf("save resume: %w", mapDBError(result.Error))
	}
	if result.RowsAffected != 1 {
		return domain.ErrConcurrencyConflict
	}
	if err := persistOverview(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistProfile(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistContacts(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistSkills(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistLanguages(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistTechnologies(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistExperiences(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistEducations(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistCourses(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistActivities(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	if err := persistProjects(tx.WithContext(ctx), resume); err != nil {
		return err
	}
	return nil
}

func (r *ResumeRepository) Delete(ctx context.Context, params ports.DeleteResumeParams) error {
	resumeID, userID, version := params.ResumeID, params.UserID, params.Version
	tx, err := database.DBFromTransactionContext(ctx)
	if err != nil {
		return fmt.Errorf("delete resume: %w", err)
	}
	res := tx.WithContext(ctx).Where("id = ? AND user_id = ? AND version = ?", resumeID.UUID(), userID.UUID(), version).Delete(&gormmodel.Resume{})
	if res.Error != nil {
		return fmt.Errorf("delete resume: %w", res.Error)
	}
	if res.RowsAffected != 1 {
		return domain.ErrConcurrencyConflict
	}
	return nil
}

func (r *ResumeRepository) loadResume(ctx context.Context, resumeID sharedvalueobjects.ID, userID sharedvalueobjects.ID) (*domain.Resume, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get resume: repository dependencies are nil")
	}
	tx := database.DBFromContext(ctx, r.db)
	var m gormmodel.Resume
	if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", resumeID.UUID(), userID.UUID()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, fmt.Errorf("get resume: %w", err)
	}
	var profileModel gormmodel.ResumeProfile
	var profile *domain.Profile
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&profileModel).Error; err == nil {
		var formats []gormmodel.ResumeProfileWorkFormat
		if err := tx.WithContext(ctx).Where("profile_id = ?", profileModel.ID).Order("position").Find(&formats).Error; err != nil {
			return nil, fmt.Errorf("get resume profile work formats: %w", err)
		}
		profile, err = mappers.ProfileToDomain(profileModel, formats)
		if err != nil {
			return nil, err
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get resume profile: %w", err)
	}
	var contactsModel gormmodel.ResumeContacts
	var contacts *domain.Contacts
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&contactsModel).Error; err == nil {
		var hidden []gormmodel.ResumeContactHidden
		if err := tx.WithContext(ctx).Where("contacts_id = ?", contactsModel.ID).Order("position").Find(&hidden).Error; err != nil {
			return nil, fmt.Errorf("get resume contact hidden: %w", err)
		}
		var textModel gormmodel.ResumeContactsText
		textErr := tx.WithContext(ctx).Where("contacts_id = ?", contactsModel.ID).First(&textModel).Error
		var text *gormmodel.ResumeContactsText
		if textErr == nil {
			text = &textModel
		} else if !errors.Is(textErr, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get resume contacts text: %w", textErr)
		}
		contacts, err = mappers.ContactsToDomain(contactsModel, hidden, text)
		if err != nil {
			return nil, err
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get resume contacts: %w", err)
	}
	skills, err := loadSkills(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	languages, err := loadLanguages(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	technologies, err := loadTechnologies(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	experiences, err := loadExperiences(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	educations, err := loadEducations(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	courses, err := loadCourses(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	activities, err := loadActivities(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	projects, err := loadProjects(tx.WithContext(ctx), resumeID)
	if err != nil {
		return nil, err
	}
	var overviewModel gormmodel.ResumeOverview
	if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).First(&overviewModel).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get resume overview: %w", err)
		}
	} else {
		var sections []gormmodel.ResumeOverviewSection
		if err := tx.WithContext(ctx).Where("resume_id = ?", m.ID).Order("position").Find(&sections).Error; err != nil {
			return nil, fmt.Errorf("get resume overview sections: %w", err)
		}
		o, err := mappers.OverviewToDomain(overviewModel, sections)
		if err != nil {
			return nil, err
		}
		snapshot, err := mappers.ResumeToDomainSnapshot(m, o, profile, contacts, skills, languages, technologies, experiences, educations, courses, activities, projects)
		if err != nil {
			return nil, err
		}
		return domain.ReconstructResume(snapshot)
	}
	snapshot, err := mappers.ResumeToDomainSnapshot(m, nil, profile, contacts, skills, languages, technologies, experiences, educations, courses, activities, projects)
	if err != nil {
		return nil, err
	}
	return domain.ReconstructResume(snapshot)
}
