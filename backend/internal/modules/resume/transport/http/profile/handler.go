package profile

import (
	"developers-cv/backend/internal/modules/resume/application/commands"
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
	replace *commands.ReplaceProfileHandler
}

func NewHandler(replace *commands.ReplaceProfileHandler) *Handler { return &Handler{replace: replace} }

type profileRequest struct {
	Version     int64    `json:"version"`
	ID          string   `json:"id,omitempty"`
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
	var id valueobjects.ID
	if r.ID != "" {
		parsed, err := valueobjects.ParseID(r.ID)
		if err != nil {
			return commands.ProfileInput{}, err
		}
		id = parsed
	}
	photo, err := valueobjects.NewExternalURL(r.PhotoPath)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	dob, err := valueobjects.ParseDate(r.DateOfBirth)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	formats := make([]enums.WorkFormat, len(r.WorkFormat))
	for i, v := range r.WorkFormat {
		formats[i] = enums.WorkFormat(v)
	}
	return commands.ProfileInput{ID: id, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic, Position: r.Position, PhotoPath: photo, DateOfBirth: dob, Location: r.Location, Relocation: r.Relocation, WorkFormats: formats}, nil
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
	var body profileRequest
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
	err = h.replace.Handle(c.Request.Context(), commands.ReplaceProfileCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input})
	if err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		writeProfileError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func writeProfileError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrDuplicateEntity) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume profile has duplicate values"))
		return
	}
	common.WriteCommandError(c, err)
}
