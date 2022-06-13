package helper

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/dto"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/internalerrors"
	"net/http"
	"strconv"
)

// ErrorResponse is used to static shape json return
type ErrorResponse struct {
	Errors []*dto.Error `json:"errors"`
} // @name Response.ErrorResponse

// ErrorResponseWithStatus is used to static shape json return
type ErrorResponseWithStatus struct {
	ErrorResponse ErrorResponse `json:"-" `
	Status        int           `json:"-"`
} // @name Response.ErrorResponseWithStatus

// BuildErrorResponse method is to inject data value to dynamic failed response
func BuildErrorResponse(ctx *gin.Context, message string, err error) ErrorResponseWithStatus {
	status := http.StatusInternalServerError
	title := "Internal Server Error"
	switch {
	case errors.Is(err, internalerrors.ErrExists) ||
		errors.Is(err, internalerrors.ErrConflict):
		status = http.StatusConflict
		title = "Conflict"
	case errors.Is(err, internalerrors.ErrNotFound) ||
		errors.Is(err, internalerrors.ErrStoreNotFound):
		status = http.StatusNotFound
		title = "Not Found"
	case errors.Is(err, internalerrors.ErrForbidden):
		status = http.StatusForbidden
		title = "Forbidden"
	case errors.Is(err, internalerrors.ErrBadRequest) ||
		errors.Is(err, internalerrors.ErrPageOffsetExceed) ||
		errors.Is(err, internalerrors.ErrPageLimitExceed):
		status = http.StatusBadRequest
		title = "Bad Request"
	case errors.Is(err, internalerrors.ErrUnprocessableEntity) ||
		errors.Is(err, internalerrors.ErrStoreIDMismatch):
		status = http.StatusUnprocessableEntity
		title = "Unprocessable Entity"
	case errors.Is(err, internalerrors.ErrMethodNotAllowed):
		status = http.StatusMethodNotAllowed
		title = "Method Not Allowed"
	}

	return buildErrorWithStatus(ctx, status, title, message, err)
}

// buildErrorWithStatus log full stack for 500 errors and build final structure for error
func buildErrorWithStatus(ctx *gin.Context, status int, title string, message string, err error) ErrorResponseWithStatus {
	var errorResponse ErrorResponse

	if status != http.StatusInternalServerError {
		errorResponse = ErrorResponse{
			Errors: []*dto.Error{
				{
					Status: strconv.Itoa(status),
					Title:  title,
					Detail: message,
				},
			},
		}
	} else {
		// don't include detail if something could be a server error.
		// It logs the whole stack in case it's a server error and it's the debug level
		log.Printf("Error: %+v", err)
		errorResponse = ErrorResponse{
			Errors: []*dto.Error{{
				Status: strconv.Itoa(status),
				Title:  title,
				Detail: "there was a problem processing your request",
			}},
		}
	}
	errorResponseWithStatus := ErrorResponseWithStatus{
		ErrorResponse: errorResponse,
		Status:        status,
	}

	// Add error to context which is best practice in gin
	_ = ctx.Error(fmt.Errorf("title=%s, message=%s", title, message))
	return errorResponseWithStatus
}
