package activities

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/transport/http/common"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type Handler struct {
	replace *commands.ReplaceActivitiesHandler
}

func NewHandler(replace *commands.ReplaceActivitiesHandler) *Handler {
	return &Handler{replace: replace}
}

type activityRequest struct {
	ID          string  `json:"id,omitempty"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type activitiesRequest struct {
	Version int64             `json:"version"`
	Items   []activityRequest `json:"items"`
}

func (r activitiesRequest) input() (commands.ActivitiesInput, error) {
	out := make([]commands.ActivityInput, 0, len(r.Items))
	for _, x := range r.Items {
		var id sharedvalueobjects.ID
		if x.ID != "" {
			v, err := sharedvalueobjects.ParseID(x.ID)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			id = v
		}
		var date *domain.YearMonth
		if x.Date != nil {
			v, err := domain.ParseYearMonth(*x.Date)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			date = &v
		}
		var link *domain.ExternalURL
		if x.Link != nil {
			v, err := domain.NewExternalURL(*x.Link)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			link = &v
		}
		out = append(out, commands.ActivityInput{ID: id, Type: domain.ActivityType(x.Type), Title: x.Title, Description: x.Description, Date: date, Link: link})
	}
	return commands.ActivitiesInput{Items: out}, nil
}

func (h *Handler) Replace(c *gin.Context) {
	uid, rid, ok := common.RequireResumeContext(c)
	if !ok {
		return
	}
	var body activitiesRequest
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
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceActivitiesCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input}); err != nil {
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
