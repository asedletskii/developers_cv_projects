package contacts

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/entities"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"
)

type Handler struct {
	create *commands.CreateResumeContactsHandler
	get    *queries.GetResumeContactsQueryHandler
	update *commands.UpdateResumeContactsHandler
}

func NewHandler(c *commands.CreateResumeContactsHandler, g *queries.GetResumeContactsQueryHandler, u *commands.UpdateResumeContactsHandler) *Handler {
	return &Handler{create: c, get: g, update: u}
}

type Request struct {
	Phone          *string              `json:"phone,omitempty"`
	Telegram       *string              `json:"telegram,omitempty"`
	Email          *string              `json:"email,omitempty"`
	HiddenContacts []string             `json:"hiddenContacts"`
	ContactsText   *ContactsTextRequest `json:"contactsText,omitempty"`
}
type ContactsTextRequest struct {
	Title    *LocalizedTextRequest `json:"title,omitempty"`
	Subtitle *LocalizedTextRequest `json:"subtitle,omitempty"`
}
type LocalizedTextRequest struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}
type Response struct {
	ID             string                `json:"id"`
	ResumeID       string                `json:"resumeId"`
	Phone          *string               `json:"phone,omitempty"`
	Telegram       *string               `json:"telegram,omitempty"`
	Email          *string               `json:"email,omitempty"`
	HiddenContacts []string              `json:"hiddenContacts"`
	ContactsText   *ContactsTextResponse `json:"contactsText,omitempty"`
}
type ContactsTextResponse struct {
	Title    *LocalizedTextResponse `json:"title,omitempty"`
	Subtitle *LocalizedTextResponse `json:"subtitle,omitempty"`
}
type LocalizedTextResponse struct {
	RU string `json:"ru"`
	EN string `json:"en"`
}

func (r Request) params(rid valueobjects.ID) (domain.ContactsParams, error) {
	var tg *valueobjects.Telegram
	if r.Telegram != nil {
		x, e := valueobjects.NewTelegram(*r.Telegram)
		if e != nil {
			return domain.ContactsParams{}, e
		}
		tg = &x
	}
	var em *valueobjects.Email
	if r.Email != nil {
		x, e := valueobjects.NewEmail(*r.Email)
		if e != nil {
			return domain.ContactsParams{}, e
		}
		em = &x
	}
	cs, e := entities.NewContacts(entities.ContactsParams{Phone: stringValue(r.Phone), Telegram: tg, Email: em})
	if e != nil {
		return domain.ContactsParams{}, e
	}
	hidden := make([]enums.ContactKey, len(r.HiddenContacts))
	for i, v := range r.HiddenContacts {
		hidden[i] = enums.ContactKey(v)
	}
	var text *valueobjects.ContactsText
	if r.ContactsText != nil {
		var title, subtitle *valueobjects.LocalizedText
		if r.ContactsText.Title != nil {
			x, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: r.ContactsText.Title.RU, EN: r.ContactsText.Title.EN})
			if e != nil {
				return domain.ContactsParams{}, e
			}
			title = &x
		}
		if r.ContactsText.Subtitle != nil {
			x, e := valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: r.ContactsText.Subtitle.RU, EN: r.ContactsText.Subtitle.EN})
			if e != nil {
				return domain.ContactsParams{}, e
			}
			subtitle = &x
		}
		x, e := valueobjects.NewContactsText(valueobjects.ContactsTextParams{Title: title, Subtitle: subtitle})
		if e != nil {
			return domain.ContactsParams{}, e
		}
		text = &x
	}
	return domain.ContactsParams{ResumeID: rid, Contacts: cs, HiddenContacts: hidden, Text: text}, nil
}
func (h *Handler) Create(c *gin.Context) { h.handleCreate(c) }
func (h *Handler) handleCreate(c *gin.Context) {
	if h == nil || h.create == nil || h.get == nil {
		common.WriteError(c, errors.New("create contacts handler dependencies are nil"))
		return
	}
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var b Request
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	p, e := b.params(rid)
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	_, e = h.create.Handle(c.Request.Context(), commands.CreateResumeContactsCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Contacts: p})
	if e != nil {
		writeError(c, e)
		return
	}
	a, e := h.get.Handle(c.Request.Context(), queries.GetResumeContactsQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusCreated, a)
}
func (h *Handler) Get(c *gin.Context) {
	if h == nil || h.get == nil {
		common.WriteError(c, errors.New("get contacts handler dependencies are nil"))
		return
	}
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
	if h == nil || h.update == nil || h.get == nil {
		common.WriteError(c, errors.New("update contacts handler dependencies are nil"))
		return
	}
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var b Request
	if common.BindJSON(c, &b) != nil {
		common.WriteInvalidBody(c)
		return
	}
	p, e := b.params(rid)
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	if e = h.update.Handle(c.Request.Context(), commands.UpdateResumeContactsCommand{ResumeID: rid, UserID: uid, Contacts: p}); e != nil {
		writeError(c, e)
		return
	}
	a, e := h.get.Handle(c.Request.Context(), queries.GetResumeContactsQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, a)
}
func writeError(c *gin.Context, e error) {
	if errors.Is(e, domain.ErrDuplicateEntity) {
		common.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume contacts already exists"))
		return
	}
	common.WriteCommandError(c, e)
}
func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
