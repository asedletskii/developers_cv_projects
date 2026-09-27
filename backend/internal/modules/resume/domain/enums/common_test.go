package enums

import "testing"

func TestEnumsAcceptContractValuesAndRejectUnknownValues(t *testing.T) {
	cases := []struct {
		name  string
		valid func() bool
		bad   func() bool
	}{
		{"work format", func() bool { return WorkFormatRemote.IsValid() }, func() bool { return WorkFormat("invalid").IsValid() }},
		{"resume language", func() bool { return ResumeLanguageEN.IsValid() }, func() bool { return ResumeLanguage("invalid").IsValid() }},
		{"spoken language", func() bool { return SpokenLanguageEN.IsValid() }, func() bool { return SpokenLanguage("invalid").IsValid() }},
		{"language grade", func() bool { return LanguageGradeC1.IsValid() }, func() bool { return LanguageGrade("invalid").IsValid() }},
		{"technology grade", func() bool { return TechnologyGradeSenior.IsValid() }, func() bool { return TechnologyGrade("invalid").IsValid() }},
		{"education level", func() bool { return EducationLevelMaster.IsValid() }, func() bool { return EducationLevel(99).IsValid() }},
		{"contact key", func() bool { return ContactKeyEmail.IsValid() }, func() bool { return ContactKey("invalid").IsValid() }},
		{"activity type", func() bool { return ActivityTypeOpenSource.IsValid() }, func() bool { return ActivityType("invalid").IsValid() }},
		{"availability", func() bool { return AvailabilityStatusOpen.IsValid() }, func() bool { return AvailabilityStatus("invalid").IsValid() }},
		{"section", func() bool { return SectionAbout.IsValid() }, func() bool { return SectionID("invalid").IsValid() }},
		{"hard skill", func() bool { return HardSkillAPI.IsValid() }, func() bool { return HardSkill("invalid").IsValid() }},
		{"soft skill", func() bool { return SoftSkillOwnership.IsValid() }, func() bool { return SoftSkill("invalid").IsValid() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.valid() || tc.bad() {
				t.Fatal("enum validation contract is broken")
			}
		})
	}
}
