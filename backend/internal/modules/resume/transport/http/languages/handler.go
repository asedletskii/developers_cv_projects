package languages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	create  *commands.CreateLanguageHandler
	get     *queries.GetLanguageQueryHandler
	list    *queries.ListLanguagesQueryHandler
	update  *commands.UpdateLanguageHandler
	reorder *commands.ReorderLanguagesHandler
	delete  *commands.DeleteLanguageHandler
}

func NewHandler(create *commands.CreateLanguageHandler, list *queries.ListLanguagesQueryHandler, get *queries.GetLanguageQueryHandler, update *commands.UpdateLanguageHandler, reorder *commands.ReorderLanguagesHandler, delete *commands.DeleteLanguageHandler) *Handler {
	return &Handler{create: create, list: list, get: get, update: update, reorder: reorder, delete: delete}
}

type languageRequest struct {
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
}

type reorderLanguagesRequest struct {
	IDs []string `json:"ids"`
}

func (h *Handler) Create(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body languageRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	id, err := h.create.Handle(c.Request.Context(), commands.CreateLanguageCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Input: langInput(body)})
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
	xs, err := h.list.Handle(c.Request.Context(), queries.ListLanguagesQuery{ResumeID: rid, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, xs)
}

func (h *Handler) Get(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "languageId")
	if !ok {
		return
	}
	x, err := h.get.Handle(c.Request.Context(), queries.GetLanguageQuery{ResumeID: rid, LanguageID: id, UserID: uid})
	if err != nil {
		common.WriteQueryError(c, err)
		return
	}
	c.JSON(http.StatusOK, x)
}

func (h *Handler) Update(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "languageId")
	if !ok {
		return
	}
	var body languageRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	if err := h.update.Handle(c.Request.Context(), commands.UpdateLanguageCommand{ResumeID: rid, LanguageID: id, UserID: uid, Input: langInput(body)}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Reorder(c *gin.Context) {
	uid, rid, ok := ids(c)
	if !ok {
		return
	}
	var body reorderLanguagesRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	ids := make([]valueobjects.ID, len(body.IDs))
	for i, raw := range body.IDs {
		id, err := valueobjects.ParseID(raw)
		if err != nil {
			common.WriteInvalidBody(c)
			return
		}
		ids[i] = id
	}
	if err := h.reorder.Handle(c.Request.Context(), commands.ReorderLanguagesCommand{ResumeID: rid, UserID: uid, LanguageIDs: ids}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	uid, rid, id, ok := itemIDs(c, "languageId")
	if !ok {
		return
	}
	if err := h.delete.Handle(c.Request.Context(), commands.DeleteLanguageCommand{ResumeID: rid, LanguageID: id, UserID: uid}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func langInput(body languageRequest) commands.LanguageInput {
	input := commands.LanguageInput{ID: valueobjects.NewID(), Language: enums.SpokenLanguage(body.Language)}
	if body.Grade != nil {
		grade := enums.LanguageGrade(*body.Grade)
		input.Grade = &grade
	}
	return input
}

func ids(c *gin.Context) (useridentity.ID, valueobjects.ID, bool) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return useridentity.ID{}, valueobjects.ID{}, false
	}
	rid, ok := common.ParseResumeID(c)
	return uid, rid, ok
}

func itemIDs(c *gin.Context, name string) (useridentity.ID, valueobjects.ID, valueobjects.ID, bool) {
	uid, rid, ok := ids(c)
	if !ok {
		return uid, rid, valueobjects.ID{}, false
	}
	id, err := valueobjects.ParseID(c.Param(name))
	if err != nil {
		common.WriteInvalidBody(c)
		return uid, rid, valueobjects.ID{}, false
	}
	return uid, rid, id, true
}
