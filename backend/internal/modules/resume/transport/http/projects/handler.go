package projects

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct {
	create *commands.CreateProjectHandler
	list   *commands.ListProjectsHandler
	get    *commands.GetProjectHandler
	update *commands.UpdateProjectHandler
	delete *commands.DeleteProjectHandler
}

func NewHandler(create *commands.CreateProjectHandler, list *commands.ListProjectsHandler, get *commands.GetProjectHandler, update *commands.UpdateProjectHandler, delete *commands.DeleteProjectHandler) *Handler {
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
	input, err := projectInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateProjectCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: input})
	if err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), commands.GetProjectQuery{ResumeID: rid, ProjectID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusCreated, projectView(rid, item))
}

func (h *Handler) List(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	items, err := h.list.Handle(c.Request.Context(), commands.ListProjectsQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	out := make([]ProjectResponse, len(items))
	for i, item := range items {
		out[i] = projectView(rid, item)
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c)
	if !ok {
		return
	}
	item, err := h.get.Handle(c.Request.Context(), commands.GetProjectQuery{ResumeID: rid, ProjectID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusOK, projectView(rid, item))
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
	input, err := projectInput(body)
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateProjectCommand{ResumeID: rid, ProjectID: id, UserID: uid, Input: input}); err != nil {
		write(c, err)
		return
	}
	item, err := h.get.Handle(c.Request.Context(), commands.GetProjectQuery{ResumeID: rid, ProjectID: id, UserID: uid})
	if err != nil {
		write(c, err)
		return
	}
	c.JSON(http.StatusOK, projectView(rid, item))
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

func projectInput(body ProjectRequest) (commands.ProjectInput, error) {
	stack := make([]entities.Technology, 0, len(body.Stack))
	for _, item := range body.Stack {
		var grade *enums.TechnologyGrade
		if item.Grade != nil {
			value := enums.TechnologyGrade(*item.Grade)
			grade = &value
		}
		technology, err := entities.NewTechnology(entities.TechnologyParams{Name: item.Name, Grade: grade})
		if err != nil {
			return commands.ProjectInput{}, err
		}
		stack = append(stack, technology)
	}

	screenshots := make([]valueobjects.AssetPath, 0, len(body.Screenshots))
	for _, path := range body.Screenshots {
		value, err := valueobjects.NewAssetPath(path)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		screenshots = append(screenshots, value)
	}

	var link, gitLink *valueobjects.ExternalURL
	if body.Link != nil {
		value, err := valueobjects.NewExternalURL(*body.Link)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		link = &value
	}
	if body.GitLink != nil {
		value, err := valueobjects.NewExternalURL(*body.GitLink)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		gitLink = &value
	}

	contributors := make([]entities.Contributor, 0, len(body.Contributors))
	for _, item := range body.Contributors {
		contacts, err := contributorContacts(item.Contacts)
		if err != nil {
			return commands.ProjectInput{}, err
		}
		contributor, err := entities.NewContributor(entities.ContributorParams{Name: item.Name, Surname: item.Surname, Role: item.Role, Contacts: contacts})
		if err != nil {
			return commands.ProjectInput{}, err
		}
		contributors = append(contributors, contributor)
	}

	return commands.ProjectInput{Name: body.Name, Description: body.Description, Role: body.Role, Screenshots: screenshots, Link: link, GitLink: gitLink, Stack: stack, Contributors: contributors}, nil
}

func contributorContacts(input ContactsRequest) (entities.Contacts, error) {
	var telegram *valueobjects.Telegram
	if input.Telegram != nil {
		value, err := valueobjects.NewTelegram(*input.Telegram)
		if err != nil {
			return entities.Contacts{}, err
		}
		telegram = &value
	}
	var email *valueobjects.Email
	if input.Email != nil {
		value, err := valueobjects.NewEmail(*input.Email)
		if err != nil {
			return entities.Contacts{}, err
		}
		email = &value
	}
	if input.Phone == "" && telegram == nil && email == nil {
		return entities.EmptyContacts(), nil
	}
	return entities.NewContacts(entities.ContactsParams{Phone: input.Phone, Telegram: telegram, Email: email})
}

func projectView(resumeID valueobjects.ID, item *domain.Project) ProjectResponse {
	stack := item.Stack()
	stackView := make([]TechnologyResponse, len(stack))
	for i, technology := range stack {
		var grade *string
		if value := technology.Grade(); value != nil {
			text := string(*value)
			grade = &text
		}
		stackView[i] = TechnologyResponse{Name: technology.Name(), Grade: grade}
	}
	contributors := item.Contributors()
	contributorsView := make([]ContributorResponse, len(contributors))
	for i, contributor := range contributors {
		contacts := contributor.Contacts()
		contributorsView[i] = ContributorResponse{Name: contributor.Name(), Surname: contributor.Surname(), Role: contributor.Role(), Contacts: ContactsResponse{Phone: contacts.Phone(), Telegram: telegramString(contacts.Telegram()), Email: emailString(contacts.Email())}}
	}
	screenshots := item.Screenshots()
	screenshotView := make([]string, len(screenshots))
	for i, screenshot := range screenshots {
		screenshotView[i] = screenshot.String()
	}
	return ProjectResponse{ID: item.ID().String(), ResumeID: resumeID.String(), Name: item.Name(), Description: item.Description(), Role: item.Role(), Screenshots: screenshotView, Link: externalURLString(item.Link()), GitLink: externalURLString(item.GitLink()), Stack: stackView, Contributors: contributorsView}
}

func telegramString(value *valueobjects.Telegram) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}
func emailString(value *valueobjects.Email) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}
func externalURLString(value *valueobjects.ExternalURL) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
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
