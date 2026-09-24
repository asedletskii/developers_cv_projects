package request

import (
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

func (r ResumeRequest) ToCreateCommand(userID uuid.UUID) (commands.CreateResumeCommand, error) {
	userIdentity, err := useridentity.IDFromUUID(userID)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	dateOfBirth, err := valueobjects.ParseDate(r.DateOfBirth)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	photoPath, err := valueobjects.NewAssetPath(r.PhotoPath)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}

	greeting, err := newGreeting(r.Greeting)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	availability, err := newAvailability(r.Availability)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	theme, err := newTheme(r.Theme)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	sections, err := newSections(r.Sections)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	contactsText, err := newContactsText(r.ContactsText)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	contacts, err := newContacts(r.Contacts)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}

	resumeLink, err := optionalAssetPath(r.ResumeLink)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}

	languages := make([]entities.ResumeLanguageParams, 0, len(r.Languages))
	for _, item := range r.Languages {
		var grade *enums.LanguageGrade
		if item.Grade != nil {
			v := enums.LanguageGrade(*item.Grade)
			grade = &v
		}
		languages = append(languages, entities.ResumeLanguageParams{
			Language: enums.LanguageType(item.Language),
			Grade:    grade,
		})
	}

	experience := make([]entities.ExperienceParams, 0, len(r.Experience))
	for _, item := range r.Experience {
		itemID, e := optionalID(item.ID)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		start, e := valueobjects.ParseYearMonth(item.StartDate)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		end, e := optionalYearMonth(item.EndDate)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		stack, e := newTechnologies(item.Stack)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		experience = append(experience, entities.ExperienceParams{
			ID:           itemID,
			Company:      item.Company,
			Position:     item.Position,
			Location:     item.Location,
			StartDate:    start,
			EndDate:      end,
			Description:  item.Description,
			Achievements: item.Achievements,
			Stack:        stack,
		})
	}

	education := make([]entities.EducationParams, 0, len(r.Education))
	for _, item := range r.Education {
		itemID, e := optionalID(item.ID)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		end, e := optionalYear(item.EndDate)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		education = append(education, entities.EducationParams{
			ID:          itemID,
			Location:    item.Location,
			University:  item.University,
			Grade:       enums.EducationLevel(item.Grade),
			ProgramName: item.ProgramName,
			ProfileName: item.ProfileName,
			Faculty:     item.Faculty,
			Course:      item.Course,
			EndDate:     end,
		})
	}

	courses := make([]entities.CourseParams, 0, len(r.Courses))
	for _, item := range r.Courses {
		itemID, e := optionalID(item.ID)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		end, e := valueobjects.ParseYearMonth(item.EndDate)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		var link *valueobjects.ExternalURL
		if item.CertificateLink != nil {
			v, e := valueobjects.NewExternalURL(*item.CertificateLink)
			if e != nil {
				return commands.CreateResumeCommand{}, e
			}
			link = &v
		}
		courses = append(courses, entities.CourseParams{
			ID:              itemID,
			Name:            item.Name,
			Organization:    item.Organization,
			EndDate:         end,
			Hours:           item.Hours,
			CertificateLink: link,
		})
	}

	activities := make([]entities.ActivityParams, 0, len(r.Activities))
	for _, item := range r.Activities {
		itemID, e := optionalID(item.ID)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		date, e := optionalYearMonth(item.Date)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		link, e := optionalExternalURL(item.Link)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		activities = append(activities, entities.ActivityParams{
			ID:          itemID,
			Type:        enums.ActivityType(item.Type),
			Title:       item.Title,
			Description: item.Description,
			Date:        date,
			Link:        link,
		})
	}

	technologies, err := newTechnologyParams(r.Technologies)
	if err != nil {
		return commands.CreateResumeCommand{}, err
	}
	projects := make([]entities.ProjectParams, 0, len(r.Projects))
	for _, item := range r.Projects {
		itemID, e := optionalID(item.ID)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		link, e := optionalExternalURL(item.Link)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		gitLink, e := optionalExternalURL(item.GitLink)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		stack, e := newTechnologies(item.Stack)
		if e != nil {
			return commands.CreateResumeCommand{}, e
		}
		screenshots := make([]valueobjects.AssetPath, 0, len(item.Screenshots))
		for _, path := range item.Screenshots {
			v, e := valueobjects.NewAssetPath(path)
			if e != nil {
				return commands.CreateResumeCommand{}, e
			}
			screenshots = append(screenshots, v)
		}
		contributors := make([]entities.Contributor, 0, len(item.Contributors))
		for _, contributor := range item.Contributors {
			cs, e := newContacts(contributor.Contacts)
			if e != nil {
				return commands.CreateResumeCommand{}, e
			}
			v, e := entities.NewContributor(entities.ContributorParams{
				Name:     contributor.Name,
				Surname:  contributor.Surname,
				Contacts: cs,
				Role:     contributor.Role,
			})
			if e != nil {
				return commands.CreateResumeCommand{}, e
			}
			contributors = append(contributors, v)
		}
		projects = append(projects, entities.ProjectParams{
			ID:           itemID,
			Name:         item.Name,
			Description:  item.Description,
			Role:         item.Role,
			Screenshots:  screenshots,
			Link:         link,
			GitLink:      gitLink,
			Stack:        stack,
			Contributors: contributors,
		})
	}

	return commands.CreateResumeCommand{
		UserID:         userIdentity,
		Name:           r.Name,
		Surname:        r.Surname,
		Patronymic:     r.Patronymic,
		Position:       r.Position,
		Availability:   availability,
		ResumeLink:     resumeLink,
		Theme:          theme,
		Sections:       sections,
		HiddenContacts: r.HiddenContacts,
		ContactsText:   contactsText,
		Greeting:       greeting,
		PhotoPath:      photoPath,
		DateOfBirth:    dateOfBirth,
		Location:       r.Location,
		Relocation:     r.Relocation,
		WorkFormats:    r.WorkFormat,
		Contacts:       contacts,
		Languages:      languages,
		Experience:     experience,
		Education:      education,
		Courses:        courses,
		Activities:     activities,
		HardSkills:     r.HardSkills,
		SoftSkills:     r.SoftSkills,
		Technologies:   technologies,
		Projects:       projects,
	}, nil
}

func newLocalizedText(v LocalizedTextRequest) (valueobjects.LocalizedText, error) {
	return valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: v.RU, EN: v.EN})
}
func newGreeting(v GreetingRequest) (valueobjects.Greeting, error) {
	first, e := newLocalizedText(v.First)
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	second, e := newLocalizedText(v.Second)
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	highlight, e := newLocalizedText(v.SecondHighlight)
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	third, e := newLocalizedText(v.Third)
	if e != nil {
		return valueobjects.Greeting{}, e
	}
	return valueobjects.NewGreeting(valueobjects.GreetingParams{
		First:           first,
		Second:          second,
		SecondHighlight: highlight,
		Third:           third,
	}), nil
}
func newAvailability(v *AvailabilityRequest) (*valueobjects.Availability, error) {
	if v == nil {
		return nil, nil
	}
	var from *valueobjects.Date
	if v.From != nil {
		d, e := valueobjects.ParseDate(*v.From)
		if e != nil {
			return nil, e
		}
		from = &d
	}
	var note *valueobjects.LocalizedText
	if v.Note != nil {
		n, e := newLocalizedText(*v.Note)
		if e != nil {
			return nil, e
		}
		note = &n
	}
	a, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{
		Status: enums.AvailabilityStatus(v.Status),
		From:   from,
		Note:   note,
	})
	if e != nil {
		return nil, e
	}
	return &a, nil
}
func newTheme(v *ThemeRequest) (*valueobjects.ThemeSettings, error) {
	if v == nil {
		return nil, nil
	}
	var primary, secondary *valueobjects.Color
	if v.Primary != nil {
		c, e := valueobjects.NewColor(*v.Primary)
		if e != nil {
			return nil, e
		}
		primary = &c
	}
	if v.Secondary != nil {
		c, e := valueobjects.NewColor(*v.Secondary)
		if e != nil {
			return nil, e
		}
		secondary = &c
	}
	t, e := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{
		Primary:   primary,
		Secondary: secondary,
	})
	if e != nil {
		return nil, e
	}
	return &t, nil
}
func newSections(values []SectionRequest) ([]valueobjects.SectionSettings, error) {
	result := make([]valueobjects.SectionSettings, 0, len(values))
	for _, v := range values {
		s, e := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{
			ID:        enums.SectionID(v.ID),
			Enabled:   v.Enabled,
			Collapsed: v.Collapsed,
		})
		if e != nil {
			return nil, e
		}
		result = append(result, s)
	}
	return result, nil
}
func newContactsText(v *ContactsTextRequest) (*valueobjects.ContactsText, error) {
	if v == nil {
		return nil, nil
	}
	var title, subtitle *valueobjects.LocalizedText
	if v.Title != nil {
		x, e := newLocalizedText(*v.Title)
		if e != nil {
			return nil, e
		}
		title = &x
	}
	if v.Subtitle != nil {
		x, e := newLocalizedText(*v.Subtitle)
		if e != nil {
			return nil, e
		}
		subtitle = &x
	}
	t, e := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: title, Subtitle: subtitle})
	if e != nil {
		return nil, e
	}
	return &t, nil
}
func newContacts(v ContactsRequest) (entities.Contacts, error) {
	var telegram *valueobjects.Telegram
	if v.Telegram != nil {
		x, e := valueobjects.NewTelegram(*v.Telegram)
		if e != nil {
			return entities.Contacts{}, e
		}
		telegram = &x
	}
	var email *valueobjects.Email
	if v.Email != nil {
		x, e := valueobjects.NewEmail(*v.Email)
		if e != nil {
			return entities.Contacts{}, e
		}
		email = &x
	}
	return entities.NewContacts(entities.ContactsParams{
		Phone:    stringValue(v.Phone),
		Telegram: telegram,
		Email:    email,
	})
}
func newTechnologies(values []TechnologyRequest) ([]entities.Technology, error) {
	result := make([]entities.Technology, 0, len(values))
	for _, v := range values {
		t, e := entities.NewTechnology(entities.TechnologyParams{Name: v.Name})
		if e != nil {
			return nil, e
		}
		result = append(result, t)
	}
	return result, nil
}
func newTechnologyParams(values []TechnologyRequest) ([]entities.TechnologyParams, error) {
	result := make([]entities.TechnologyParams, 0, len(values))
	for _, v := range values {
		result = append(result, entities.TechnologyParams{Name: v.Name})
	}
	return result, nil
}
func optionalID(v string) (valueobjects.ID, error) {
	if v == "" {
		return valueobjects.ID{}, nil
	}
	return valueobjects.ParseID(v)
}
func optionalAssetPath(v *string) (*valueobjects.AssetPath, error) {
	if v == nil {
		return nil, nil
	}
	x, e := valueobjects.NewAssetPath(*v)
	if e != nil {
		return nil, e
	}
	return &x, nil
}
func optionalExternalURL(v *string) (*valueobjects.ExternalURL, error) {
	if v == nil {
		return nil, nil
	}
	x, e := valueobjects.NewExternalURL(*v)
	if e != nil {
		return nil, e
	}
	return &x, nil
}
func optionalYearMonth(v *string) (*valueobjects.YearMonth, error) {
	if v == nil {
		return nil, nil
	}
	x, e := valueobjects.ParseYearMonth(*v)
	if e != nil {
		return nil, e
	}
	return &x, nil
}
func optionalYear(v *string) (*valueobjects.Year, error) {
	if v == nil {
		return nil, nil
	}
	x, e := valueobjects.ParseYear(*v)
	if e != nil {
		return nil, e
	}
	return &x, nil
}
func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
