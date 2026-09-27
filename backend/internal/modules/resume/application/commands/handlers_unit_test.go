package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/ports"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type fakeCommandUnitOfWork struct {
	called bool
	err    error
}

func (u *fakeCommandUnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	if u.err != nil {
		return u.err
	}
	return fn(ctx)
}

type fakeCommandRepository struct {
	called  bool
	err     error
	create  any
	update  any
	delete  any
	reorder any
}

func (r *fakeCommandRepository) CreateActivity(_ context.Context, cmd ports.CreateActivityRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateActivity(_ context.Context, cmd ports.UpdateActivityRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) DeleteActivity(_ context.Context, cmd ports.DeleteActivityRepositoryCommand) error {
	r.called, r.delete = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateContacts(_ context.Context, cmd ports.CreateResumeContactsRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateContacts(_ context.Context, cmd ports.UpdateResumeContactsRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateCourse(_ context.Context, cmd ports.CreateCourseRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateCourse(_ context.Context, cmd ports.UpdateCourseRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) DeleteCourse(_ context.Context, cmd ports.DeleteCourseRepositoryCommand) error {
	r.called, r.delete = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateEducation(_ context.Context, cmd ports.CreateEducationRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateEducation(_ context.Context, cmd ports.UpdateEducationRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) DeleteEducation(_ context.Context, cmd ports.DeleteEducationRepositoryCommand) error {
	r.called, r.delete = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateExperience(_ context.Context, cmd ports.CreateExperienceRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateExperience(_ context.Context, cmd ports.UpdateExperienceRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) DeleteExperience(_ context.Context, cmd ports.DeleteExperienceRepositoryCommand) error {
	r.called, r.delete = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateLanguage(_ context.Context, cmd ports.CreateLanguageRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateLanguage(_ context.Context, cmd ports.UpdateLanguageRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) ReorderLanguages(_ context.Context, cmd ports.ReorderLanguagesRepositoryCommand) error {
	r.called, r.reorder = true, cmd
	return r.err
}
func (r *fakeCommandRepository) DeleteLanguage(_ context.Context, cmd ports.DeleteLanguageRepositoryCommand) error {
	r.called, r.delete = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateOverview(_ context.Context, cmd ports.UpdateOverviewRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateProfile(_ context.Context, cmd ports.CreateResumeProfileRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateProfile(_ context.Context, cmd ports.UpdateResumeProfileRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateProject(_ context.Context, cmd ports.CreateProjectRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateProject(_ context.Context, cmd ports.UpdateProjectRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) DeleteProject(_ context.Context, cmd ports.DeleteProjectRepositoryCommand) error {
	r.called, r.delete = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateSkills(_ context.Context, cmd ports.CreateSkillsRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateSkills(_ context.Context, cmd ports.UpdateSkillsRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) CreateTechnology(_ context.Context, cmd ports.CreateTechnologyRepositoryCommand) error {
	r.called, r.create = true, cmd
	return r.err
}
func (r *fakeCommandRepository) UpdateTechnology(_ context.Context, cmd ports.UpdateTechnologyRepositoryCommand) error {
	r.called, r.update = true, cmd
	return r.err
}
func (r *fakeCommandRepository) ReorderTechnologies(_ context.Context, cmd ports.ReorderTechnologiesRepositoryCommand) error {
	r.called, r.reorder = true, cmd
	return r.err
}
func (r *fakeCommandRepository) DeleteTechnology(_ context.Context, cmd ports.DeleteTechnologyRepositoryCommand) error {
	r.called, r.delete = true, cmd
	return r.err
}

func testDate(t *testing.T) valueobjects.Date {
	t.Helper()
	v, err := valueobjects.ParseDate("2020-01-02")
	require.NoError(t, err)
	return v
}

func testYearMonth(t *testing.T) valueobjects.YearMonth {
	t.Helper()
	v, err := valueobjects.ParseYearMonth("2020-01")
	require.NoError(t, err)
	return v
}

func testURL(t *testing.T) valueobjects.ExternalURL {
	t.Helper()
	v, err := valueobjects.NewExternalURL("https://example.com/item")
	require.NoError(t, err)
	return v
}

func testEmail(t *testing.T) valueobjects.Email {
	t.Helper()
	v, err := valueobjects.NewEmail("user@example.com")
	require.NoError(t, err)
	return v
}

func testTelegram(t *testing.T) valueobjects.Telegram {
	t.Helper()
	v, err := valueobjects.NewTelegram("developer")
	require.NoError(t, err)
	return v
}

func mustCommand[T any](t *testing.T, value any) T {
	t.Helper()
	command, ok := value.(T)
	require.True(t, ok)
	return command
}

func TestActivityHandlersCreateUpdateDelete(t *testing.T) {
	ctx := context.Background()
	resumeID, userID := valueobjects.NewID(), useridentity.NewID()
	repo, uow := &fakeCommandRepository{}, &fakeCommandUnitOfWork{}
	created, err := NewCreateActivityHandler(repo, uow).Handle(ctx, CreateActivityCommand{ResumeID: resumeID, UserID: userID, Input: ActivityInput{Type: enums.ActivityTypeOpenSource, Title: "  OSS  "}})
	require.NoError(t, err)
	require.False(t, created.IsZero())
	require.True(t, repo.called)
	require.True(t, uow.called)
	require.Equal(t, resumeID, mustCommand[ports.CreateActivityRepositoryCommand](t, repo.create).ResumeID)

	repo.called = false
	activityID := valueobjects.NewID()
	err = NewUpdateActivityHandler(repo, uow).Handle(ctx, UpdateActivityCommand{ResumeID: resumeID, ActivityID: activityID, UserID: userID, Input: ActivityInput{Type: enums.ActivityTypeSpeaker, Title: "Talk"}})
	require.NoError(t, err)
	require.Equal(t, activityID, mustCommand[ports.UpdateActivityRepositoryCommand](t, repo.update).Activity.ID())

	repo.called = false
	err = NewDeleteActivityHandler(repo, uow).Handle(ctx, DeleteActivityCommand{ResumeID: resumeID, ActivityID: activityID, UserID: userID})
	require.NoError(t, err)
	require.Equal(t, activityID, mustCommand[ports.DeleteActivityRepositoryCommand](t, repo.delete).ActivityID)
}

func TestCourseEducationAndExperienceHandlers(t *testing.T) {
	ctx := context.Background()
	resumeID, userID := valueobjects.NewID(), useridentity.NewID()
	uow := &fakeCommandUnitOfWork{}

	courseRepo := &fakeCommandRepository{}
	courseID, err := NewCreateCourseHandler(courseRepo, uow).Handle(ctx, CreateCourseCommand{ResumeID: resumeID, UserID: userID, Input: CourseInput{Name: " Go", Organization: "Org", EndDate: testYearMonth(t)}})
	require.NoError(t, err)
	require.Equal(t, courseID, mustCommand[ports.CreateCourseRepositoryCommand](t, courseRepo.create).Course.ID())
	err = NewUpdateCourseHandler(courseRepo, uow).Handle(ctx, UpdateCourseCommand{ResumeID: resumeID, CourseID: courseID, UserID: userID, Input: CourseInput{Name: "Updated", Organization: "Org", EndDate: testYearMonth(t)}})
	require.NoError(t, err)
	err = NewDeleteCourseHandler(courseRepo, uow).Handle(ctx, DeleteCourseCommand{ResumeID: resumeID, CourseID: courseID, UserID: userID})
	require.NoError(t, err)

	educationRepo := &fakeCommandRepository{}
	educationID, err := NewCreateEducationHandler(educationRepo, uow).Handle(ctx, CreateEducationCommand{ResumeID: resumeID, UserID: userID, Input: EducationInput{Location: "Moscow", University: "University", Grade: enums.EducationLevelBachelor, ProgramName: "CS", Faculty: "IT"}})
	require.NoError(t, err)
	require.Equal(t, educationID, mustCommand[ports.CreateEducationRepositoryCommand](t, educationRepo.create).Education.ID())
	err = NewUpdateEducationHandler(educationRepo, uow).Handle(ctx, UpdateEducationCommand{ResumeID: resumeID, EducationID: educationID, UserID: userID, Input: EducationInput{Location: "SPb", University: "University", Grade: enums.EducationLevelMaster, ProgramName: "SE", Faculty: "IT"}})
	require.NoError(t, err)
	err = NewDeleteEducationHandler(educationRepo, uow).Handle(ctx, DeleteEducationCommand{ResumeID: resumeID, EducationID: educationID, UserID: userID})
	require.NoError(t, err)

	experienceRepo := &fakeCommandRepository{}
	start := testYearMonth(t)
	experienceID, err := NewCreateExperienceHandler(experienceRepo, uow).Handle(ctx, CreateExperienceCommand{ResumeID: resumeID, UserID: userID, Input: ExperienceInput{Company: "Company", Position: "Go Developer", StartDate: start, Achievements: []string{" shipped"}, Stack: []ExperienceStackItemInput{{Name: " Go ", Grade: ptr(enums.TechnologyGradeSenior)}}}})
	require.NoError(t, err)
	require.Equal(t, experienceID, mustCommand[ports.CreateExperienceRepositoryCommand](t, experienceRepo.create).Experience.ID())
	err = NewUpdateExperienceHandler(experienceRepo, uow).Handle(ctx, UpdateExperienceCommand{ResumeID: resumeID, ExperienceID: experienceID, UserID: userID, Input: ExperienceInput{Company: "Company 2", Position: "Staff", StartDate: start}})
	require.NoError(t, err)
	err = NewDeleteExperienceHandler(experienceRepo, uow).Handle(ctx, DeleteExperienceCommand{ResumeID: resumeID, ExperienceID: experienceID, UserID: userID})
	require.NoError(t, err)
}

func TestLanguageAndTechnologyHandlers(t *testing.T) {
	ctx := context.Background()
	resumeID, userID := valueobjects.NewID(), useridentity.NewID()
	uow := &fakeCommandUnitOfWork{}

	languageRepo := &fakeCommandRepository{}
	languageID, err := NewCreateLanguageHandler(languageRepo, uow).Handle(ctx, CreateLanguageCommand{ResumeID: resumeID, UserID: userID, Input: LanguageInput{Language: enums.SpokenLanguageEN, Grade: ptr(enums.LanguageGradeC1)}})
	require.NoError(t, err)
	err = NewUpdateLanguageHandler(languageRepo, uow).Handle(ctx, UpdateLanguageCommand{ResumeID: resumeID, LanguageID: languageID, UserID: userID, Input: LanguageInput{Language: enums.SpokenLanguageFR}})
	require.NoError(t, err)
	order := []valueobjects.ID{languageID, valueobjects.NewID()}
	err = NewReorderLanguagesHandler(languageRepo, uow).Handle(ctx, ReorderLanguagesCommand{ResumeID: resumeID, UserID: userID, LanguageIDs: order})
	require.NoError(t, err)
	require.Equal(t, order, mustCommand[ports.ReorderLanguagesRepositoryCommand](t, languageRepo.reorder).LanguageIDs)
	err = NewDeleteLanguageHandler(languageRepo, uow).Handle(ctx, DeleteLanguageCommand{ResumeID: resumeID, LanguageID: languageID, UserID: userID})
	require.NoError(t, err)

	technologyRepo := &fakeCommandRepository{}
	technologyID, err := NewCreateTechnologyHandler(technologyRepo, uow).Handle(ctx, CreateTechnologyCommand{ResumeID: resumeID, UserID: userID, Input: TechnologyInput{Name: " Go ", Grade: ptr(enums.TechnologyGradeSenior)}})
	require.NoError(t, err)
	err = NewUpdateTechnologyHandler(technologyRepo, uow).Handle(ctx, UpdateTechnologyCommand{ResumeID: resumeID, TechnologyID: technologyID, UserID: userID, Input: TechnologyInput{Name: "Rust"}})
	require.NoError(t, err)
	techOrder := []valueobjects.ID{technologyID, valueobjects.NewID()}
	err = NewReorderTechnologiesHandler(technologyRepo, uow).Handle(ctx, ReorderTechnologiesCommand{ResumeID: resumeID, UserID: userID, TechnologyIDs: techOrder})
	require.NoError(t, err)
	require.Equal(t, techOrder, mustCommand[ports.ReorderTechnologiesRepositoryCommand](t, technologyRepo.reorder).TechnologyIDs)
	err = NewDeleteTechnologyHandler(technologyRepo, uow).Handle(ctx, DeleteTechnologyCommand{ResumeID: resumeID, TechnologyID: technologyID, UserID: userID})
	require.NoError(t, err)
}

func TestProfileContactsSkillsAndOverviewHandlers(t *testing.T) {
	ctx := context.Background()
	resumeID, userID := valueobjects.NewID(), useridentity.NewID()
	uow := &fakeCommandUnitOfWork{}
	photo := testURL(t)
	birth := testDate(t)

	profileRepo := &fakeCommandRepository{}
	profileID, err := NewCreateResumeProfileHandler(profileRepo, uow).Handle(ctx, CreateResumeProfileCommand{ResumeID: resumeID, UserID: userID, Profile: ProfileInput{Name: "John", Surname: "Doe", Position: "Developer", PhotoPath: photo, DateOfBirth: birth, Location: "Paris", WorkFormats: []enums.WorkFormat{enums.WorkFormatRemote}}})
	require.NoError(t, err)
	require.Equal(t, profileID, mustCommand[ports.CreateResumeProfileRepositoryCommand](t, profileRepo.create).Profile.ID())
	err = NewUpdateResumeProfileHandler(profileRepo, uow).Handle(ctx, UpdateResumeProfileCommand{ResumeID: resumeID, UserID: userID, Profile: ProfileInput{Name: "Jane", Surname: "Doe", Position: "Senior Developer", PhotoPath: photo, DateOfBirth: birth, Location: "Paris"}})
	require.NoError(t, err)

	email := testEmail(t)
	telegram := testTelegram(t)
	contactsRepo := &fakeCommandRepository{}
	_, err = NewCreateResumeContactsHandler(contactsRepo, uow).Handle(ctx, CreateResumeContactsCommand{ResumeID: resumeID, UserID: userID, Contacts: ContactsInput{Phone: ptr("+33123456789"), Telegram: &telegram, Email: &email, HiddenContacts: []enums.ContactKey{enums.ContactKeyPhone}}})
	require.NoError(t, err)
	err = NewUpdateResumeContactsHandler(contactsRepo, uow).Handle(ctx, UpdateResumeContactsCommand{ResumeID: resumeID, UserID: userID, Contacts: ContactsInput{Email: &email}})
	require.NoError(t, err)

	skillsRepo := &fakeCommandRepository{}
	_, err = NewCreateSkillsHandler(skillsRepo, uow).Handle(ctx, CreateSkillsCommand{ResumeID: resumeID, UserID: userID, Hard: []enums.HardSkill{enums.HardSkillAPI}, Soft: []enums.SoftSkill{enums.SoftSkillMentoring}})
	require.NoError(t, err)
	err = NewUpdateSkillsHandler(skillsRepo, uow).Handle(ctx, UpdateSkillsCommand{ResumeID: resumeID, UserID: userID, Hard: []enums.HardSkill{enums.HardSkillTesting}, Soft: []enums.SoftSkill{enums.SoftSkillTeamwork}})
	require.NoError(t, err)

	color, err := valueobjects.NewColor("#fff")
	require.NoError(t, err)
	overviewRepo := &fakeCommandRepository{}
	err = NewUpdateOverviewHandler(overviewRepo, uow).Handle(ctx, UpdateOverviewCommand{ResumeID: resumeID, UserID: userID, Input: OverviewInput{Theme: &ThemeInput{Primary: &color}, Greeting: GreetingInput{First: "Hello"}, Sections: []SectionInput{{ID: enums.SectionAbout, Enabled: true}}}})
	require.NoError(t, err)
	overviewCommand := mustCommand[ports.UpdateOverviewRepositoryCommand](t, overviewRepo.update)
	overview := overviewCommand.Overview
	require.Equal(t, resumeID, overviewCommand.ResumeID)
	require.Equal(t, enums.SectionAbout, overview.Sections()[0].ID())
}

func TestProjectHandlers(t *testing.T) {
	ctx := context.Background()
	resumeID, userID := valueobjects.NewID(), useridentity.NewID()
	uow := &fakeCommandUnitOfWork{}
	email := testEmail(t)
	repo := &fakeCommandRepository{}
	projectID, err := NewCreateProjectHandler(repo, uow).Handle(ctx, CreateProjectCommand{ResumeID: resumeID, UserID: userID, Input: ProjectInput{Name: "Project", Stack: []ProjectTechnologyInput{{Name: "Go"}}, Contributors: []ProjectContributorInput{{Name: "John", Surname: "Doe", Role: "Developer", Contacts: ProjectContactsInput{Email: &email}}}}})
	require.NoError(t, err)
	createdCommand := mustCommand[ports.CreateProjectRepositoryCommand](t, repo.create)
	created := createdCommand.Project
	require.Equal(t, projectID, created.ID())
	require.Len(t, created.Stack(), 1)
	require.Len(t, created.Contributors(), 1)

	err = NewUpdateProjectHandler(repo, uow).Handle(ctx, UpdateProjectCommand{ResumeID: resumeID, ProjectID: projectID, UserID: userID, Input: ProjectInput{Name: "Updated"}})
	require.NoError(t, err)
	require.Equal(t, projectID, mustCommand[ports.UpdateProjectRepositoryCommand](t, repo.update).Project.ID())
	err = NewDeleteProjectHandler(repo, uow).Handle(ctx, DeleteProjectCommand{ResumeID: resumeID, ProjectID: projectID, UserID: userID})
	require.NoError(t, err)
}

func TestProjectHandlersRejectContributorWithoutContactsBeforeRepository(t *testing.T) {
	ctx := context.Background()
	repo := &fakeCommandRepository{}
	uow := &fakeCommandUnitOfWork{}

	_, err := NewCreateProjectHandler(repo, uow).Handle(ctx, CreateProjectCommand{
		ResumeID: valueobjects.NewID(),
		UserID:   useridentity.NewID(),
		Input: ProjectInput{
			Name: "Project",
			Contributors: []ProjectContributorInput{{
				Name: "John", Surname: "Doe", Role: "Developer",
			}},
		},
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "at least one contact")
	require.False(t, repo.called)
}

func TestCommandHandlersPropagateRepositoryErrors(t *testing.T) {
	repositoryErr := errors.New("repository failure")
	repo := &fakeCommandRepository{err: repositoryErr}
	uow := &fakeCommandUnitOfWork{}
	_, err := NewCreateActivityHandler(repo, uow).Handle(context.Background(), CreateActivityCommand{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID(), Input: ActivityInput{Type: enums.ActivityTypeSpeaker, Title: "Talk"}})
	require.ErrorIs(t, err, repositoryErr)
	require.True(t, repo.called)

	repo = &fakeCommandRepository{err: repositoryErr}
	err = NewUpdateTechnologyHandler(repo, uow).Handle(context.Background(), UpdateTechnologyCommand{ResumeID: valueobjects.NewID(), TechnologyID: valueobjects.NewID(), UserID: useridentity.NewID(), Input: TechnologyInput{Name: "Go"}})
	require.ErrorIs(t, err, repositoryErr)
}

func TestCommandHandlersStopBeforeRepositoryOnInvalidDomainInput(t *testing.T) {
	repo := &fakeCommandRepository{}
	uow := &fakeCommandUnitOfWork{}
	_, err := NewCreateActivityHandler(repo, uow).Handle(context.Background(), CreateActivityCommand{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID(), Input: ActivityInput{Title: "missing type"}})
	require.Error(t, err)
	require.False(t, repo.called)

	repo = &fakeCommandRepository{}
	_, err = NewCreateExperienceHandler(repo, uow).Handle(context.Background(), CreateExperienceCommand{ResumeID: valueobjects.NewID(), UserID: useridentity.NewID(), Input: ExperienceInput{Company: "Company", Position: "Developer", StartDate: testYearMonth(t), Achievements: []string{"   "}}})
	require.Error(t, err)
	require.False(t, repo.called)
}

func TestUpdateHandlersRejectMissingIdentity(t *testing.T) {
	uow := &fakeCommandUnitOfWork{}
	repo := &fakeCommandRepository{}

	err := NewUpdateResumeProfileHandler(repo, uow).Handle(context.Background(), UpdateResumeProfileCommand{
		UserID: useridentity.NewID(),
		Profile: ProfileInput{
			Name:        "John",
			Surname:     "Doe",
			Position:    "Developer",
			PhotoPath:   testURL(t),
			DateOfBirth: testDate(t),
			Location:    "Paris",
		},
	})
	require.Error(t, err)
	require.False(t, repo.called)
	require.False(t, uow.called)
}

func ptr[T any](value T) *T {
	return &value
}

func TestDeleteHandlersRejectInvalidIDs(t *testing.T) {
	ctx := context.Background()
	validResumeID := valueobjects.NewID()
	validEntityID := valueobjects.NewID()
	validUserID := useridentity.NewID()
	tests := []struct {
		name   string
		handle func(DeleteCommandIDs) error
	}{
		{
			name: "activity",
			handle: func(c DeleteCommandIDs) error {
				return NewDeleteActivityHandler(&fakeCommandRepository{}, &fakeCommandUnitOfWork{}).Handle(ctx, DeleteActivityCommand{ResumeID: c.resumeID, ActivityID: c.entityID, UserID: c.userID})
			},
		},
		{
			name: "course",
			handle: func(c DeleteCommandIDs) error {
				return NewDeleteCourseHandler(&fakeCommandRepository{}, &fakeCommandUnitOfWork{}).Handle(ctx, DeleteCourseCommand{ResumeID: c.resumeID, CourseID: c.entityID, UserID: c.userID})
			},
		},
		{
			name: "education",
			handle: func(c DeleteCommandIDs) error {
				return NewDeleteEducationHandler(&fakeCommandRepository{}, &fakeCommandUnitOfWork{}).Handle(ctx, DeleteEducationCommand{ResumeID: c.resumeID, EducationID: c.entityID, UserID: c.userID})
			},
		},
		{
			name: "experience",
			handle: func(c DeleteCommandIDs) error {
				return NewDeleteExperienceHandler(&fakeCommandRepository{}, &fakeCommandUnitOfWork{}).Handle(ctx, DeleteExperienceCommand{ResumeID: c.resumeID, ExperienceID: c.entityID, UserID: c.userID})
			},
		},
		{
			name: "language",
			handle: func(c DeleteCommandIDs) error {
				return NewDeleteLanguageHandler(&fakeCommandRepository{}, &fakeCommandUnitOfWork{}).Handle(ctx, DeleteLanguageCommand{ResumeID: c.resumeID, LanguageID: c.entityID, UserID: c.userID})
			},
		},
		{
			name: "technology",
			handle: func(c DeleteCommandIDs) error {
				return NewDeleteTechnologyHandler(&fakeCommandRepository{}, &fakeCommandUnitOfWork{}).Handle(ctx, DeleteTechnologyCommand{ResumeID: c.resumeID, TechnologyID: c.entityID, UserID: c.userID})
			},
		},
		{
			name: "project",
			handle: func(c DeleteCommandIDs) error {
				return NewDeleteProjectHandler(&fakeCommandRepository{}, &fakeCommandUnitOfWork{}).Handle(ctx, DeleteProjectCommand{ResumeID: c.resumeID, ProjectID: c.entityID, UserID: c.userID})
			},
		},
	}

	invalidCases := []struct {
		name string
		ids  DeleteCommandIDs
	}{
		{name: "zero resume id", ids: DeleteCommandIDs{resumeID: valueobjects.ID{}, entityID: validEntityID, userID: validUserID}},
		{name: "zero entity id", ids: DeleteCommandIDs{resumeID: validResumeID, entityID: valueobjects.ID{}, userID: validUserID}},
		{name: "zero user id", ids: DeleteCommandIDs{resumeID: validResumeID, entityID: validEntityID, userID: useridentity.ID{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, invalid := range invalidCases {
				t.Run(invalid.name, func(t *testing.T) {
					err := tt.handle(invalid.ids)
					require.ErrorIs(t, err, domain.ErrInvalidEntity)
				})
			}
		})
	}
}

type DeleteCommandIDs struct {
	resumeID valueobjects.ID
	entityID valueobjects.ID
	userID   useridentity.ID
}
