package profile

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	create *commands.CreateResumeProfileHandler
	get    *queries.GetResumeProfileQueryHandler
	update *commands.UpdateResumeProfileHandler
}

func NewHandler(c *commands.CreateResumeProfileHandler, g *queries.GetResumeProfileQueryHandler, u *commands.UpdateResumeProfileHandler) *Handler {
	return &Handler{create: c, get: g, update: u}
}

type profileRequest struct {
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Patronymic  *string  `json:"patronymic,omitempty"`
	Position    string   `json:"position"`
	PhotoPath   string   `json:"photoPath"`
	DateOfBirth string   `json:"dateOfBirth"`
	Location    string   `json:"location"`
	Relocation  bool     `json:"relocation"`
	WorkFormat  []string `json:"workFormat"`
}

func (r profileRequest) input() (commands.ProfileInput, error) {
	photoURL, err := valueobjects.NewExternalURL(r.PhotoPath)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	dateOfBirth, err := valueobjects.ParseDate(r.DateOfBirth)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	workFormats := make([]enums.WorkFormat, len(r.WorkFormat))
	for i, value := range r.WorkFormat {
		workFormats[i] = enums.WorkFormat(value)
	}
	return commands.ProfileInput{
		Name:        r.Name,
		Surname:     r.Surname,
		Patronymic:  r.Patronymic,
		Position:    r.Position,
		PhotoPath:   photoURL,
		DateOfBirth: dateOfBirth,
		Location:    r.Location,
		Relocation:  r.Relocation,
		WorkFormats: workFormats,
	}, nil
}

func (h *Handler) Create(c *gin.Context) {
	uid, ok := common.RequireUserID(c)
	if !ok {
		return
	}
	rid, ok := common.ParseResumeID(c)
	if !ok {
		return
	}
	var body profileRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	params, e := body.input()
	if e != nil {
		writeProfileError(c, e)
		return
	}
	id, e := h.create.Handle(c.Request.Context(), commands.CreateResumeProfileCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Profile: params})
	if e != nil {
		writeProfileError(c, e)
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
	p, e := h.get.Handle(c.Request.Context(), queries.GetResumeProfileQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, p)
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
	var body profileRequest
	if err := common.BindJSON(c, &body); err != nil {
		common.WriteBindError(c, err)
		return
	}
	params, e := body.input()
	if e != nil {
		writeProfileError(c, e)
		return
	}
	if e = h.update.Handle(c.Request.Context(), commands.UpdateResumeProfileCommand{ResumeID: rid, UserID: uid, Profile: params}); e != nil {
		writeProfileError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func writeProfileError(c *gin.Context, e error) {
	if errors.Is(e, domain.ErrDuplicateEntity) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume profile already exists"))
		return
	}
	common.WriteCommandError(c, e)
}
