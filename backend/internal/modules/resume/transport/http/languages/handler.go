package languages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	replace *commands.ReplaceLanguagesHandler
}

func NewHandler(replace *commands.ReplaceLanguagesHandler) *Handler {
	return &Handler{replace: replace}
}

type languagesRequest struct {
	Version int64             `json:"version"`
	Items   []languageRequest `json:"items"`
}

type languageRequest struct {
	ID       string  `json:"id,omitempty"`
	Language string  `json:"language"`
	Grade    *string `json:"grade,omitempty"`
	Position int     `json:"position"`
}

func (r languagesRequest) input() (commands.LanguagesInput, error) {
	items := make([]commands.LanguageInput, 0, len(r.Items))
	for _, item := range r.Items {
		var id sharedvalueobjects.ID
		if item.ID != "" {
			parsed, err := sharedvalueobjects.ParseID(item.ID)
			if err != nil {
				return commands.LanguagesInput{}, err
			}
			id = parsed
		}
		var grade *domain.LanguageGrade
		if item.Grade != nil {
			value := domain.LanguageGrade(*item.Grade)
			grade = &value
		}
		items = append(items, commands.LanguageInput{
			ID: id, Language: domain.SpokenLanguage(item.Language), Grade: grade, Position: item.Position,
		})
	}
	return commands.LanguagesInput{Entries: items}, nil
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body languagesRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if !common.RequirePositiveVersion(c, body.Version) {
		return
	}
	input, err := body.input()
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceLanguagesCommand{
		ResumeID: rid, UserID: uid, Version: body.Version, Input: input,
	}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
