package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/queries"
	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	apperrors "developers-cv/backend/internal/platform/http/errors"
	"developers-cv/backend/internal/platform/http/middleware"
	useridentity "developers-cv/backend/internal/shared/domain/valueobjects"
)

var ErrRequestBodyTooLarge = errors.New("request body too large")

func RequireUserID(c *gin.Context) (useridentity.ID, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeUnauthorized, http.StatusUnauthorized, "Authentication required"))
		return useridentity.ID{}, false
	}
	identity, err := useridentity.IDFromUUID(userID)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return useridentity.ID{}, false
	}
	return identity, true
}

func ParseResumeID(c *gin.Context) (valueobjects.ID, bool) {
	parsed, err := uuid.Parse(c.Param("resumeId"))
	if err != nil {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid resume id"))
		return valueobjects.ID{}, false
	}
	id, err := valueobjects.IDFromUUID(parsed)
	if err != nil {
		apperrors.WriteError(c, apperrors.Internal(err))
		return valueobjects.ID{}, false
	}
	return id, true
}

func WriteCommandError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrEntityNotFound) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume resource not found"))
		return
	}
	if errors.Is(err, domain.ErrDuplicateEntity) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeConflict, http.StatusConflict, "Resume resource already exists"))
		return
	}
	if IsValidationError(err) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, err.Error()))
		return
	}
	apperrors.WriteError(c, apperrors.Internal(err))
}

func WriteInvalidBody(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeBadRequest, http.StatusBadRequest, "Invalid request body"))
}

func WriteBindError(c *gin.Context, err error) {
	if errors.Is(err, ErrRequestBodyTooLarge) {
		apperrors.WriteError(c, apperrors.New(apperrors.CodeRequestTooLarge, http.StatusRequestEntityTooLarge, "Request body too large"))
		return
	}
	WriteInvalidBody(c)
}

func WriteQueryError(c *gin.Context, err error) {
	if errors.Is(err, queries.ErrResumeNotFound) {
		WriteResumeNotFound(c)
		return
	}
	apperrors.WriteError(c, err)
}

func WriteResumeNotFound(c *gin.Context) {
	apperrors.WriteError(c, apperrors.New(apperrors.CodeNotFound, http.StatusNotFound, "Resume not found"))
}

func BindJSON(c *gin.Context, target any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			return ErrRequestBodyTooLarge
		}
		return err
	}
	if err := validateJSONStructure(body); err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("request body must contain a single JSON value")
		}
		return err
	}
	return nil
}

func validateJSONStructure(body []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok || delimiter != '{' {
		return errors.New("request body must be a JSON object")
	}
	if err := validateJSONObject(decoder); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("request body must contain a single JSON value")
		}
		return err
	}
	return nil
}

func validateJSONObject(decoder *json.Decoder) error {
	seen := map[string]struct{}{}
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return err
		}
		key, ok := keyToken.(string)
		if !ok {
			return fmt.Errorf("invalid JSON object key")
		}
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate JSON field %q", key)
		}
		seen[key] = struct{}{}
		if err := validateJSONValue(decoder); err != nil {
			return err
		}
	}
	_, err := decoder.Token()
	return err
}

func validateJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		return validateJSONObject(decoder)
	case '[':
		for decoder.More() {
			if err := validateJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	default:
		return nil
	}
}
