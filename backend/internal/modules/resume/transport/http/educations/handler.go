package educations

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	create *commands.CreateEducationHandler
	list   *queries.ListEducationsQueryHandler
	get    *queries.GetEducationQueryHandler
	update *commands.UpdateEducationHandler
	delete *commands.DeleteEducationHandler
}

func NewHandler(create *commands.CreateEducationHandler, list *queries.ListEducationsQueryHandler, get *queries.GetEducationQueryHandler, update *commands.UpdateEducationHandler, delete *commands.DeleteEducationHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

type EducationRequest struct {
	Location    string  `json:"location"`
	University  string  `json:"university"`
	Grade       int     `json:"grade"`
	ProgramName string  `json:"programName"`
	ProfileName *string `json:"profileName,omitempty"`
	Faculty     string  `json:"faculty"`
	Course      *int    `json:"course,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body EducationRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input, err := educationInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateEducationCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: input})
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
	items, err := h.list.Handle(c.Request.Context(), queries.ListEducationsQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}
func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetEducationQuery{ResumeID: rid, EducationID: id, UserID: uid})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	var body EducationRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input, err := educationInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateEducationCommand{ResumeID: rid, EducationID: id, UserID: uid, Input: input}); err != nil {
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
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteEducationCommand{ResumeID: rid, EducationID: id, UserID: uid}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func educationInput(body EducationRequest) (commands.EducationInput, error) {
	var end *valueobjects.Year
	if body.EndDate != nil {
		v, err := valueobjects.ParseYear(*body.EndDate)
		if err != nil {
			return commands.EducationInput{}, err
		}
		end = &v
	}
	return commands.EducationInput{Location: body.Location, University: body.University, Grade: enums.EducationLevel(body.Grade), ProgramName: body.ProgramName, ProfileName: body.ProfileName, Faculty: body.Faculty, Course: body.Course, EndDate: end}, nil
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
	id, err := valueobjects.ParseID(c.Param("educationId"))
	if err != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
}
