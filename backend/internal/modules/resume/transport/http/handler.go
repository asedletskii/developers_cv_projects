package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"

	"developers-cv/backend/internal/modules/resume/application/commands"
	"developers-cv/backend/internal/modules/resume/application/queries"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
)

type Handler struct {
	getResume    *queries.GetResumeHandler
	listResumes  *queries.ListResumesHandler
	createResume *commands.CreateResumeHandler
}

func NewHandler(getResume *queries.GetResumeHandler, deps ...any) *Handler {
	var listResumes *queries.ListResumesHandler
	var createResume *commands.CreateResumeHandler
	for _, dep := range deps {
		switch value := dep.(type) {
		case *queries.ListResumesHandler:
			listResumes = value
		case *commands.CreateResumeHandler:
			createResume = value
		}
	}
	return &Handler{getResume: getResume, listResumes: listResumes, createResume: createResume}
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
	items, err := h.listResumes.Handle(c.Request.Context(), queries.GetListResumesQuery{UserID: id})
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
