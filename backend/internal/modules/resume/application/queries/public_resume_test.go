package queries

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	resumecontracts "developers-cv/backend/internal/modules/resume/contracts"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestPublicResumeQueryHandlerReturnsCanonicalResumeView(t *testing.T) {
	view := ResumeView{
		ID:       "resume-1",
		Title:    "Backend Developer",
		Language: "en",
		Overview: ResumeOverviewView{
			Greeting: OverviewGreeting{First: "Hello"},
		},
		Profile: &ResumeProfileView{
			Name:       "Ada",
			Surname:    "Lovelace",
			WorkFormat: []string{"remote"},
		},
		Contacts: &ResumeContactsView{HiddenContacts: []string{"email"}},
		Skills:   &SkillsView{HardSkills: []string{"Go"}, SoftSkills: []string{}},
		Projects: []ProjectView{{
			ID:           "project-1",
			Name:         "CV",
			Role:         stringPtr("Backend engineer"),
			GitLink:      stringPtr("https://github.com/example/cv"),
			Screenshots:  []string{"https://example.com/shot.png"},
			Stack:        []ProjectTechnologyView{{Name: "Go"}},
			Contributors: []ProjectContributorView{{Name: "Grace", Surname: "Hopper", Role: "Advisor"}},
		}},
	}

	result, err := NewPublicResumeQueryHandler(publicResumeReaderStub{view: view}).Handle(
		context.Background(),
		resumecontracts.PublicResumeQuery{ID: sharedvalueobjects.NewID()},
	)
	require.NoError(t, err)
	require.Equal(t, view.Title, result.Title)
	require.Equal(t, view.Profile.Name, result.Profile.Name)
	require.Equal(t, view.Profile.WorkFormat, result.Profile.WorkFormat)
	require.Equal(t, view.Contacts.HiddenContacts, result.Contacts.HiddenContacts)
	require.Equal(t, view.Skills.HardSkills, result.Skills.HardSkills)
	require.Equal(t, view.Skills.SoftSkills, result.Skills.SoftSkills)
	require.Len(t, result.Projects, 1)
	require.Equal(t, "project-1", result.Projects[0].ID)
	require.NotNil(t, result.Projects[0].Screenshots)
	require.Empty(t, result.Projects[0].Screenshots)
	require.NotNil(t, result.Projects[0].Contributors)
	require.Empty(t, result.Projects[0].Contributors)
}

func TestToPublicDoesNotMutateSourceProjects(t *testing.T) {
	role := "Backend engineer"
	gitLink := "https://github.com/example/cv"
	view := ResumeView{Projects: []ProjectView{{
		Role:         &role,
		GitLink:      &gitLink,
		Screenshots:  []string{"https://example.com/shot.png"},
		Contributors: []ProjectContributorView{{Name: "Grace"}},
	}}}

	result := ToPublicResumeView(view)

	require.NotSame(t, &view.Projects[0], &result.Projects[0])
	require.NotNil(t, view.Projects[0].Role)
	require.NotNil(t, view.Projects[0].GitLink)
	require.Len(t, view.Projects[0].Screenshots, 1)
	require.Len(t, view.Projects[0].Contributors, 1)
}

func TestToPublicRedactsEverySupportedHiddenContact(t *testing.T) {
	phone := "+33123456789"
	telegram := "developer"
	email := "private@example.com"
	view := ResumeView{
		Contacts: &ResumeContactsView{
			Phone:          &phone,
			Telegram:       &telegram,
			Email:          &email,
			HiddenContacts: []string{"phone", "telegram", "email"},
		},
	}

	result := ToPublicResumeView(view)

	require.NotNil(t, result.Contacts)
	require.Nil(t, result.Contacts.Phone)
	require.Nil(t, result.Contacts.Telegram)
	require.Nil(t, result.Contacts.Email)
	require.Equal(t, []string{"phone", "telegram", "email"}, result.Contacts.HiddenContacts)

	require.Equal(t, &phone, view.Contacts.Phone)
	require.Equal(t, &telegram, view.Contacts.Telegram)
	require.Equal(t, &email, view.Contacts.Email)
}

func TestToPublicPreservesVisibleContactValues(t *testing.T) {
	phone := "+33123456789"
	telegram := "developer"
	email := "public@example.com"
	view := ResumeView{
		Contacts: &ResumeContactsView{
			Phone:          &phone,
			Telegram:       &telegram,
			Email:          &email,
			HiddenContacts: []string{"email"},
		},
	}

	result := ToPublicResumeView(view)

	require.Equal(t, &phone, result.Contacts.Phone)
	require.Equal(t, &telegram, result.Contacts.Telegram)
	require.Nil(t, result.Contacts.Email)
}

func TestPublicProjectionDoesNotContainHiddenValues(t *testing.T) {
	view := ResumeView{
		Contacts: &ResumeContactsView{
			Phone:          stringPtr("+33123456789"),
			Telegram:       stringPtr("developer"),
			Email:          stringPtr("private@example.com"),
			HiddenContacts: []string{"email"},
		},
	}

	result := ToPublicResumeView(view)

	require.NotNil(t, result.Contacts)
	require.Equal(t, "+33123456789", *result.Contacts.Phone)
	require.Equal(t, "developer", *result.Contacts.Telegram)
	require.Nil(t, result.Contacts.Email)
}

func TestPublicResumeQueryHandlerRejectsInvalidID(t *testing.T) {
	_, err := NewPublicResumeQueryHandler(publicResumeReaderStub{}).Handle(context.Background(), resumecontracts.PublicResumeQuery{})
	require.Error(t, err)
}

func TestPublicResumeQueryHandlerPropagatesReaderError(t *testing.T) {
	want := context.Canceled
	_, err := NewPublicResumeQueryHandler(publicResumeReaderStub{err: want}).Handle(
		context.Background(),
		resumecontracts.PublicResumeQuery{ID: sharedvalueobjects.NewID()},
	)
	require.ErrorIs(t, err, want)
}

func TestPublicResumeProjectionNormalizesNilCollectionsToEmptySlices(t *testing.T) {
	view := ResumeView{
		Profile:  &ResumeProfileView{},
		Contacts: &ResumeContactsView{},
		Skills:   &SkillsView{},
		Projects: []ProjectView{{}},
	}

	result := ToPublicResumeView(view)

	require.NotNil(t, result.Languages)
	require.NotNil(t, result.Technologies)
	require.NotNil(t, result.Experiences)
	require.NotNil(t, result.Educations)
	require.NotNil(t, result.Courses)
	require.NotNil(t, result.Activities)
	require.Len(t, result.Projects, 1)
	require.NotNil(t, result.Projects[0].Screenshots)
	require.NotNil(t, result.Projects[0].Stack)
	require.NotNil(t, result.Projects[0].Contributors)
	require.NotNil(t, result.Overview.Sections)
	require.NotNil(t, result.Profile.WorkFormat)
	require.NotNil(t, result.Skills.HardSkills)
	require.NotNil(t, result.Skills.SoftSkills)
	require.NotNil(t, result.Contacts.HiddenContacts)
}

func TestPublicContractContainsChildIDsButNoResumeID(t *testing.T) {
	view := ResumeView{
		ID:        "resume-1",
		Projects:  []ProjectView{{ID: "project-1"}},
		Languages: []LanguageView{{ID: "language-1", Language: "en"}},
	}

	result := ToPublicResumeView(view)

	require.Equal(t, "project-1", result.Projects[0].ID)
	require.Equal(t, "language-1", result.Languages[0].ID)
	_, hasResumeID := reflect.TypeOf(result).FieldByName("ResumeID")
	require.False(t, hasResumeID)
}

func TestPrivateContractKeepsOnlyRootResumeID(t *testing.T) {
	view := ResumeView{
		ID:       "resume-1",
		Profile:  &ResumeProfileView{ID: "profile-1"},
		Contacts: &ResumeContactsView{ID: "contacts-1"},
		Skills:   &SkillsView{ID: "skills-1"},
		Projects: []ProjectView{{ID: "project-1"}},
	}

	payload, err := json.Marshal(view)
	require.NoError(t, err)
	require.Contains(t, string(payload), `"id":"resume-1"`)
	require.Contains(t, string(payload), `"id":"profile-1"`)
	require.Contains(t, string(payload), `"id":"contacts-1"`)
	require.Contains(t, string(payload), `"id":"skills-1"`)
	require.Contains(t, string(payload), `"id":"project-1"`)
	require.NotContains(t, string(payload), `"resumeId"`)
}

func TestPublicProjectionPreservesExplicitEmptyContactsText(t *testing.T) {
	empty := ""
	result := ToPublicResumeView(ResumeView{
		Contacts: &ResumeContactsView{
			ContactsText: &ResumeContactsTextView{Title: &empty, Subtitle: &empty},
		},
	})

	require.NotNil(t, result.Contacts)
	require.NotNil(t, result.Contacts.ContactsText)
	require.NotNil(t, result.Contacts.ContactsText.Title)
	require.NotNil(t, result.Contacts.ContactsText.Subtitle)
	require.Empty(t, *result.Contacts.ContactsText.Title)
	require.Empty(t, *result.Contacts.ContactsText.Subtitle)
}

func TestPublicProjectionKeepsOptionalFieldsNil(t *testing.T) {
	result := ToPublicResumeView(ResumeView{
		Overview:    ResumeOverviewView{Sections: nil},
		Profile:     &ResumeProfileView{},
		Contacts:    &ResumeContactsView{},
		Skills:      &SkillsView{},
		Experiences: []ExperienceView{{}},
		Educations:  []EducationView{{}},
		Courses:     []CourseView{{}},
		Activities:  []ActivityView{{}},
		Projects:    []ProjectView{{}},
	})

	require.Nil(t, result.Profile.Patronymic)
	require.Nil(t, result.Contacts.Phone)
	require.Nil(t, result.Contacts.ContactsText)
	require.Nil(t, result.Experiences[0].Location)
	require.Nil(t, result.Experiences[0].EndDate)
	require.Nil(t, result.Experiences[0].Description)
	require.Nil(t, result.Educations[0].ProfileName)
	require.Nil(t, result.Educations[0].Course)
	require.Nil(t, result.Courses[0].Hours)
	require.Nil(t, result.Courses[0].CertificateLink)
	require.Nil(t, result.Activities[0].Description)
	require.Nil(t, result.Activities[0].Date)
	require.Nil(t, result.Activities[0].Link)
	require.NotNil(t, result.Projects[0].Screenshots)
	require.Empty(t, result.Projects[0].Screenshots)
	require.NotNil(t, result.Projects[0].Contributors)
	require.Empty(t, result.Projects[0].Contributors)
}

func TestPublicDTOsAreIndependentFromPrivateDTOs(t *testing.T) {
	publicType := reflect.TypeOf(resumecontracts.PublicResumeView{})
	privateType := reflect.TypeOf(ResumeView{})
	require.NotEqual(t, publicType, privateType)
	for _, fieldName := range []string{"ID", "ResumeID", "Status", "Version"} {
		_, exists := publicType.FieldByName(fieldName)
		require.False(t, exists, "public DTO must not expose private field %s", fieldName)
	}
	publicProjectType := reflect.TypeOf(resumecontracts.PublicProjectView{})
	for _, fieldName := range []string{"Role", "GitLink"} {
		_, exists := publicProjectType.FieldByName(fieldName)
		require.False(t, exists, "public project DTO must not expose private field %s", fieldName)
	}
}
