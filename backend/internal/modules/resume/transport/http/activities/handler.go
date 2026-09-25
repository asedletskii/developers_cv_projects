package activities

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct {
	create *commands.CreateActivityHandler
	list   *queries.ListActivitiesQueryHandler
	get    *queries.GetActivityQueryHandler
	update *commands.UpdateActivityHandler
	delete *commands.DeleteActivityHandler
}

func NewHandler(create *commands.CreateActivityHandler, list *queries.ListActivitiesQueryHandler, get *queries.GetActivityQueryHandler, update *commands.UpdateActivityHandler, delete *commands.DeleteActivityHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

type ActivityRequest struct {
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body ActivityRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input, err := activityInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateActivityCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: input})
	if err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetActivityQuery{ResumeID: rid, ActivityID: id, UserID: uid})
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
	items, err := h.list.Handle(c.Request.Context(), queries.ListActivitiesQuery{ResumeID: rid, UserID: uid})
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
	item, err := h.get.Handle(c.Request.Context(), queries.GetActivityQuery{ResumeID: rid, ActivityID: id, UserID: uid})
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
	var body ActivityRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input, err := activityInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateActivityCommand{ResumeID: rid, ActivityID: id, UserID: uid, Input: input}); err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetActivityQuery{ResumeID: rid, ActivityID: id, UserID: uid})
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
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteActivityCommand{ResumeID: rid, ActivityID: id, UserID: uid}); err != nil {
		write(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func activityInput(body ActivityRequest) (commands.ActivityInput, error) {
	var date *valueobjects.YearMonth
	if body.Date != nil {
		v, err := valueobjects.ParseYearMonth(*body.Date)
		if err != nil {
			return commands.ActivityInput{}, err
		}
		date = &v
	}
	var link *valueobjects.ExternalURL
	if body.Link != nil {
		v, err := valueobjects.NewExternalURL(*body.Link)
		if err != nil {
			return commands.ActivityInput{}, err
		}
		link = &v
	}
	return commands.ActivityInput{Type: enums.ActivityType(body.Type), Title: body.Title, Description: body.Description, Date: date, Link: link}, nil
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
	id, err := valueobjects.ParseID(c.Param("activityId"))
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
