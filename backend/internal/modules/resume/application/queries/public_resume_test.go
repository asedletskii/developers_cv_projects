package queries

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
)

type publicResumeReaderStub struct {
	view ResumeView
	err  error
}

func (s publicResumeReaderStub) GetPublic(context.Context, PublicResumeQuery) (ResumeView, error) {
	if s.err != nil {
		return ResumeView{}, s.err
	}
	return s.view, nil
}

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
		Skills:   &SkillsView{HardSkills: []string{"Go"}},
		Projects: []ProjectView{{
			ID:           "project-1",
			ResumeID:     "resume-1",
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
		PublicResumeQuery{ID: valueobjects.NewID()},
	)
	require.NoError(t, err)
	require.Equal(t, view.ID, result.ID)
	require.Nil(t, result.Version)
	require.Equal(t, view.Title, result.Title)
	require.Equal(t, view.Profile, result.Profile)
	require.Equal(t, view.Contacts, result.Contacts)
	require.Equal(t, view.Skills, result.Skills)
	require.Len(t, result.Projects, 1)
	require.Equal(t, "project-1", result.Projects[0].ID)
	require.Equal(t, "resume-1", result.Projects[0].ResumeID)
	require.Nil(t, result.Projects[0].Role)
	require.Nil(t, result.Projects[0].GitLink)
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

	result := toPublic(view)

	require.NotSame(t, &view.Projects[0], &result.Projects[0])
	require.NotNil(t, view.Projects[0].Role)
	require.NotNil(t, view.Projects[0].GitLink)
	require.Len(t, view.Projects[0].Screenshots, 1)
	require.Len(t, view.Projects[0].Contributors, 1)
}

func TestToPublicRedactsHiddenContactValuesWithoutMutatingSource(t *testing.T) {
	phone := "+33123456789"
	telegram := "developer"
	email := "user@example.com"
	view := ResumeView{
		Contacts: &ResumeContactsView{
			Phone:          &phone,
			Telegram:       &telegram,
			Email:          &email,
			HiddenContacts: []string{"phone", "email"},
		},
	}

	result := toPublic(view)

	require.NotNil(t, result.Contacts)
	require.Nil(t, result.Contacts.Phone)
	require.Equal(t, &telegram, result.Contacts.Telegram)
	require.Nil(t, result.Contacts.Email)
	require.Equal(t, []string{"phone", "email"}, result.Contacts.HiddenContacts)
	require.NotNil(t, view.Contacts.Phone)
	require.NotNil(t, view.Contacts.Email)
	require.Equal(t, &phone, view.Contacts.Phone)
	require.Equal(t, &email, view.Contacts.Email)
}

func TestToPublicRedactsEverySupportedHiddenContact(t *testing.T) {
	phone := "+33123456789"
	telegram := "developer"
	email := "user@example.com"
	view := ResumeView{
		Contacts: &ResumeContactsView{
			Phone:          &phone,
			Telegram:       &telegram,
			Email:          &email,
			HiddenContacts: []string{"phone", "telegram", "email"},
		},
	}

	result := toPublic(view)

	require.Nil(t, result.Contacts.Phone)
	require.Nil(t, result.Contacts.Telegram)
	require.Nil(t, result.Contacts.Email)

	payload, err := json.Marshal(result)
	require.NoError(t, err)
	require.NotContains(t, string(payload), phone)
	require.NotContains(t, string(payload), telegram)
	require.NotContains(t, string(payload), email)
	require.Contains(t, string(payload), `"hiddenContacts":["phone","telegram","email"]`)
}

func TestPublicResumeQueryHandlerRejectsInvalidID(t *testing.T) {
	_, err := NewPublicResumeQueryHandler(publicResumeReaderStub{}).Handle(context.Background(), PublicResumeQuery{})
	require.Error(t, err)
}

func TestPublicResumeQueryHandlerPropagatesReaderError(t *testing.T) {
	want := context.Canceled
	_, err := NewPublicResumeQueryHandler(publicResumeReaderStub{err: want}).Handle(
		context.Background(),
		PublicResumeQuery{ID: valueobjects.NewID()},
	)
	require.ErrorIs(t, err, want)
}

func TestPublicResumeJSONUsesCanonicalNestedShape(t *testing.T) {
	view := ResumeView{
		ID:           "resume-1",
		Title:        "CV",
		Language:     "en",
		Overview:     ResumeOverviewView{Sections: []OverviewSectionSettings{}},
		Languages:    []LanguageView{},
		Technologies: []TechnologyView{},
		Experiences:  []ExperienceView{},
		Educations:   []EducationView{},
		Courses:      []CourseView{},
		Activities:   []ActivityView{},
		Projects:     []ProjectView{},
	}

	payload, err := json.Marshal(toPublic(view))
	require.NoError(t, err)

	var object map[string]any
	require.NoError(t, json.Unmarshal(payload, &object))
	require.Contains(t, object, "overview")
	require.NotContains(t, object, "version")
	require.Contains(t, object, "experiences")
	require.Contains(t, object, "educations")
	require.NotContains(t, object, "name")
	require.NotContains(t, object, "surname")
	require.NotContains(t, object, "experience")
	require.NotContains(t, object, "education")
	require.NotContains(t, object, "presentation")
	require.NotContains(t, object, "hardSkills")
	require.NotContains(t, object, "softSkills")
}
