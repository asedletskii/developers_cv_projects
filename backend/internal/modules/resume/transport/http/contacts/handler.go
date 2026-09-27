package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
)

type Handler struct {
	create *commands.CreateResumeContactsHandler
	get    *queries.GetResumeContactsQueryHandler
	update *commands.UpdateResumeContactsHandler
}

func NewHandler(c *commands.CreateResumeContactsHandler, g *queries.GetResumeContactsQueryHandler, u *commands.UpdateResumeContactsHandler) *Handler {
	return &Handler{create: c, get: g, update: u}
}

type contactsRequest struct {
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
		text = &commands.ContactsTextInput{}
		if r.ContactsText.Title != nil {
			text.Title = r.ContactsText.Title
		}
		if r.ContactsText.Subtitle != nil {
			text.Subtitle = r.ContactsText.Subtitle
		}
	}
	return commands.ContactsInput{Phone: r.Phone, Telegram: telegram, Email: email, HiddenContacts: hidden, ContactsText: text}, nil
}

func (h *Handler) Create(c *gin.Context) { h.handleCreate(c) }
func (h *Handler) handleCreate(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var b contactsRequest
	if err := common.BindJSON(c, &b); err != nil {
		common.WriteBindError(c, err)
		return
	}
	p, e := b.input()
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	id := valueobjects.NewID()
	_, e = h.create.Handle(c.Request.Context(), commands.CreateResumeContactsCommand{ID: id, ResumeID: rid, UserID: uid, Contacts: p})
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.JSON(http.StatusCreated, common.CreatedIDResponse{ID: id.String()})
}
func (h *Handler) Get(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	a, e := h.get.Handle(c.Request.Context(), queries.GetResumeContactsQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, a)
}
func (h *Handler) Update(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var b contactsRequest
	if err := common.BindJSON(c, &b); err != nil {
		common.WriteBindError(c, err)
		return
	}
	p, e := b.input()
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	if e = h.update.Handle(c.Request.Context(), commands.UpdateResumeContactsCommand{ResumeID: rid, UserID: uid, Contacts: p}); e != nil {
		common.WriteCommandError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
