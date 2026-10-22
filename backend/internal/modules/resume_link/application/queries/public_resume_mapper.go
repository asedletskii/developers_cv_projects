package queries

import resumecontracts "developers-cv/backend/internal/modules/resume/contracts"

func mapPublicResume(view resumecontracts.PublicResumeView) PublicResumeResult {
	result := PublicResumeResult{
		Title:        view.Title,
		Language:     view.Language,
		CreatedAt:    view.CreatedAt,
		UpdatedAt:    view.UpdatedAt,
		Overview:     mapPublicOverview(view.Overview),
		Languages:    mapSlice(view.Languages, mapPublicLanguage),
		Technologies: mapSlice(view.Technologies, mapPublicTechnology),
		Experiences:  mapSlice(view.Experiences, mapPublicExperience),
		Educations:   mapSlice(view.Educations, mapPublicEducation),
		Courses:      mapSlice(view.Courses, mapPublicCourse),
		Activities:   mapSlice(view.Activities, mapPublicActivity),
		Projects:     mapSlice(view.Projects, mapPublicProject),
	}
	if view.Profile != nil {
		profile := mapPublicProfile(*view.Profile)
		result.Profile = &profile
	}
	if view.Contacts != nil {
		result.Contacts = mapPublicContacts(*view.Contacts)
	}
	if view.Skills != nil {
		result.Skills = mapPublicSkills(*view.Skills)
	}
	return result
}

func mapPublicOverview(value resumecontracts.PublicOverviewView) PublicOverviewResult {
	result := PublicOverviewResult{
		Greeting: PublicGreetingResult{
			First:           value.Greeting.First,
			Second:          value.Greeting.Second,
			SecondHighlight: value.Greeting.SecondHighlight,
			Third:           value.Greeting.Third,
		},
		Sections: mapSlice(value.Sections, mapPublicSectionSettings),
	}
	if value.Theme != nil {
		result.Theme = &PublicThemeSettingsResult{Primary: value.Theme.Primary, Secondary: value.Theme.Secondary}
	}
	if value.Availability != nil {
		result.Availability = &PublicAvailabilityResult{
			Status: value.Availability.Status,
			From:   cloneStringPointer(value.Availability.From),
			Note:   cloneStringPointer(value.Availability.Note),
		}
	}
	return result
}

func mapPublicSectionSettings(value resumecontracts.PublicSectionSettings) PublicSectionSettingsResult {
	return PublicSectionSettingsResult{ID: value.ID, Enabled: value.Enabled, Collapsed: value.Collapsed}
}

func mapPublicProfile(value resumecontracts.PublicProfileView) PublicProfileResult {
	return PublicProfileResult{
		ID:          value.ID,
		Name:        value.Name,
		Surname:     value.Surname,
		Patronymic:  cloneStringPointer(value.Patronymic),
		Position:    value.Position,
		PhotoURL:    value.PhotoURL,
		DateOfBirth: value.DateOfBirth,
		Location:    value.Location,
		Relocation:  value.Relocation,
		WorkFormat:  cloneStrings(value.WorkFormat),
	}
}

func mapPublicContacts(value resumecontracts.PublicContactsView) *PublicContactsResult {
	result := &PublicContactsResult{
		ID:             value.ID,
		Phone:          cloneStringPointer(value.Phone),
		Telegram:       cloneStringPointer(value.Telegram),
		Email:          cloneStringPointer(value.Email),
		HiddenContacts: cloneStrings(value.HiddenContacts),
	}
	if value.ContactsText != nil {
		result.ContactsText = &PublicContactsTextResult{
			Title:    cloneStringPointer(value.ContactsText.Title),
			Subtitle: cloneStringPointer(value.ContactsText.Subtitle),
		}
	}
	return result
}

func mapPublicSkills(value resumecontracts.PublicSkillsView) *PublicSkillsResult {
	return &PublicSkillsResult{
		ID:         value.ID,
		HardSkills: cloneStrings(value.HardSkills),
		SoftSkills: cloneStrings(value.SoftSkills),
	}
}

func mapPublicLanguage(value resumecontracts.PublicLanguageView) PublicLanguageResult {
	return PublicLanguageResult{
		ID:       value.ID,
		Language: value.Language,
		Grade:    cloneStringPointer(value.Grade),
		Position: value.Position,
	}
}

func mapPublicTechnology(value resumecontracts.PublicTechnologyView) PublicTechnologyResult {
	return PublicTechnologyResult{ID: value.ID, Name: value.Name, Position: value.Position}
}

func mapPublicExperience(value resumecontracts.PublicExperienceView) PublicExperienceResult {
	return PublicExperienceResult{
		ID:           value.ID,
		Company:      value.Company,
		Position:     value.Position,
		Location:     cloneStringPointer(value.Location),
		StartDate:    value.StartDate,
		EndDate:      cloneStringPointer(value.EndDate),
		Description:  cloneStringPointer(value.Description),
		Achievements: cloneStrings(value.Achievements),
		Tasks:        cloneStrings(value.Tasks),
		Stack:        mapSlice(value.Stack, mapPublicExperienceStackItem),
	}
}

func mapPublicExperienceStackItem(value resumecontracts.PublicExperienceStackItemView) PublicExperienceStackItemResult {
	return PublicExperienceStackItemResult{Name: value.Name}
}

func mapPublicEducation(value resumecontracts.PublicEducationView) PublicEducationResult {
	return PublicEducationResult{
		ID:          value.ID,
		Location:    value.Location,
		University:  value.University,
		Grade:       value.Grade,
		ProgramName: value.ProgramName,
		ProfileName: cloneStringPointer(value.ProfileName),
		Faculty:     value.Faculty,
		Course:      cloneIntPointer(value.Course),
		EndDate:     cloneStringPointer(value.EndDate),
	}
}

func mapPublicCourse(value resumecontracts.PublicCourseView) PublicCourseResult {
	return PublicCourseResult{
		ID:              value.ID,
		Name:            value.Name,
		Organization:    value.Organization,
		EndDate:         value.EndDate,
		Hours:           cloneIntPointer(value.Hours),
		CertificateLink: cloneStringPointer(value.CertificateLink),
	}
}

func mapPublicActivity(value resumecontracts.PublicActivityView) PublicActivityResult {
	return PublicActivityResult{
		ID:          value.ID,
		Type:        value.Type,
		Title:       value.Title,
		Description: cloneStringPointer(value.Description),
		Date:        cloneStringPointer(value.Date),
		Link:        cloneStringPointer(value.Link),
	}
}

func mapPublicProject(value resumecontracts.PublicProjectView) PublicProjectResult {
	return PublicProjectResult{
		ID:           value.ID,
		Name:         value.Name,
		Description:  cloneStringPointer(value.Description),
		Screenshots:  cloneStrings(value.Screenshots),
		Link:         cloneStringPointer(value.Link),
		Stack:        mapSlice(value.Stack, mapPublicProjectTechnology),
		Contributors: mapSlice(value.Contributors, mapPublicProjectContributor),
	}
}

func mapPublicProjectTechnology(value resumecontracts.PublicProjectTechnologyView) PublicProjectTechnologyResult {
	return PublicProjectTechnologyResult{Name: value.Name}
}

func mapPublicProjectContributor(value resumecontracts.PublicProjectContributorView) PublicProjectContributorResult {
	return PublicProjectContributorResult{
		Name:     value.Name,
		Surname:  value.Surname,
		Role:     value.Role,
		Contacts: mapPublicProjectContacts(value.Contacts),
	}
}

func mapPublicProjectContacts(value resumecontracts.PublicProjectContactsView) PublicProjectContactsResult {
	return PublicProjectContactsResult{
		Phone:    cloneStringPointer(value.Phone),
		Telegram: cloneStringPointer(value.Telegram),
		Email:    cloneStringPointer(value.Email),
	}
}

func mapSlice[Source, Target any](values []Source, mapValue func(Source) Target) []Target {
	result := make([]Target, len(values))
	for i, value := range values {
		result[i] = mapValue(value)
	}
	return result
}

func cloneStrings(values []string) []string {
	result := make([]string, len(values))
	copy(result, values)
	return result
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}

func cloneIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}
