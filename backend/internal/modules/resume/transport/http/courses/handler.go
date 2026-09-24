package courses

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	create *commands.CreateCourseHandler
	list   *queries.ListCoursesQueryHandler
	get    *queries.GetCourseQueryHandler
	update *commands.UpdateCourseHandler
	delete *commands.DeleteCourseHandler
}

func NewHandler(create *commands.CreateCourseHandler, list *queries.ListCoursesQueryHandler, get *queries.GetCourseQueryHandler, update *commands.UpdateCourseHandler, delete *commands.DeleteCourseHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

type CourseRequest struct {
	Name            string  `json:"name"`
	Organization    string  `json:"organization"`
	EndDate         string  `json:"endDate"`
	Hours           *int    `json:"hours,omitempty"`
	CertificateLink *string `json:"certificateLink,omitempty"`
}
type CourseResponse = queries.CourseView

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body CourseRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input, err := courseInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateCourseCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: input})
	if err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetCourseQuery{ResumeID: rid, CourseID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}
func (h *Handler) List(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	items, err := h.list.Handle(c.Request.Context(), queries.ListCoursesQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}
func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetCourseQuery{ResumeID: rid, CourseID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	var body CourseRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input, err := courseInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateCourseCommand{ResumeID: rid, CourseID: id, UserID: uid, Input: input}); err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetCourseQuery{ResumeID: rid, CourseID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteCourseCommand{ResumeID: rid, CourseID: id, UserID: uid}); err != nil {
		write(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func courseInput(body CourseRequest) (commands.CourseInput, error) {
	end, err := valueobjects.ParseYearMonth(body.EndDate)
	if err != nil {
		return commands.CourseInput{}, err
	}
	var link *valueobjects.ExternalURL
	if body.CertificateLink != nil {
		v, err := valueobjects.NewExternalURL(*body.CertificateLink)
		if err != nil {
			return commands.CourseInput{}, err
		}
		link = &v
	}
	return commands.CourseInput{Name: body.Name, Organization: body.Organization, EndDate: end, Hours: body.Hours, CertificateLink: link}, nil
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
	id, err := valueobjects.ParseID(c.Param("courseId"))
	if err != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
}
func write(c *gin.Context, err error) {
	if errors.Is(err, queries.ErrResumeNotFound) || errors.Is(err, domain.ErrEntityNotFound) {
		common.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume resource not found"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateEntity) {
		common.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume resource already exists"))
		return
	}
	common.WriteCommandError(c, err)
}
