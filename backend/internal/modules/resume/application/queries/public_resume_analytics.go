package queries

type ResumeAnalyticsSection struct {
	ID   string
	Hash string
}

type ResumeAnalyticsLayout struct {
	Sections     []ResumeAnalyticsSection
	Technologies []string
}
