package educations

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	replace *commands.ReplaceEducationsHandler
}

func NewHandler(replace *commands.ReplaceEducationsHandler) *Handler {
	return &Handler{replace: replace}
}

type educationRequest struct {
	ID          string  `json:"id,omitempty"`
	Location    string  `json:"location"`
	University  string  `json:"university"`
	Grade       int     `json:"grade"`
	ProgramName string  `json:"programName"`
	ProfileName *string `json:"profileName,omitempty"`
	Faculty     string  `json:"faculty"`
	Course      *int    `json:"course,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}
type educationsRequest struct {
	Version int64              `json:"version"`
	Items   []educationRequest `json:"items"`
}

func (r educationsRequest) input() (commands.EducationsInput, error) {
	out := make([]commands.EducationInput, 0, len(r.Items))
	for _, x := range r.Items {
		var id valueobjects.ID
		if x.ID != "" {
			v, e := valueobjects.ParseID(x.ID)
			if e != nil {
				return commands.EducationsInput{}, e
			}
			id = v
		}
		var end *valueobjects.Year
		if x.EndDate != nil {
			v, e := valueobjects.ParseYear(*x.EndDate)
			if e != nil {
				return commands.EducationsInput{}, e
			}
			end = &v
		}
		out = append(out, commands.EducationInput{ID: id, Location: x.Location, University: x.University, Grade: enums.EducationLevel(x.Grade), ProgramName: x.ProgramName, ProfileName: x.ProfileName, Faculty: x.Faculty, Course: x.Course, EndDate: end})
	}
	return commands.EducationsInput{Items: out}, nil
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
	var body educationsRequest
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
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceEducationsCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input}); err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
