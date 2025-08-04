package errors

import (
	"github.com/Kofandr/Product_Accounting_Service/internal/repository"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strings"
)

var (
	ErrInvalidID        = NewErrorResponse(http.StatusBadRequest, "Invalid id")
	ErrInvalidJSON      = NewErrorResponse(http.StatusBadRequest, "Invalid JSON format")
	ErrNotFound         = NewErrorResponse(http.StatusNotFound, "Not found")
	ErrNotFoundCategory = NewErrorResponse(http.StatusNotFound, "Not found category")
	ErrServerError      = NewErrorResponse(http.StatusInternalServerError, "Server error")
	ErrDuplicateEntry   = NewErrorResponse(http.StatusConflict, "Duplicate entry")
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e ErrorResponse) ToJSON() map[string]string {
	return map[string]string{"err": e.Message}
}

func NewErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{Status: status, Message: message}
}

func MapErrorToResponse(err error) ErrorResponse {
	switch {
	case err == nil:
		return ErrorResponse{}
	case err == pgx.ErrNoRows:
		return ErrNotFound
	case err == repository.ErrDuplicate:
		return ErrDuplicateEntry
	case strings.Contains(err.Error(), "invalid id"):
		return ErrInvalidID
	case strings.Contains(err.Error(), "json"):
		return ErrInvalidJSON
	default:
		return ErrServerError
	}

}
