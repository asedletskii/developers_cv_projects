package responses

import "developers-cv/backend/internal/modules/resume_link/application/queries"

func NewPublicResumeResponse(result queries.PublicResumeResult) PublicResumeResponse {
	return PublicResumeResponse{
		Title:        result.Title,
		Language:     result.Language,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
		Overview:     newPublicOverviewResponse(result.Overview),
		Profile:      newPublicProfileResponse(result.Profile),
		Contacts:     newPublicContactsResponse(result.Contacts),
		Skills:       newPublicSkillsResponse(result.Skills),
		Languages:    mapResponseSlice(result.Languages, newPublicLanguageResponse),
		Technologies: mapResponseSlice(result.Technologies, newPublicTechnologyResponse),
		Experiences:  mapResponseSlice(result.Experiences, newPublicExperienceResponse),
		Educations:   mapResponseSlice(result.Educations, newPublicEducationResponse),
		Courses:      mapResponseSlice(result.Courses, newPublicCourseResponse),
		Activities:   mapResponseSlice(result.Activities, newPublicActivityResponse),
		Projects:     mapResponseSlice(result.Projects, newPublicProjectResponse),
	}
}

func newPublicOverviewResponse(result queries.PublicOverviewResult) PublicOverviewResponse {
	return PublicOverviewResponse{
		Theme:        newPublicThemeSettingsResponse(result.Theme),
		Greeting:     PublicGreetingResponse{First: result.Greeting.First, Second: result.Greeting.Second, SecondHighlight: result.Greeting.SecondHighlight, Third: result.Greeting.Third},
		Availability: newPublicAvailabilityResponse(result.Availability),
		Sections:     mapResponseSlice(result.Sections, newPublicSectionSettingsResponse),
	}
}

func newPublicThemeSettingsResponse(result *queries.PublicThemeSettingsResult) *PublicThemeSettingsResponse {
	if result == nil {
		return nil
	}
	return &PublicThemeSettingsResponse{Primary: result.Primary, Secondary: result.Secondary}
}

func newPublicAvailabilityResponse(result *queries.PublicAvailabilityResult) *PublicAvailabilityResponse {
	if result == nil {
		return nil
	}
	return &PublicAvailabilityResponse{Status: result.Status, From: result.From, Note: result.Note}
}

func newPublicSectionSettingsResponse(result queries.PublicSectionSettingsResult) PublicSectionSettingsResponse {
	return PublicSectionSettingsResponse{ID: result.ID, Enabled: result.Enabled, Collapsed: result.Collapsed, Style: result.Style}
}

func newPublicProfileResponse(result *queries.PublicProfileResult) *PublicProfileResponse {
	if result == nil {
		return nil
	}
	return &PublicProfileResponse{
		ID:          result.ID,
		Name:        result.Name,
		Surname:     result.Surname,
		Patronymic:  result.Patronymic,
		Position:    result.Position,
		PhotoURL:    result.PhotoURL,
		DateOfBirth: result.DateOfBirth,
		Location:    result.Location,
		Relocation:  result.Relocation,
		WorkFormat:  copyStrings(result.WorkFormat),
	}
}

func newPublicContactsResponse(result *queries.PublicContactsResult) *PublicContactsResponse {
	if result == nil {
		return nil
	}
	return &PublicContactsResponse{
		ID:             result.ID,
		Phone:          result.Phone,
		Telegram:       result.Telegram,
		Email:          result.Email,
		HiddenContacts: copyStrings(result.HiddenContacts),
		ContactsText:   newPublicContactsTextResponse(result.ContactsText),
	}
}

func newPublicContactsTextResponse(result *queries.PublicContactsTextResult) *PublicContactsTextResponse {
	if result == nil {
		return nil
	}
	return &PublicContactsTextResponse{Title: result.Title, Subtitle: result.Subtitle}
}

func newPublicSkillsResponse(result *queries.PublicSkillsResult) *PublicSkillsResponse {
	if result == nil {
		return nil
	}
	return &PublicSkillsResponse{ID: result.ID, HardSkills: copyStrings(result.HardSkills), SoftSkills: copyStrings(result.SoftSkills)}
}

func newPublicLanguageResponse(result queries.PublicLanguageResult) PublicLanguageResponse {
	return PublicLanguageResponse{ID: result.ID, Language: result.Language, Grade: result.Grade, Position: result.Position}
}

func newPublicTechnologyResponse(result queries.PublicTechnologyResult) PublicTechnologyResponse {
	return PublicTechnologyResponse{ID: result.ID, Name: result.Name, Position: result.Position}
}

func newPublicExperienceResponse(result queries.PublicExperienceResult) PublicExperienceResponse {
	return PublicExperienceResponse{
		ID:           result.ID,
		Company:      result.Company,
		Position:     result.Position,
		Location:     result.Location,
		StartDate:    result.StartDate,
		EndDate:      result.EndDate,
		Description:  result.Description,
		Achievements: copyStrings(result.Achievements),
		Tasks:        copyStrings(result.Tasks),
		Stack:        mapResponseSlice(result.Stack, newPublicExperienceStackItemResponse),
	}
}

func newPublicExperienceStackItemResponse(result queries.PublicExperienceStackItemResult) PublicExperienceStackItemResponse {
	return PublicExperienceStackItemResponse{Name: result.Name}
}

func newPublicEducationResponse(result queries.PublicEducationResult) PublicEducationResponse {
	return PublicEducationResponse{
		ID:          result.ID,
		Location:    result.Location,
		University:  result.University,
		Grade:       result.Grade,
		ProgramName: result.ProgramName,
		ProfileName: result.ProfileName,
		Faculty:     result.Faculty,
		Course:      result.Course,
		EndDate:     result.EndDate,
	}
}

func newPublicCourseResponse(result queries.PublicCourseResult) PublicCourseResponse {
	return PublicCourseResponse{
		ID:              result.ID,
		Name:            result.Name,
		Organization:    result.Organization,
		EndDate:         result.EndDate,
		Hours:           result.Hours,
		CertificateLink: result.CertificateLink,
	}
}

func newPublicActivityResponse(result queries.PublicActivityResult) PublicActivityResponse {
	return PublicActivityResponse{
		ID:          result.ID,
		Type:        result.Type,
		Title:       result.Title,
		Description: result.Description,
		Date:        result.Date,
		Link:        result.Link,
	}
}

func newPublicProjectResponse(result queries.PublicProjectResult) PublicProjectResponse {
	return PublicProjectResponse{
		ID:           result.ID,
		Name:         result.Name,
		Description:  result.Description,
		Screenshots:  copyStrings(result.Screenshots),
		Link:         result.Link,
		Stack:        mapResponseSlice(result.Stack, newPublicProjectTechnologyResponse),
		Contributors: mapResponseSlice(result.Contributors, newPublicProjectContributorResponse),
	}
}

func newPublicProjectTechnologyResponse(result queries.PublicProjectTechnologyResult) PublicProjectTechnologyResponse {
	return PublicProjectTechnologyResponse{Name: result.Name}
}

func newPublicProjectContributorResponse(result queries.PublicProjectContributorResult) PublicProjectContributorResponse {
	return PublicProjectContributorResponse{
		Name:     result.Name,
		Surname:  result.Surname,
		Role:     result.Role,
		Contacts: newPublicProjectContactsResponse(result.Contacts),
	}
}

func newPublicProjectContactsResponse(result queries.PublicProjectContactsResult) PublicProjectContactsResponse {
	return PublicProjectContactsResponse{Phone: result.Phone, Telegram: result.Telegram, Email: result.Email}
}

func mapResponseSlice[Source, Target any](values []Source, mapValue func(Source) Target) []Target {
	result := make([]Target, len(values))
	for i, value := range values {
		result[i] = mapValue(value)
	}
	return result
}

func copyStrings(values []string) []string {
	result := make([]string, len(values))
	copy(result, values)
	return result
}
