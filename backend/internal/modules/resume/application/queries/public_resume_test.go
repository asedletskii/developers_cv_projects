package queries

import (
	"context"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
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
		PublicResumeQuery{ID: sharedvalueobjects.NewID()},
	)
	require.NoError(t, err)
	payload := resultJSON(t, result)
	var object map[string]any
	require.NoError(t, json.Unmarshal([]byte(payload), &object))
	require.NotContains(t, object, "id")
	require.Equal(t, view.Title, result.Title)
	require.Equal(t, view.Profile, result.Profile)
	require.Equal(t, view.Contacts, result.Contacts)
	require.Equal(t, view.Skills, result.Skills)
	require.Len(t, result.Projects, 1)
	require.Equal(t, "project-1", result.Projects[0].ID)
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

	result := toPublic(view)

	require.NotNil(t, result.Contacts)
	require.Nil(t, result.Contacts.Phone)
	require.Nil(t, result.Contacts.Telegram)
	require.Nil(t, result.Contacts.Email)
	require.Equal(t, []string{"phone", "telegram", "email"}, result.Contacts.HiddenContacts)

	// Redaction must not mutate the source view.
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

	result := toPublic(view)

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

	payload, err := json.Marshal(toPublic(view))
	require.NoError(t, err)
	require.Contains(t, string(payload), "+33123456789")
	require.Contains(t, string(payload), "developer")
	require.NotContains(t, string(payload), "private@example.com")
}

func TestPublicResumeQueryHandlerRejectsInvalidID(t *testing.T) {
	_, err := NewPublicResumeQueryHandler(publicResumeReaderStub{}).Handle(context.Background(), PublicResumeQuery{})
	require.Error(t, err)
}

func TestPublicResumeQueryHandlerPropagatesReaderError(t *testing.T) {
	want := context.Canceled
	_, err := NewPublicResumeQueryHandler(publicResumeReaderStub{err: want}).Handle(
		context.Background(),
		PublicResumeQuery{ID: sharedvalueobjects.NewID()},
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
	require.NotContains(t, object, "id")
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

func TestPublicContractContainsChildIDsButNoResumeID(t *testing.T) {
	view := ResumeView{
		ID:        "resume-1",
		Projects:  []ProjectView{{ID: "project-1"}},
		Languages: []LanguageView{{ID: "language-1", Language: "en"}},
	}

	payload, err := json.Marshal(toPublic(view))
	require.NoError(t, err)
	require.Contains(t, string(payload), `"id":"project-1"`)
	require.Contains(t, string(payload), `"id":"language-1"`)
	require.NotContains(t, string(payload), `"resumeId"`)
	require.NotContains(t, string(payload), `"id":"resume-1"`)
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

func resultJSON(t *testing.T, result PublicResumeView) string {
	t.Helper()
	payload, err := json.Marshal(result)
	require.NoError(t, err)
	return string(payload)
}
