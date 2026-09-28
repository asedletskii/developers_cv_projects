package activities

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
		var id valueobjects.ID
		if x.ID != "" {
			v, err := valueobjects.ParseID(x.ID)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			id = v
		}
		var date *valueobjects.YearMonth
		if x.Date != nil {
			v, err := valueobjects.ParseYearMonth(*x.Date)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			date = &v
		}
		var link *valueobjects.ExternalURL
		if x.Link != nil {
			v, err := valueobjects.NewExternalURL(*x.Link)
			if err != nil {
				return commands.ActivitiesInput{}, err
			}
			link = &v
		}
		out = append(out, commands.ActivityInput{ID: id, Type: enums.ActivityType(x.Type), Title: x.Title, Description: x.Description, Date: date, Link: link})
	}
	return commands.ActivitiesInput{Items: out}, nil
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
	var body activitiesRequest
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
	if err := h.replace.Handle(c.Request.Context(), commands.ReplaceActivitiesCommand{ResumeID: rid, UserID: uid, Version: body.Version, Input: input}); err != nil {
		if errors.Is(err, domain.ErrConcurrencyConflict) {
			c.Status(http.StatusPreconditionFailed)
			return
		}
		common.WriteCommandError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
