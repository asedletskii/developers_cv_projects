package courses

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
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
		var id valueobjects.ID
		if x.ID != "" {
			v, e := valueobjects.ParseID(x.ID)
			if e != nil {
				return commands.CoursesInput{}, e
			}
			id = v
		}
		end, e := valueobjects.ParseYearMonth(x.EndDate)
		if e != nil {
			return commands.CoursesInput{}, e
		}
		var link *valueobjects.ExternalURL
		if x.CertificateLink != nil {
			v, e := valueobjects.NewExternalURL(*x.CertificateLink)
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
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var body coursesRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if body.Version < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be positive"})
		return
	}
	input, err := body.input()
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceCoursesCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input}); err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
