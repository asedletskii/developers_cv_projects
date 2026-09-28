package languages

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
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
		var id valueobjects.ID
		if item.ID != "" {
			parsed, err := valueobjects.ParseID(item.ID)
			if err != nil {
				return commands.LanguagesInput{}, err
			}
			id = parsed
		}
		var grade *enums.LanguageGrade
		if item.Grade != nil {
			value := enums.LanguageGrade(*item.Grade)
			grade = &value
		}
		items = append(items, commands.LanguageInput{
			ID: id, Language: enums.SpokenLanguage(item.Language), Grade: grade, Position: item.Position,
		})
	}
	return commands.LanguagesInput{Entries: items}, nil
}

func (h *Handler) Replace(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var body languagesRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	if body.Version < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version must be positive"})
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
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
