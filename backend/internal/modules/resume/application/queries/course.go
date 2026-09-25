package queries

import (
	"context"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"fmt"
)

type ListCoursesQuery struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type GetCourseQuery struct {
	ResumeID, CourseID valueobjects.ID
	UserID             useridentity.ID
}
type CourseView struct {
	ID              string  `json:"id"`
	ResumeID        string  `json:"resumeId"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}
type CoursesReader interface {
	ListCourses(context.Context, ListCoursesQuery) ([]CourseView, error)
	GetCourse(context.Context, GetCourseQuery) (CourseView, error)
}
type ListCoursesQueryHandler struct{ reader CoursesReader }

func NewListCoursesQueryHandler(reader CoursesReader) *ListCoursesQueryHandler {
	return &ListCoursesQueryHandler{reader: reader}
}
func (h *ListCoursesQueryHandler) Handle(ctx context.Context, q ListCoursesQuery) ([]CourseView, error) {
	if h == nil || h.reader == nil {
		return nil, fmt.Errorf("list courses query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.UserID.IsZero() {
		return nil, fmt.Errorf("invalid list courses query")
	}
	return h.reader.ListCourses(ctx, q)
}

type GetCourseQueryHandler struct{ reader CoursesReader }

func NewGetCourseQueryHandler(reader CoursesReader) *GetCourseQueryHandler {
	return &GetCourseQueryHandler{reader: reader}
}
func (h *GetCourseQueryHandler) Handle(ctx context.Context, q GetCourseQuery) (CourseView, error) {
	if h == nil || h.reader == nil {
		return CourseView{}, fmt.Errorf("get course query handler dependencies are nil")
	}
	if q.ResumeID.IsZero() || q.CourseID.IsZero() || q.UserID.IsZero() {
		return CourseView{}, fmt.Errorf("invalid get course query")
	}
	return h.reader.GetCourse(ctx, q)
}
