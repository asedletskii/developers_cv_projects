package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	identity "developers-cv/backend/internal/shared/domain/valueobjects"
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
	var id identity.ID
	if r.ID != "" {
		parsed, err := identity.ParseID(r.ID)
		if err != nil {
			return commands.ProfileInput{}, err
		}
		id = parsed
	}
	photo, err := domain.NewExternalURL(r.PhotoPath)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	dob, err := domain.ParseDate(r.DateOfBirth)
	if err != nil {
		return commands.ProfileInput{}, err
	}
	formats := make([]domain.WorkFormat, len(r.WorkFormat))
	for i, v := range r.WorkFormat {
		formats[i] = domain.WorkFormat(v)
	}
	return commands.ProfileInput{ID: id, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic, Position: r.Position, PhotoPath: photo, DateOfBirth: dob, Location: r.Location, Relocation: r.Relocation, WorkFormats: formats}, nil
}
func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body profileRequest
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
	err = h.replace.Handle(c.Request.Context(), commands.ReplaceProfileCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input})
	if err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
