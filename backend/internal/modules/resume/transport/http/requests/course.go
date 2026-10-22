package requests

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type CourseRequest struct {
	ID              string  `json:"id,omitempty"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}

type ReplaceCoursesRequest struct {
	Version int64           `json:"version"`
	Items   []CourseRequest `json:"items"`
}

func (r ReplaceCoursesRequest) ToInput() (commands.CoursesInput, error) {
	out := make([]commands.CourseInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id sharedvalueobjects.ID
		if item.ID != "" {
			parsed, err := sharedvalueobjects.ParseID(item.ID)
			if err != nil {
				return commands.CoursesInput{}, err
			}
			id = parsed
		}
		endDate, err := domain.ParseYearMonth(item.EndDate)
		if err != nil {
			return commands.CoursesInput{}, err
		}
		var link *domain.ExternalURL
		if item.CertificateLink != nil {
			parsed, err := domain.NewExternalURL(*item.CertificateLink)
			if err != nil {
				return commands.CoursesInput{}, err
			}
			link = &parsed
		}
		out = append(out, commands.CourseInput{
			ID: id, Name: item.Name, Organization: item.Organization,
			EndDate: endDate, Hours: item.Hours, CertificateLink: link,
		})
	}
	return commands.CoursesInput{Items: out}, nil
}
