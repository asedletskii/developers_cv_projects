package queries

type ResumeOverviewView struct {
	Theme        *OverviewThemeSettings    `json:"theme,omitempty"`
	Greeting     OverviewGreeting          `json:"greeting"`
	Availability *OverviewAvailability     `json:"availability,omitempty"`
	ResumeLink   *string                   `json:"resumeLink,omitempty"`
	Sections     []OverviewSectionSettings `json:"sections"`
}

type OverviewGreeting struct {
	First           string `json:"first"`
	Second          string `json:"second"`
	SecondHighlight string `json:"secondHighlight"`
	Third           string `json:"third"`
}

type OverviewAvailability struct {
	Status string  `json:"status"`
	From   *string `json:"from,omitempty"`
	Note   *string `json:"note,omitempty"`
}

type OverviewThemeSettings struct {
	Primary   string `json:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty"`
}

type OverviewSectionSettings struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Collapsed bool   `json:"collapsed"`
}
