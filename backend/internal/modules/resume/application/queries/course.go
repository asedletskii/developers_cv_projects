package queries

type CourseView struct {
	ID              string  `json:"id"`
	ResumeID        string  `json:"resumeId"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}
