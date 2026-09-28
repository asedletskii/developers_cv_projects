package contacts

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
	replace *commands.ReplaceContactsHandler
}

func NewHandler(replace *commands.ReplaceContactsHandler) *Handler {
	return &Handler{replace: replace}
}

type contactsRequest struct {
	Version        int64                `json:"version"`
	ID             string               `json:"id,omitempty"`
	Phone          *string              `json:"phone,omitempty"`
	Telegram       *string              `json:"telegram,omitempty"`
	Email          *string              `json:"email,omitempty"`
	HiddenContacts []string             `json:"hiddenContacts"`
	ContactsText   *contactsTextRequest `json:"contactsText,omitempty"`
}

type contactsTextRequest struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}

func (r contactsRequest) input() (commands.ContactsInput, error) {
	var telegram *valueobjects.Telegram
	if r.Telegram != nil {
		value, err := valueobjects.NewTelegram(*r.Telegram)
		if err != nil {
			return commands.ContactsInput{}, err
		}
		telegram = &value
	}

	var email *valueobjects.Email
	if r.Email != nil {
		value, err := valueobjects.NewEmail(*r.Email)
		if err != nil {
			return commands.ContactsInput{}, err
		}
		email = &value
	}

	hidden := make([]enums.ContactKey, len(r.HiddenContacts))
	for i, value := range r.HiddenContacts {
		hidden[i] = enums.ContactKey(value)
	}

	var text *commands.ContactsTextInput
	if r.ContactsText != nil {
		text = &commands.ContactsTextInput{Title: r.ContactsText.Title, Subtitle: r.ContactsText.Subtitle}
	}

	return commands.ContactsInput{
		Phone:          r.Phone,
		Telegram:       telegram,
		Email:          email,
		HiddenContacts: hidden,
		ContactsText:   text,
	}, nil
}

func (h *Handler) Replace(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	resumeID, ok := common.ParseResumeID(c)
	if !ok {
		return
	}

	var body contactsRequest
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

	if body.ID != "" {
		id, err := valueobjects.ParseID(body.ID)
		if err != nil {
			common.WriteCommandError(c, err)
			return
		}
		input.ID = id
	}

	err = h.replace.Handle(c.Request.Context(), commands.ReplaceContactsCommand{
		ResumeID: resumeID,
		UserID:   uid,
		Version:  body.Version,
		Input:    input,
	})
	if err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
