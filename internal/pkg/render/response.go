package render

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/buni/tx-parser/internal/pkg/sloglog"
	"github.com/go-chi/chi/v5/middleware"
)

// ErrorResponse is a generic error response.
type ErrorResponse struct {
	Error *Error `json:"error"`
}

// NewSuccessHandlerResponse is a generic success response writer.
func NewSuccessResponse[Response any](ctx context.Context, w http.ResponseWriter, status int, response Response) {
	write(ctx, w, status, response)
}

// NewErrorResponse is a generic error response writer.
func NewErrorResponse(ctx context.Context, w http.ResponseWriter, code int, status string, err error) {
	responseError := NewError(status, err.Error())

	errors.As(err, &responseError)

	write(
		ctx,
		w,
		code,
		ErrorResponse{Error: responseError},
	)
}

// NewSuccessOKResponse is a success response writer with status code 200.
func NewSuccessOKResponse[Response any](ctx context.Context, w http.ResponseWriter, response Response) {
	NewSuccessResponse(ctx, w, http.StatusOK, response)
}

// NewSuccessCreatedResponse is a success response writer with status code 201.
func NewSuccessCreatedResponse[Response any](ctx context.Context, w http.ResponseWriter, response Response) {
	NewSuccessResponse(ctx, w, http.StatusCreated, response)
}

// NewInternalServerErrorResponse is an internal server error response writer.
func NewInternalServerErrorResponse(ctx context.Context, w http.ResponseWriter, _ error) {
	NewErrorResponse(ctx, w, http.StatusInternalServerError, InternalServerError, errors.New("internal server error")) //nolint:goerr113
}

// NewNotFoundErrorResponse is a not found error response writer.
func NewNotFoundErrorResponse(ctx context.Context, w http.ResponseWriter, _ error) {
	NewErrorResponse(ctx, w, http.StatusNotFound, NotFoundError, errors.New("not found")) //nolint:goerr113
}

// NewNotFoundErrorResponse is a bad request error response writer.
func NewBadRequestErrorResponse(ctx context.Context, w http.ResponseWriter, _ error) {
	NewErrorResponse(ctx, w, http.StatusBadRequest, BadRequestError, errors.New("bad request")) //nolint:goerr113
}

// NewConflictErrorResponse is a conflict error response writer.
func NewConflictErrorResponse(ctx context.Context, w http.ResponseWriter, _ error) {
	NewErrorResponse(ctx, w, http.StatusConflict, ConflictError, errors.New("conflict with existing resource")) //nolint:goerr113
}

// NewValidationErrorResponse is a validation error response writer.
func NewValidationErrorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	NewErrorResponse(ctx, w, http.StatusBadRequest, RequestValidationError, err) //nolint:goerr113
}

func write[Response any](ctx context.Context, w http.ResponseWriter, status int, v Response) {
	log := sloglog.FromContext(ctx)

	ww, ok := w.(middleware.WrapResponseWriter)
	if !ok {
		ww = middleware.NewWrapResponseWriter(w, 1)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if ww.BytesWritten() != 0 {
		log.DebugContext(ctx, "response already written, in renderer")
		return
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	err := enc.Encode(v)
	if err != nil {
		log.ErrorContext(ctx, "error encoding response, in renderer", sloglog.Error(err))
		http.Error(ww, InternalServerError, http.StatusInternalServerError)
		return
	}

	if ww.Status() == 0 {
		ww.WriteHeader(status)
	} else {
		log.DebugContext(ctx, "status resp already written, in renderer")
	}

	_, err = ww.Write(buf.Bytes())
	if err != nil {
		log.ErrorContext(ctx, "error writing response, in renderer", sloglog.Error(err))
		return
	}
}
