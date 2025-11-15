package errors

import (
	"context"
	"encoding/json"

	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type ApplicationError interface {
	InternalCode() string
	StatusCode() int
	Message() string
	OriginalMessage() string
	AddExtraFields(fields map[string]interface{}) ApplicationError
	IsBasedOn(mappings.ErrorDetails) bool
	Log(context.Context)
	json.Marshaler
	error
}

type applicationError struct {
	errorDetails    mappings.ErrorDetails
	originalMessage string
	extraFields     map[string]interface{}
}

func (r *applicationError) InternalCode() string {
	return r.errorDetails.InternalCode
}

func (r *applicationError) StatusCode() int {
	return r.errorDetails.StatusCode
}

func (r *applicationError) Message() string {
	return r.errorDetails.Message
}

func (r *applicationError) OriginalMessage() string {
	return r.originalMessage
}

func (r *applicationError) Error() string {
	return r.errorDetails.Message
}

func (r *applicationError) AddExtraFields(fields map[string]interface{}) ApplicationError {
	if r.extraFields == nil {
		r.extraFields = make(map[string]interface{})
	}
	for k, v := range fields {
		r.extraFields[k] = v
	}
	return r
}

func (r *applicationError) IsBasedOn(errorDetails mappings.ErrorDetails) bool {
	return r.errorDetails == errorDetails
}

func (r *applicationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    r.InternalCode(),
		Message: r.Message(),
	})
}

func NewApplicationError(
	code mappings.ErrorDetails,
	originalError error,
) ApplicationError {
	var originalMessage string

	if originalError != nil {
		originalMessage = originalError.Error()
	}

	return &applicationError{code, originalMessage, make(map[string]interface{})}
}
