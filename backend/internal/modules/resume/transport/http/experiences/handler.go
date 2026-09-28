package experiences

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	replace *commands.ReplaceExperiencesHandler
}

func NewHandler(replace *commands.ReplaceExperiencesHandler) *Handler {
	return &Handler{replace: replace}
}

type technologyRequest struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}
type experienceRequest struct {
	ID           string              `json:"id,omitempty"`
	Company      string              `json:"company"`
	Position     string              `json:"position"`
	Location     *string             `json:"location,omitempty"`
	StartDate    string              `json:"startDate"`
	EndDate      *string             `json:"endDate,omitempty"`
	Description  *string             `json:"description,omitempty"`
	Achievements []string            `json:"achievements"`
	Stack        []technologyRequest `json:"stack"`
}
type replaceRequest struct {
	Version int64               `json:"version"`
	Items   []experienceRequest `json:"items"`
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body replaceRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if body.Version < 1 {
		common.WriteInvalidBody(c)
		return
	}
	items, err := experienceInputs(body.Items)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceExperiencesCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: commands.ExperiencesInput{Items: items}}); err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func experienceInputs(values []experienceRequest) ([]commands.ExperienceInput, error) {
	out := make([]commands.ExperienceInput, 0, len(values))
	for _, body := range values {
		var id valueobjects.ID
		if body.ID != "" {
			parsed, err := valueobjects.ParseID(body.ID)
			if err != nil {
				return nil, err
			}
			id = parsed
		}
		start, err := valueobjects.ParseYearMonth(body.StartDate)
		if err != nil {
			return nil, err
		}
		var end *valueobjects.YearMonth
		if body.EndDate != nil {
			parsed, err := valueobjects.ParseYearMonth(*body.EndDate)
			if err != nil {
				return nil, err
			}
			end = &parsed
		}
		stack := make([]commands.ExperienceStackItemInput, len(body.Stack))
		for i, item := range body.Stack {
			var grade *enums.TechnologyGrade
			if item.Grade != nil {
				value := enums.TechnologyGrade(*item.Grade)
				grade = &value
			}
			stack[i] = commands.ExperienceStackItemInput{Name: item.Name, Grade: grade}
		}
		out = append(out, commands.ExperienceInput{ID: id, Company: body.Company, Position: body.Position, Location: body.Location, StartDate: start, EndDate: end, Description: body.Description, Achievements: append([]string(nil), body.Achievements...), Stack: stack})
	}
	return out, nil
}

func ids(c *gin.Context) (useridentity.ID, valueobjects.ID, bool) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return uid, valueobjects.ID{}, false
	}
	rid, ok := common.ParseResumeID(c)
	return uid, rid, ok
}
