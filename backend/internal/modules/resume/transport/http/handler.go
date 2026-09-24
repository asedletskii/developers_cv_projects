package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
)

type Handler struct {
	getResume    *queries.GetResumeHandler
	listResumes  *queries.ListResumesHandler
	createResume *commands.CreateResumeHandler
	updateResume *commands.UpdateResumeHandler
	patchResume  *commands.PatchResumeHandler
}

func NewHandler(
	getResume *queries.GetResumeHandler,
	listResumes *queries.ListResumesHandler,
	createResume *commands.CreateResumeHandler,
	updateResume *commands.UpdateResumeHandler,
	patchResume *commands.PatchResumeHandler,
) *Handler {
	return &Handler{
		getResume:    getResume,
		listResumes:  listResumes,
		createResume: createResume,
		updateResume: updateResume,
		patchResume:  patchResume,
	}
}

func (h *Handler) ListResumes(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		writeError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return
	}
	id, err := useridentity.IDFromUUID(userID)
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}
	items, err := h.listResumes.Handle(c.Request.Context(), queries.ListResumesQuery{UserID: id})
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetResume(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		writeError(c, apperrors.New(
			apperrors.CodeUnauthorized,
			http.StatusUnauthorized,
			"Authentication required",
		))
		return
	}

	query, ok := buildGetResumeQuery(c, userIDFromUUID(userID))
	if !ok {
		return
	}
	h.getResumeResponse(c, query)
}

func (h *Handler) GetPublicResume(c *gin.Context) {
	query, ok := buildGetResumeQuery(c, nil)
	if !ok {
		return
	}
	h.getResumeResponse(c, query)
}

func buildGetResumeQuery(c *gin.Context, userID *useridentity.ID) (queries.GetResumeQuery, bool) {
	id, err := uuid.Parse(c.Param("resumeId"))
	if err != nil {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return queries.GetResumeQuery{}, false
	}

	lang := queries.Language(c.DefaultQuery("lang", string(queries.LanguageRU)))
	if lang != queries.LanguageRU && lang != queries.LanguageEN {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid language"))
		return queries.GetResumeQuery{}, false
	}

	return queries.GetResumeQuery{
		UserID: userID,
		ID:     id,
		Lang:   lang,
	}, true
}

func (h *Handler) getResumeResponse(c *gin.Context, query queries.GetResumeQuery) {
	resume, err := h.getResume.Handle(c.Request.Context(), query)
	if err != nil {
		if errors.Is(err, queries.ErrResumeNotFound) {
			writeError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
			return
		}
		writeError(c, apperrors.Internal(err))
		return
	}

	c.JSON(http.StatusOK, resume)
}

func (h *Handler) UpdateResume(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		writeError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return
	}
	resumeID, err := uuid.Parse(c.Param("resumeId"))
	if err != nil {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return
	}
	var request createResumeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}
	data, err := request.toCommand(userID)
	if err != nil {
		if isValidationError(err) {
			writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, err.Error()))
			return
		}
		writeError(c, apperrors.Internal(err))
		return
	}
	identity, err := useridentity.IDFromUUID(userID)
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}
	resumeIdentity, err := valueobjects.IDFromUUID(resumeID)
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}
	_, err = h.updateResume.Handle(c.Request.Context(), commands.UpdateResumeCommand{
		ResumeID:       resumeIdentity,
		UserID:         identity,
		Name:           data.Name,
		Surname:        data.Surname,
		Patronymic:     data.Patronymic,
		Position:       data.Position,
		Availability:   data.Availability,
		ResumeLink:     data.ResumeLink,
		Theme:          data.Theme,
		Sections:       data.Sections,
		HiddenContacts: data.HiddenContacts,
		ContactsText:   data.ContactsText,
		Greeting:       data.Greeting,
		PhotoPath:      data.PhotoPath,
		DateOfBirth:    data.DateOfBirth,
		Location:       data.Location,
		Relocation:     data.Relocation,
		WorkFormats:    data.WorkFormats,
		Contacts:       data.Contacts,
		Languages:      data.Languages,
		Experience:     data.Experience,
		Education:      data.Education,
		Courses:        data.Courses,
		Activities:     data.Activities,
		HardSkills:     data.HardSkills,
		SoftSkills:     data.SoftSkills,
		Technologies:   data.Technologies,
		Projects:       data.Projects,
	})
	if err != nil {
		writeCommandError(c, err)
		return
	}

	resume, err := h.getResume.Handle(c.Request.Context(), queries.GetResumeQuery{
		UserID: &identity,
		ID:     resumeID,
		Lang:   queries.LanguageRU,
	})
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}
	c.JSON(http.StatusOK, resume)
}

type patchResumeRequest struct {
	Title *string `json:"title"`
}

func (h *Handler) PatchResume(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		writeError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return
	}
	resumeID, err := uuid.Parse(c.Param("resumeId"))
	if err != nil {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return
	}
	var request patchResumeRequest
	if err := c.ShouldBindJSON(&request); err != nil || request.Title == nil {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
		return
	}
	identity, err := useridentity.IDFromUUID(userID)
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}
	resumeIdentity, err := valueobjects.IDFromUUID(resumeID)
	if err != nil {
		writeError(c, apperrors.Internal(err))
		return
	}
	result, err := h.patchResume.Handle(c.Request.Context(), commands.PatchResumeCommand{
		ResumeID: resumeIdentity,
		UserID:   identity,
		Title:    *request.Title,
	})
	if err != nil {
		writeCommandError(c, err)
		return
	}
	c.JSON(http.StatusOK, queries.ResumeListItem{ID: result.ResumeID.String(), Title: result.Title})
}

func writeCommandError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrEntityNotFound) {
		writeError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
		return
	}
	if isValidationError(err) {
		writeError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, err.Error()))
		return
	}
	writeError(c, apperrors.Internal(err))
}

func userIDFromUUID(id uuid.UUID) *useridentity.ID {
	value, err := useridentity.IDFromUUID(id)
	if err != nil {
		return nil
	}
	return &value
}

func writeError(c *gin.Context, err error) {
	appErr := apperrors.From(err)
	requestID := middleware.GetRequestID(c)
	c.AbortWithStatusJSON(
		appErr.Status,
		gin.H{
			"error": gin.H{
				"code":      appErr.Code,
				"message":   appErr.Message,
				"requestId": requestID,
			},
		},
	)
}
