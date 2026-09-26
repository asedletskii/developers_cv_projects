//go:build integration

package readers_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/infrastructure/persistence/readers"
	"developers-cv/backend/internal/platform/database"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestPostgresPrivateReadersEnforceResumeOwnership(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx := context.Background()
	container, err := testpostgres.Run(ctx, "postgres:17-alpine", testpostgres.WithDatabase("developers_cv"), testpostgres.WithUsername("developers_cv"), testpostgres.WithPassword("test-password"), testpostgres.BasicWaitStrategies())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, testcontainers.TerminateContainer(container)) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	applyReaderTestMigrations(t, ctx, dsn)

	gormDB, err := database.NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, database.Close(gormDB)) })

	owner := useridentity.NewID()
	foreign := useridentity.NewID()
	resumeID := valueobjects.NewID()
	ids := insertReaderTestData(t, ctx, gormDB, owner, foreign, resumeID)

	overview := readers.NewOverviewReader(gormDB)
	profile := readers.NewProfileReader(gormDB)
	contacts := readers.NewContactsReader(gormDB)
	skills := readers.NewSkillsReader(gormDB)
	languages := readers.NewLanguagesReader(gormDB)
	technologies := readers.NewTechnologiesReader(gormDB)
	experiences := readers.NewExperiencesReader(gormDB)
	education := readers.NewEducationReader(gormDB)
	courses := readers.NewCourseReader(gormDB)
	activities := readers.NewActivityReader(gormDB)
	projects := readers.NewProjectsReader(gormDB)
	metadata := readers.NewResumeMetadataReader(gormDB)
	resumeReader := readers.NewResumeReader(readers.ResumeReaderDependencies{Metadata: metadata, Overview: overview, Profile: profile, Contacts: contacts, Skills: skills, Languages: languages, Technologies: technologies, Experiences: experiences, Educations: education, Courses: courses, Activities: activities, Projects: projects})

	_, err = resumeReader.GetResume(ctx, queries.GetResumeQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	_, err = resumeReader.GetResume(ctx, queries.GetResumeQuery{ResumeID: resumeID, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)

	_, err = profile.GetResumeProfile(ctx, queries.GetResumeProfileQuery{ResumeID: resumeID, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	_, err = contacts.GetResumeContacts(ctx, queries.GetResumeContactsQuery{ResumeID: resumeID, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	_, err = skills.GetSkills(ctx, queries.GetSkillsQuery{ResumeID: resumeID, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)

	languagesView, err := languages.ListLanguages(ctx, queries.ListLanguagesQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Len(t, languagesView, 1)
	_, err = languages.GetLanguage(ctx, queries.GetLanguageQuery{ResumeID: resumeID, LanguageID: ids.language, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	foreignLanguages, err := languages.ListLanguages(ctx, queries.ListLanguagesQuery{ResumeID: resumeID, UserID: foreign})
	require.NoError(t, err)
	require.Empty(t, foreignLanguages)

	technologiesView, err := technologies.ListTechnologies(ctx, queries.ListTechnologiesQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Len(t, technologiesView, 1)
	_, err = technologies.GetTechnology(ctx, queries.GetTechnologyQuery{ResumeID: resumeID, TechnologyID: ids.technology, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	foreignTechnologies, err := technologies.ListTechnologies(ctx, queries.ListTechnologiesQuery{ResumeID: resumeID, UserID: foreign})
	require.NoError(t, err)
	require.Empty(t, foreignTechnologies)

	experienceView, err := experiences.ListExperiences(ctx, queries.ListExperiencesQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Len(t, experienceView, 1)
	_, err = experiences.GetExperience(ctx, queries.GetExperienceQuery{ResumeID: resumeID, ExperienceID: ids.experience, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	foreignExperiences, err := experiences.ListExperiences(ctx, queries.ListExperiencesQuery{ResumeID: resumeID, UserID: foreign})
	require.NoError(t, err)
	require.Empty(t, foreignExperiences)

	educationView, err := education.ListEducations(ctx, queries.ListEducationsQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Len(t, educationView, 1)
	_, err = education.GetEducation(ctx, queries.GetEducationQuery{ResumeID: resumeID, EducationID: ids.education, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	foreignEducation, err := education.ListEducations(ctx, queries.ListEducationsQuery{ResumeID: resumeID, UserID: foreign})
	require.NoError(t, err)
	require.Empty(t, foreignEducation)

	courseView, err := courses.ListCourses(ctx, queries.ListCoursesQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Len(t, courseView, 1)
	_, err = courses.GetCourse(ctx, queries.GetCourseQuery{ResumeID: resumeID, CourseID: ids.course, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	foreignCourses, err := courses.ListCourses(ctx, queries.ListCoursesQuery{ResumeID: resumeID, UserID: foreign})
	require.NoError(t, err)
	require.Empty(t, foreignCourses)

	activityView, err := activities.ListActivities(ctx, queries.ListActivitiesQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Len(t, activityView, 1)
	_, err = activities.GetActivity(ctx, queries.GetActivityQuery{ResumeID: resumeID, ActivityID: ids.activity, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	foreignActivities, err := activities.ListActivities(ctx, queries.ListActivitiesQuery{ResumeID: resumeID, UserID: foreign})
	require.NoError(t, err)
	require.Empty(t, foreignActivities)

	projectView, err := projects.ListProjects(ctx, queries.ListProjectsQuery{ResumeID: resumeID, UserID: owner})
	require.NoError(t, err)
	require.Len(t, projectView, 1)
	_, err = projects.GetProject(ctx, queries.GetProjectQuery{ResumeID: resumeID, ProjectID: ids.project, UserID: foreign})
	require.ErrorIs(t, err, queries.ErrResumeNotFound)
	foreignProjects, err := projects.ListProjects(ctx, queries.ListProjectsQuery{ResumeID: resumeID, UserID: foreign})
	require.NoError(t, err)
	require.Empty(t, foreignProjects)
}

type readerFixtureIDs struct {
	language, technology, experience, education, course, activity, project valueobjects.ID
}

func insertReaderTestData(t *testing.T, ctx context.Context, db *gorm.DB, owner, foreign useridentity.ID, resumeID valueobjects.ID) readerFixtureIDs {
	t.Helper()
	ids := readerFixtureIDs{
		language: valueobjects.NewID(), technology: valueobjects.NewID(), experience: valueobjects.NewID(),
		education: valueobjects.NewID(), course: valueobjects.NewID(), activity: valueobjects.NewID(), project: valueobjects.NewID(),
	}
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO users (id) VALUES ($1), ($2)`, owner.UUID(), foreign.UUID()).Error)
	now := time.Now().UTC()
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resumes (id,user_id,title,created_at,updated_at) VALUES ($1,$2,'Ownership test',$3,$3)`, resumeID.UUID(), owner.UUID(), now).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_overviews (resume_id,greeting_first,greeting_second,greeting_highlight,greeting_third) VALUES ($1,'a','b','c','d')`, resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_profiles (id,resume_id,name,surname,position,photo_path,date_of_birth,location,relocation) VALUES ($1,$2,'A','B','Developer','/photo', '1998-01-01','Paris',false)`, uuid.New(), resumeID.UUID()).Error)
	contactsID := uuid.New()
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_contacts (id,resume_id,email) VALUES ($1,$2,'test@example.com')`, contactsID, resumeID.UUID()).Error)
	skillsID := uuid.New()
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_skills (id,resume_id) VALUES ($1,$2)`, skillsID, resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_skill_hard (skills_id,value,position) VALUES ($1,'Go',0)`, skillsID).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_languages (id,resume_id,language,position,grade) VALUES ($1,$2,'en',0,'B2')`, ids.language.UUID(), resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_technologies (id,resume_id,name,position,grade) VALUES ($1,$2,'Go',0,'senior')`, ids.technology.UUID(), resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_experiences (id,resume_id,sort_order,company,position,start_date) VALUES ($1,$2,0,'Acme','Developer','2020-01-01')`, ids.experience.UUID(), resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_experience_achievements (experience_id,position,value) VALUES ($1,0,'Built API')`, ids.experience.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_educations (id,resume_id,sort_order,location,university,grade,program_name,faculty) VALUES ($1,$2,0,'Paris','University',5,'Computer Science','Engineering')`, ids.education.UUID(), resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_courses (id,resume_id,sort_order,name,organization,end_date) VALUES ($1,$2,0,'Go Course','Acme','2024-01-01')`, ids.course.UUID(), resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_activities (id,resume_id,sort_order,type,title) VALUES ($1,$2,0,'openSource','Contributor')`, ids.activity.UUID(), resumeID.UUID()).Error)
	require.NoError(t, db.WithContext(ctx).Exec(`INSERT INTO resume_projects (id,resume_id,sort_order,name) VALUES ($1,$2,0,'CV')`, ids.project.UUID(), resumeID.UUID()).Error)

	return ids
}

func applyReaderTestMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()
	config, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	conn, err := pgx.ConnectConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	entries, err := os.ReadDir(filepath.Join("..", "..", "..", "..", "..", "..", "db", "migrations"))
	require.NoError(t, err)
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	for _, name := range files {
		content, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "..", "db", "migrations", name))
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(content))
		require.NoError(t, err, name)
	}
}
