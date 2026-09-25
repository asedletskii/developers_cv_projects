package enums

type WorkFormat string

const (
	WorkFormatFullTime WorkFormat = "fullTime"
	WorkFormatPartTime WorkFormat = "partTime"
	WorkFormatOffice   WorkFormat = "office"
	WorkFormatRemote   WorkFormat = "remote"
	WorkFormatHybrid   WorkFormat = "hybrid"
)

func (v WorkFormat) IsValid() bool {
	switch v {
	case WorkFormatFullTime, WorkFormatPartTime, WorkFormatOffice, WorkFormatRemote, WorkFormatHybrid:
		return true
	default:
		return false
	}
}

type ResumeLanguage string

const (
	ResumeLanguageRU ResumeLanguage = "ru"
	ResumeLanguageEN ResumeLanguage = "en"
	ResumeLanguageFR ResumeLanguage = "fr"
	ResumeLanguageDE ResumeLanguage = "de"
	ResumeLanguageKZ ResumeLanguage = "kz"
	ResumeLanguageZH ResumeLanguage = "zh"
	ResumeLanguageES ResumeLanguage = "es"
)

func (v ResumeLanguage) IsValid() bool {
	switch v {
	case ResumeLanguageRU, ResumeLanguageEN, ResumeLanguageFR, ResumeLanguageDE, ResumeLanguageKZ, ResumeLanguageZH, ResumeLanguageES:
		return true
	default:
		return false
	}
}

type SpokenLanguage string

const (
	SpokenLanguageRU SpokenLanguage = "ru"
	SpokenLanguageEN SpokenLanguage = "en"
	SpokenLanguageFR SpokenLanguage = "fr"
	SpokenLanguageDE SpokenLanguage = "de"
	SpokenLanguageKZ SpokenLanguage = "kz"
	SpokenLanguageZH SpokenLanguage = "zh"
	SpokenLanguageES SpokenLanguage = "es"
)

func (v SpokenLanguage) IsValid() bool {
	switch v {
	case SpokenLanguageRU, SpokenLanguageEN, SpokenLanguageFR, SpokenLanguageDE, SpokenLanguageKZ, SpokenLanguageZH, SpokenLanguageES:
		return true
	default:
		return false
	}
}

type LanguageGrade string

const (
	LanguageGradeA1 LanguageGrade = "A1"
	LanguageGradeA2 LanguageGrade = "A2"
	LanguageGradeB1 LanguageGrade = "B1"
	LanguageGradeB2 LanguageGrade = "B2"
	LanguageGradeC1 LanguageGrade = "C1"
	LanguageGradeC2 LanguageGrade = "C2"
)

func (v LanguageGrade) IsValid() bool {
	switch v {
	case LanguageGradeA1, LanguageGradeA2, LanguageGradeB1, LanguageGradeB2, LanguageGradeC1, LanguageGradeC2:
		return true
	default:
		return false
	}
}

type TechnologyGrade string

const (
	TechnologyGradeIntern TechnologyGrade = "intern"
	TechnologyGradeJunior TechnologyGrade = "junior"
	TechnologyGradeMiddle TechnologyGrade = "middle"
	TechnologyGradeSenior TechnologyGrade = "senior"
)

func (v TechnologyGrade) IsValid() bool {
	switch v {
	case TechnologyGradeIntern, TechnologyGradeJunior, TechnologyGradeMiddle, TechnologyGradeSenior:
		return true
	default:
		return false
	}
}

type EducationLevel int

const (
	EducationLevelSecondary    EducationLevel = 0
	EducationLevelBachelor     EducationLevel = 1
	EducationLevelSpecialist   EducationLevel = 2
	EducationLevelMaster       EducationLevel = 3
	EducationLevelPostgraduate EducationLevel = 4
	EducationLevelCandidate    EducationLevel = 5
	EducationLevelDoctor       EducationLevel = 6
)

func (v EducationLevel) IsValid() bool {
	return v >= EducationLevelSecondary && v <= EducationLevelDoctor
}

type ContactKey string

const (
	ContactKeyTelegram ContactKey = "telegram"
	ContactKeyEmail    ContactKey = "email"
	ContactKeyPhone    ContactKey = "phone"
)

func (v ContactKey) IsValid() bool {
	switch v {
	case ContactKeyTelegram, ContactKeyEmail, ContactKeyPhone:
		return true
	default:
		return false
	}
}

type ActivityType string

const (
	ActivityTypeOpenSource   ActivityType = "openSource"
	ActivityTypeHackathon    ActivityType = "hackathon"
	ActivityTypeSpeaker      ActivityType = "speaker"
	ActivityTypeVolunteering ActivityType = "volunteering"
	ActivityTypeMentoring    ActivityType = "mentoring"
	ActivityTypeCommunity    ActivityType = "community"
)

func (v ActivityType) IsValid() bool {
	switch v {
	case ActivityTypeOpenSource, ActivityTypeHackathon, ActivityTypeSpeaker,
		ActivityTypeVolunteering, ActivityTypeMentoring, ActivityTypeCommunity:
		return true
	default:
		return false
	}
}

type AvailabilityStatus string

const (
	AvailabilityStatusOpen        AvailabilityStatus = "open"
	AvailabilityStatusConsidering AvailabilityStatus = "considering"
	AvailabilityStatusBusy        AvailabilityStatus = "busy"
)

func (v AvailabilityStatus) IsValid() bool {
	switch v {
	case AvailabilityStatusOpen, AvailabilityStatusConsidering, AvailabilityStatusBusy:
		return true
	default:
		return false
	}
}

type SectionID string

const (
	SectionAbout      SectionID = "about"
	SectionTechStack  SectionID = "techStack"
	SectionExperience SectionID = "experience"
	SectionSkills     SectionID = "skills"
	SectionEducation  SectionID = "education"
	SectionCourses    SectionID = "courses"
	SectionProjects   SectionID = "projects"
	SectionActivities SectionID = "activities"
	SectionContacts   SectionID = "contacts"
)

func (v SectionID) IsValid() bool {
	switch v {
	case SectionAbout, SectionTechStack, SectionExperience, SectionSkills, SectionEducation,
		SectionCourses, SectionProjects, SectionActivities, SectionContacts:
		return true
	default:
		return false
	}
}

type HardSkill string

const (
	HardSkillArchitecture  HardSkill = "architecture"
	HardSkillTesting       HardSkill = "testing"
	HardSkillCICD          HardSkill = "cicd"
	HardSkillPerformance   HardSkill = "performance"
	HardSkillAccessibility HardSkill = "accessibility"
	HardSkillCodeReview    HardSkill = "codeReview"
	HardSkillUIDesign      HardSkill = "uiDesign"
	HardSkillAPI           HardSkill = "api"
)

func (v HardSkill) IsValid() bool {
	switch v {
	case HardSkillArchitecture, HardSkillTesting, HardSkillCICD, HardSkillPerformance,
		HardSkillAccessibility, HardSkillCodeReview, HardSkillUIDesign, HardSkillAPI:
		return true
	default:
		return false
	}
}

type SoftSkill string

const (
	SoftSkillCommunication  SoftSkill = "communication"
	SoftSkillTeamwork       SoftSkill = "teamwork"
	SoftSkillProblemSolving SoftSkill = "problemSolving"
	SoftSkillTimeManagement SoftSkill = "timeManagement"
	SoftSkillMentoring      SoftSkill = "mentoring"
	SoftSkillAdaptability   SoftSkill = "adaptability"
	SoftSkillOwnership      SoftSkill = "ownership"
)

func (v SoftSkill) IsValid() bool {
	switch v {
	case SoftSkillCommunication, SoftSkillTeamwork, SoftSkillProblemSolving, SoftSkillTimeManagement,
		SoftSkillMentoring, SoftSkillAdaptability, SoftSkillOwnership:
		return true
	default:
		return false
	}
}
