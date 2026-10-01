package courses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	replace *commands.ReplaceCoursesHandler
}

func NewHandler(replace *commands.ReplaceCoursesHandler) *Handler { return &Handler{replace: replace} }

type courseRequest struct {
	ID              string  `json:"id,omitempty"`
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}
type coursesRequest struct {
	Version int64           `json:"version"`
	Items   []courseRequest `json:"items"`
}

func (r coursesRequest) input() (commands.CoursesInput, error) {
	out := make([]commands.CourseInput, 0, len(r.Items))
	for _, x := range r.Items {
		var id identity.ID
		if x.ID != "" {
			v, e := identity.ParseID(x.ID)
			if e != nil {
				return commands.CoursesInput{}, e
			}
			id = v
		}
		end, e := domain.ParseYearMonth(x.EndDate)
		if e != nil {
			return commands.CoursesInput{}, e
		}
		var link *domain.ExternalURL
		if x.CertificateLink != nil {
			v, e := domain.NewExternalURL(*x.CertificateLink)
			if e != nil {
				return commands.CoursesInput{}, e
			}
			link = &v
		}
		out = append(out, commands.CourseInput{ID: id, Name: x.Name, Organization: x.Organization, EndDate: end, Hours: x.Hours, CertificateLink: link})
	}
	return commands.CoursesInput{Items: out}, nil
}
func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body coursesRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	input, err := body.input()
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceCoursesCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
