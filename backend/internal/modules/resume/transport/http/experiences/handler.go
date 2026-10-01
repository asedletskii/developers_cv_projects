package experiences

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
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
	Tasks        []string            `json:"tasks"`
	Stack        []technologyRequest `json:"stack"`
}
type replaceRequest struct {
	Version int64               `json:"version"`
	Items   []experienceRequest `json:"items"`
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body replaceRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	items, err := experienceInputs(body.Items)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceExperiencesCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: commands.ExperiencesInput{Items: items}}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func experienceInputs(values []experienceRequest) ([]commands.ExperienceInput, error) {
	out := make([]commands.ExperienceInput, 0, len(values))
	for _, body := range values {
		var id identity.ID
		if body.ID != "" {
			parsed, err := identity.ParseID(body.ID)
			if err != nil {
				return nil, err
			}
			id = parsed
		}
		start, err := domain.ParseYearMonth(body.StartDate)
		if err != nil {
			return nil, err
		}
		var end *domain.YearMonth
		if body.EndDate != nil {
			parsed, err := domain.ParseYearMonth(*body.EndDate)
			if err != nil {
				return nil, err
			}
			end = &parsed
		}
		stack := make([]commands.ExperienceStackItemInput, len(body.Stack))
		for i, item := range body.Stack {
			var grade *domain.TechnologyGrade
			if item.Grade != nil {
				value := domain.TechnologyGrade(*item.Grade)
				grade = &value
			}
			stack[i] = commands.ExperienceStackItemInput{Name: item.Name, Grade: grade}
		}
		out = append(out, commands.ExperienceInput{ID: id, Company: body.Company, Position: body.Position, Location: body.Location, StartDate: start, EndDate: end, Description: body.Description, Achievements: append([]string(nil), body.Achievements...), Tasks: append([]string(nil), body.Tasks...), Stack: stack})
	}
	return out, nil
}
