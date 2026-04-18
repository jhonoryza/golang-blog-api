package exception

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
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
	resp := map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "BAD REQUEST",
		"data":    errString,
	}
	sentry.CaptureException(errors.New(errString))
	sentry.Flush(2 * time.Second)
	log.Printf("bad request exception: %v\n", err)
	json.NewEncoder(w).Encode(resp)
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	errString := fmt.Sprintf("%v", err)
	resp := map[string]interface{}{
		"code":    http.StatusUnauthorized,
		"message": "UNAUTHORIZED",
		"data":    errString,
	}
	sentry.CaptureException(errors.New(errString))
	sentry.Flush(2 * time.Second)
	log.Printf("unauthorized exception: %v\n", err)
	json.NewEncoder(w).Encode(resp)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	errString := fmt.Sprintf("%v", err)
	resp := map[string]interface{}{
		"code":    http.StatusInternalServerError,
		"message": "INTERNAL_SERVER_ERROR",
		"data":    errString,
	}
	sentry.CaptureException(errors.New(errString))
	sentry.Flush(2 * time.Second)
	log.Printf("internal server exception: %v\n", err)
	json.NewEncoder(w).Encode(resp)
}

func NotFoundServerError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if !ok {
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	resp := map[string]interface{}{
		"code":    http.StatusNotFound,
		"message": "NOT FOUND",
		"data":    nil,
	}
	log.Printf("not found exception: %v\n", exception.Error)
	json.NewEncoder(w).Encode(resp)
	return true
}

func ValidationError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp := map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "BAD REQUEST",
		"data":    exception.Error(),
	}
	sentry.CaptureException(errors.New(exception.Error()))
	sentry.Flush(2 * time.Second)
	json.NewEncoder(w).Encode(resp)
	return true
}

type PanicHandler func(http.ResponseWriter, *http.Request, interface{})

func NewPanicHandler() PanicHandler {
	return ErrorHandler
}

func NewNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":404,"message":"Not Found","data":null}`))
	})
}

func NotFoundHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":404,"message":"Not Found","data":null}`))
	}
}