package exception

import (
	"api_blog/response"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if NotFoundServerError(w, r, err) {
		return
	}
	if ValidationError(w, r, err) {
		return
	}
	InternalServerError(w, r, err)
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	errString := fmt.Sprintf("%v", err)
	resp := response.ApiResponse{
		Code:    http.StatusBadRequest,
		Message: "BAD REQUEST",
		Data:    errString,
	}
	sentry.CaptureException(errors.New(errString))
	sentry.Flush(2 * time.Second)
	log.Printf("bad request exception: %v\n", err)
	resp.ToJson(w)
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	errString := fmt.Sprintf("%v", err)
	resp := response.ApiResponse{
		Code:    http.StatusUnauthorized,
		Message: "UNAUTHORIZED",
		Data:    errString,
	}
	sentry.CaptureException(errors.New(errString))
	sentry.Flush(2 * time.Second)
	log.Printf("unauthorized exception: %v\n", err)
	resp.ToJson(w)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	errString := fmt.Sprintf("%v", err)
	resp := response.ApiResponse{
		Code:    http.StatusInternalServerError,
		Message: "INTERNAL_SERVER_ERROR",
		Data:    errString,
	}
	sentry.CaptureException(errors.New(errString))
	sentry.Flush(2 * time.Second)
	log.Printf("internal server exception: %v\n", err)
	resp.ToJson(w)
}

func NotFoundServerError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if !ok {
		return false
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := response.ApiResponse{
			Code:    http.StatusNotFound,
			Message: "NOT FOUND",
			Data:    nil,
		}
		log.Printf("not found exception: %v\n", exception.Error)
		resp.ToJson(w)
		return true
	}
}

func ValidationError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := response.ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "BAD REQUEST",
			Data:    exception.Error(),
		}
		sentry.CaptureException(errors.New(exception.Error()))
		sentry.Flush(2 * time.Second)
		resp.ToJson(w)
		return true
	}
}
