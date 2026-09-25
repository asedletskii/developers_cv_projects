package profile

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
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

type ProfileRequest struct {
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

func (r ProfileRequest) input() commands.ProfileInput {
	return commands.ProfileInput{Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic, Position: r.Position, PhotoPath: r.PhotoPath, DateOfBirth: r.DateOfBirth, Location: r.Location, Relocation: r.Relocation, WorkFormats: r.WorkFormat}
}

func (h *Handler) Create(c *gin.Context) {
	if h == nil || h.create == nil {
		common.WriteError(c, errors.New("create profile handler dependencies are nil"))
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
	var body ProfileRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	params := body.input()
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
	if h == nil || h.get == nil {
		common.WriteError(c, errors.New("get profile handler dependencies are nil"))
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
	if h == nil || h.update == nil {
		common.WriteError(c, errors.New("update profile handler dependencies are nil"))
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
	var body ProfileRequest
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	params := body.input()
	if e := h.update.Handle(c.Request.Context(), commands.UpdateResumeProfileCommand{ResumeID: rid, UserID: uid, Profile: params}); e != nil {
		writeProfileError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func writeProfileError(c *gin.Context, e error) {
	if errors.Is(e, domain.ErrDuplicateEntity) {
		common.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume profile already exists"))
		return
	}
	common.WriteCommandError(c, e)
}
