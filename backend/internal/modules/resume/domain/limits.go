package domain

import "unicode/utf8"

const (
	MaxProfileNameLength      = 100
	MaxCompanyLength          = 120
	MaxProfilePositionLength  = 120
	MaxLocationLength         = 200
	MaxDescriptionLength      = 2000
	MaxShortTextLength        = 200
	MaxListItemTextLength     = 500
	MaxSkillLength            = 60
	MaxSkillsPerResume        = 50
	MaxExperiencesPerResume   = 30
	MaxExperienceTasks        = 30
	MaxExperienceAchievements = 30
	MaxExperienceStack        = 50
	MaxProjectsPerResume      = 30
	MaxProjectScreenshots     = 10
	MaxProjectContributors    = 20
	MaxProjectStack           = 50
	MaxEducationsPerResume    = 30
	MaxCoursesPerResume       = 30
	MaxActivitiesPerResume    = 30
	MaxTechnologiesPerResume  = 50
	MaxLanguagesPerResume     = 7
	MaxContactPhoneLength     = 32
	MaxContactTelegramLength  = 64
	MaxContactEmailLength     = 254
)

func textLengthWithin(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func optionalTextLengthWithin(value *string, max int) bool {
	return value == nil || textLengthWithin(*value, max)
}
