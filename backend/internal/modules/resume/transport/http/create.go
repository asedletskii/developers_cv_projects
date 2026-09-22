package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type createResumeRequest struct {
	Name           string              `json:"name"`
	Surname        string              `json:"surname"`
	Patronymic     *string             `json:"patronymic,omitempty"`
	Position       string              `json:"position"`
	Availability   *createAvailability `json:"availability,omitempty"`
	ResumeLink     *string             `json:"resumeLink,omitempty"`
	Theme          *createTheme        `json:"theme,omitempty"`
	Sections       []createSection     `json:"sections"`
	HiddenContacts []string            `json:"hiddenContacts"`
	ContactsText   *createContactsText `json:"contactsText,omitempty"`
	Greeting       createGreeting      `json:"greeting"`
	PhotoPath      string              `json:"photoPath"`
	DateOfBirth    string              `json:"dateOfBirth"`
	Location       string              `json:"location"`
	Relocation     bool                `json:"relocation"`
	WorkFormat     []string            `json:"workFormat"`
	Contacts       createContacts      `json:"contacts"`
	Languages      []createLanguage    `json:"languages"`
	Experience     []createExperience  `json:"experience"`
	Education      []createEducation   `json:"education"`
	Courses        []createCourse      `json:"courses"`
	Activities     []createActivity    `json:"activities"`
	HardSkills     []string            `json:"hardSkills"`
	SoftSkills     []string            `json:"softSkills"`
	Technologies   []createTechnology  `json:"technologies"`
	Projects       []createProject     `json:"projects"`
}

type createLocalizedText struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}
type createGreeting struct {
	First           createLocalizedText `json:"first"`
	Second          createLocalizedText `json:"second"`
	SecondHighlight createLocalizedText `json:"secondHighlight"`
	Third           createLocalizedText `json:"third"`
}
type createAvailability struct {
	Status string               `json:"status"`
	From   *string              `json:"from,omitempty"`
	Note   *createLocalizedText `json:"note,omitempty"`
}
type createTheme struct {
	Primary   *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
}
type createSection struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}
type createContactsText struct {
	Title    *createLocalizedText `json:"title,omitempty"`
	Subtitle *createLocalizedText `json:"subtitle,omitempty"`
}
type createContacts struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}
type createTechnology struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}
type createLanguage struct {
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
}
type createExperience struct {
	ID           string             `json:"id,omitempty"`
	Company      string             `json:"company"`
	Position     string             `json:"position"`
	Location     *string            `json:"location,omitempty"`
	StartDate    string             `json:"startDate"`
	EndDate      *string            `json:"endDate,omitempty"`
	Description  *string            `json:"description,omitempty"`
	Achievements []string           `json:"achievements"`
	Stack        []createTechnology `json:"stack"`
}
type createEducation struct {
	ID          string  `json:"id,omitempty"`
	Location    string  `json:"location"`
	University  string  `json:"university"`
	Grade       int     `json:"grade"`
	ProgramName string  `json:"programName"`
	ProfileName *string `json:"profileName,omitempty"`
	Faculty     string  `json:"faculty"`
	Course      *int    `json:"course,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}
type createCourse struct {
	ID              string  `json:"id,omitempty"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}
type createActivity struct {
	ID          string  `json:"id,omitempty"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}
type createContributor struct {
	Name     string         `json:"name"`
	Surname  string         `json:"surname"`
	Role     string         `json:"role"`
	Contacts createContacts `json:"contacts"`
}
type createProject struct {
	ID           string              `json:"id,omitempty"`
	Name         string              `json:"name"`
	Description  *string             `json:"description,omitempty"`
	Role         *string             `json:"role,omitempty"`
	Screenshots  []string            `json:"screenshots"`
	Link         *string             `json:"link,omitempty"`
	GitLink      *string             `json:"gitLink,omitempty"`
	Stack        []createTechnology  `json:"stack"`
	Contributors []createContributor `json:"contributors"`
}

func (h *Handler) CreateResume(c *gin.Context) {
	var request createResumeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, apperrors.New(
			apperrors.CodeBadRequest,
			http.StatusBadRequest,
			"Invalid request body",
		))
		return
	}

	command, err := request.toCommand()
	if err != nil {
		if isValidationError(err) {
			writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, err.Error()))
			return
		}
		writeError(c, apperrors.Internal(err))
		return
	}

	if err := h.createResume.Handle(c.Request.Context(), command); err != nil {
		if isValidationError(err) {
			writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, err.Error()))
			return
		}
		writeError(c, apperrors.Internal(err))
		return
	}

	c.Status(http.StatusCreated)
}

func (r createResumeRequest) toCommand() (commands.CreateResumeCommand, error) {
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

func newLocalizedText(v createLocalizedText) (valueobjects.LocalizedText, error) {
	return valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: v.RU, EN: v.EN})
}
func newGreeting(v createGreeting) (valueobjects.Greeting, error) {
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
func newAvailability(v *createAvailability) (*valueobjects.Availability, error) {
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
func newTheme(v *createTheme) (*valueobjects.ThemeSettings, error) {
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
func newSections(values []createSection) ([]valueobjects.SectionSettings, error) {
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
func newContactsText(v *createContactsText) (*valueobjects.ContactsText, error) {
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
func newContacts(v createContacts) (entities.Contacts, error) {
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
func newTechnologies(values []createTechnology) ([]entities.Technology, error) {
	result := make([]entities.Technology, 0, len(values))
	for _, v := range values {
		var grade *enums.TechnologyGrade
		if v.Grade != nil {
			x := enums.TechnologyGrade(*v.Grade)
			grade = &x
		}
		t, e := entities.NewTechnology(entities.TechnologyParams{Name: v.Name, Grade: grade})
		if e != nil {
			return nil, e
		}
		result = append(result, t)
	}
	return result, nil
}
func newTechnologyParams(values []createTechnology) ([]entities.TechnologyParams, error) {
	result := make([]entities.TechnologyParams, 0, len(values))
	for _, v := range values {
		var grade *enums.TechnologyGrade
		if v.Grade != nil {
			x := enums.TechnologyGrade(*v.Grade)
			grade = &x
		}
		result = append(result, entities.TechnologyParams{Name: v.Name, Grade: grade})
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
func isValidationError(err error) bool { return err != nil && (domainErr(err) || valueErr(err)) }
func domainErr(err error) bool {
	return fmt.Sprintf("%v", err) != "" && (containsError(err, domain.ErrInvalidResume) ||
		containsError(err, domain.ErrInvalidEntity) || containsError(err, domain.ErrResumeInvariant) ||
		containsError(err, domain.ErrDuplicateEntity))
}
func valueErr(err error) bool { return containsError(err, valueobjects.ErrInvalidValue) }
func containsError(err, target error) bool {
	for err != nil {
		if err == target {
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := err.(unwrapper)
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}
