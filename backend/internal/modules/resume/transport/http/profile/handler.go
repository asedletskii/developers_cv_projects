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

type Request struct {
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
type Response struct {
	ID          string   `json:"id"`
	ResumeID    string   `json:"resumeId"`
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

func (r Request) params(resumeID valueobjects.ID) (domain.ProfileParams, error) {
	dob, e := valueobjects.ParseDate(r.DateOfBirth)
	if e != nil {
		return domain.ProfileParams{}, e
	}
	photo, e := valueobjects.NewAssetPath(r.PhotoPath)
	if e != nil {
		return domain.ProfileParams{}, e
	}
	fs := make([]enums.WorkFormat, len(r.WorkFormat))
	for i, v := range r.WorkFormat {
		fs[i] = enums.WorkFormat(v)
	}
	return domain.ProfileParams{ResumeID: resumeID, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic, Position: r.Position, PhotoPath: photo, DateOfBirth: dob, Location: r.Location, Relocation: r.Relocation, WorkFormats: fs}, nil
}
func (h *Handler) Create(c *gin.Context) {
	if h == nil || h.create == nil || h.get == nil {
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
	var body Request
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	params, e := body.params(rid)
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	id, e := h.create.Handle(c.Request.Context(), commands.CreateResumeProfileCommand{ID: valueobjects.NewID(), ResumeID: rid, UserID: uid, Profile: params})
	if e != nil {
		writeProfileError(c, e)
		return
	}
	p, e := h.get.Handle(c.Request.Context(), queries.GetResumeProfileQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	_ = id
	c.JSON(http.StatusCreated, p)
}
func (h *Handler) Get(c *gin.Context) {
	if h == nil || h.get == nil {
		common.WriteError(c, errors.New("get profile handler dependencies are nil"))
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
	p, e := h.get.Handle(c.Request.Context(), queries.GetResumeProfileQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, p)
}
func (h *Handler) Update(c *gin.Context) {
	if h == nil || h.update == nil || h.get == nil {
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
	var body Request
	if common.BindJSON(c, &body) != nil {
		common.WriteInvalidBody(c)
		return
	}
	params, e := body.params(rid)
	if e != nil {
		common.WriteCommandError(c, e)
		return
	}
	if e = h.update.Handle(c.Request.Context(), commands.UpdateResumeProfileCommand{ResumeID: rid, UserID: uid, Profile: params}); e != nil {
		writeProfileError(c, e)
		return
	}
	p, e := h.get.Handle(c.Request.Context(), queries.GetResumeProfileQuery{ResumeID: rid, UserID: uid})
	if e != nil {
		common.WriteQueryError(c, e)
		return
	}
	c.JSON(http.StatusOK, p)
}
func writeProfileError(c *gin.Context, e error) {
	if errors.Is(e, domain.ErrDuplicateEntity) {
		common.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume profile already exists"))
		return
	}
	common.WriteCommandError(c, e)
}
