package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type Code string

const (
	BadRequest      Code = "BAD_REQUEST"
	ValidationError Code = "VALIDATION_ERROR"
	InternalError   Code = "INTERNAL_ERROR"
	NotFound        Code = "NOT_FOUND"
	Conflict        Code = "CONFLICT"
)

type ErrorResponse struct {
	Code    Code   `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func Error(w http.ResponseWriter, r *http.Request, status int, code Code, err string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorResponse{
		Code:  code,
		Error: err,
	})
}

func OK(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}

func Created(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, data)
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}
