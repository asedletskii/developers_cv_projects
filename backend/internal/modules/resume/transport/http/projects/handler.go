package projects

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct {
	create *commands.CreateProjectHandler
	list   *queries.ListProjectsQueryHandler
	get    *queries.GetProjectQueryHandler
	update *commands.UpdateProjectHandler
	delete *commands.DeleteProjectHandler
}

func NewHandler(create *commands.CreateProjectHandler, list *queries.ListProjectsQueryHandler, get *queries.GetProjectQueryHandler, update *commands.UpdateProjectHandler, delete *commands.DeleteProjectHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, delete: delete}
}

type TechnologyRequest struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type ContributorRequest struct {
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	Role     string          `json:"role"`
	Contacts ContactsRequest `json:"contacts"`
}

type ContactsRequest struct {
	Phone    string  `json:"phone"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type ProjectRequest struct {
	Name         string               `json:"name"`
	Description  *string              `json:"description,omitempty"`
	Role         *string              `json:"role,omitempty"`
	Screenshots  []string             `json:"screenshots"`
	Link         *string              `json:"link,omitempty"`
	GitLink      *string              `json:"gitLink,omitempty"`
	Stack        []TechnologyRequest  `json:"stack"`
	Contributors []ContributorRequest `json:"contributors"`
}

type TechnologyResponse struct {
	Name  string  `json:"name"`
	Grade *string `json:"grade,omitempty"`
}

type ContributorResponse struct {
	Name     string           `json:"name"`
	Surname  string           `json:"surname"`
	Role     string           `json:"role"`
	Contacts ContactsResponse `json:"contacts"`
}

type ContactsResponse struct {
	Phone    *string `json:"phone,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type ProjectResponse struct {
	ID           string                `json:"id"`
	ResumeID     string                `json:"resumeId"`
	Name         string                `json:"name"`
	Description  *string               `json:"description,omitempty"`
	Role         *string               `json:"role,omitempty"`
	Screenshots  []string              `json:"screenshots"`
	Link         *string               `json:"link,omitempty"`
	GitLink      *string               `json:"gitLink,omitempty"`
	Stack        []TechnologyResponse  `json:"stack"`
	Contributors []ContributorResponse `json:"contributors"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body ProjectRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input := projectInput(body)
	id, err := h.create.Handle(c.Request.Context(), commands.CreateProjectCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: input})
	if err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetProjectQuery{ResumeID: rid, ProjectID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusCreated, projectQueryView(item))
}

func (h *Handler) List(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	items, err := h.list.Handle(c.Request.Context(), queries.ListProjectsQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	out := make([]ProjectResponse, len(items))
	for i, item := range items {
		out[i] = projectQueryView(item)
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetProjectQuery{ResumeID: rid, ProjectID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusOK, projectQueryView(item))
}

func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	var body ProjectRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	input := projectInput(body)
	if err := h.update.Handle(c.Request.Context(), commands.UpdateProjectCommand{ResumeID: rid, ProjectID: id, UserID: uid, Input: input}); err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), queries.GetProjectQuery{ResumeID: rid, ProjectID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusOK, projectQueryView(item))
}

func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteProjectCommand{ResumeID: rid, ProjectID: id, UserID: uid}); err != nil {
		write(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func projectInput(body ProjectRequest) commands.ProjectInput {
	stack := make([]commands.ProjectTechnologyInput, len(body.Stack))
	for i, v := range body.Stack {
		stack[i] = commands.ProjectTechnologyInput{Name: v.Name, Grade: v.Grade}
	}
	contributors := make([]commands.ProjectContributorInput, len(body.Contributors))
	for i, v := range body.Contributors {
		contributors[i] = commands.ProjectContributorInput{Name: v.Name, Surname: v.Surname, Role: v.Role, Contacts: commands.ProjectContactsInput{Phone: v.Contacts.Phone, Telegram: v.Contacts.Telegram, Email: v.Contacts.Email}}
	}
	return commands.ProjectInput{Name: body.Name, Description: body.Description, Role: body.Role, Screenshots: body.Screenshots, Link: body.Link, GitLink: body.GitLink, Stack: stack, Contributors: contributors}
}

func projectQueryView(item queries.ProjectView) ProjectResponse {
	stack := make([]TechnologyResponse, len(item.Stack))
	for i, technology := range item.Stack {
		stack[i] = TechnologyResponse{Name: technology.Name, Grade: technology.Grade}
	}
	contributors := make([]ContributorResponse, len(item.Contributors))
	for i, contributor := range item.Contributors {
		contributors[i] = ContributorResponse{
			Name: contributor.Name, Surname: contributor.Surname, Role: contributor.Role,
			Contacts: ContactsResponse{Phone: contributor.Contacts.Phone, Telegram: contributor.Contacts.Telegram, Email: contributor.Contacts.Email},
		}
	}
	return ProjectResponse{
		ID: item.ID, ResumeID: item.ResumeID, Name: item.Name, Description: item.Description, Role: item.Role,
		Screenshots: item.Screenshots, Link: item.Link, GitLink: item.GitLink, Stack: stack, Contributors: contributors,
	}
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
	id, err := valueobjects.ParseID(c.Param("projectId"))
	if err != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
}

func write(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrEntityNotFound) {
		common.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume resource not found"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateEntity) {
		common.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume resource already exists"))
		return
	}
	common.WriteCommandError(c, err)
}
