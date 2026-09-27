package experiences

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	create *commands.CreateExperienceHandler
	list   *queries.ListExperiencesQueryHandler
	get    *queries.GetExperienceQueryHandler
	update *commands.UpdateExperienceHandler
	delete *commands.DeleteExperienceHandler
}

func NewHandler(create *commands.CreateExperienceHandler, list *queries.ListExperiencesQueryHandler, get *queries.GetExperienceQueryHandler, update *commands.UpdateExperienceHandler, delete *commands.DeleteExperienceHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

type technologyRequest struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}
type experienceRequest struct {
	Company      string              `json:"company"`
	Position     string              `json:"position"`
	Location     *string             `json:"location,omitempty"`
	StartDate    string              `json:"startDate"`
	EndDate      *string             `json:"endDate,omitempty"`
	Description  *string             `json:"description,omitempty"`
	Achievements []string            `json:"achievements"`
	Stack        []technologyRequest `json:"stack"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body experienceRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	input, err := experienceInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateExperienceCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: input})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}
func (h *Handler) List(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	items, err := h.list.Handle(c.Request.Context(), queries.ListExperiencesQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}
func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetExperienceQuery{ResumeID: rid, ExperienceID: id, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	var body experienceRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	input, err := experienceInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateExperienceCommand{ResumeID: rid, ExperienceID: id, UserID: uid, Input: input}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteExperienceCommand{ResumeID: rid, ExperienceID: id, UserID: uid}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func experienceInput(body experienceRequest) (commands.ExperienceInput, error) {
	start, err := valueobjects.ParseYearMonth(body.StartDate)
	if err != nil {
		return commands.ExperienceInput{}, err
	}
	var end *valueobjects.YearMonth
	if body.EndDate != nil {
		parsed, err := valueobjects.ParseYearMonth(*body.EndDate)
		if err != nil {
			return commands.ExperienceInput{}, err
		}
		end = &parsed
	}
	stack := make([]commands.ExperienceStackItemInput, len(body.Stack))
	for i, t := range body.Stack {
		var grade *enums.TechnologyGrade
		if t.Grade != nil {
			v := enums.TechnologyGrade(*t.Grade)
			grade = &v
		}
		stack[i] = commands.ExperienceStackItemInput{Name: t.Name, Grade: grade}
	}
	return commands.ExperienceInput{Company: body.Company, Position: body.Position, Location: body.Location, StartDate: start, EndDate: end, Description: body.Description, Achievements: body.Achievements, Stack: stack}, nil
}
func ids(c *gin.Context) (useridentity.ID, valueobjects.ID, bool) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return uid, valueobjects.ID{}, false
	}
	rid, ok := common.ParseResumeID(c)
	return uid, rid, ok
}
func itemIDs(c *gin.Context) (useridentity.ID, valueobjects.ID, valueobjects.ID, bool) {
	uid, rid, ok := ids(c)
	if !ok {
		return uid, rid, valueobjects.ID{}, false
	}
	id, err := valueobjects.ParseID(c.Param("experienceId"))
	if err != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
}
