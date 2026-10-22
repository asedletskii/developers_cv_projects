package queries

type ResumeOverviewView struct {
	Theme        *OverviewThemeSettings
	Greeting     OverviewGreeting
	Availability *OverviewAvailability
	Sections     []OverviewSectionSettings
}

type OverviewGreeting struct {
	First           string
	Second          string
	SecondHighlight string
	Third           string
}

type OverviewAvailability struct {
	Status string
	From   *string
	Note   *string
}

type OverviewThemeSettings struct {
	Primary   string
	Secondary string
}

type OverviewSectionSettings struct {
	ID        string
	Enabled   bool
	Collapsed bool
}

type ResumeSectionSettings struct {
	ID        string
	Enabled   bool
	Collapsed bool
}
