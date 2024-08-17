package exception

import (
	"api_blog/response"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if notFoundError(w, r, err) {
		return
	}
	if validationError(w, r, err) {
		return
	}
	internalServerError(w, r, err)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	errString := fmt.Sprintf("%v", err)
	resp := response.ApiResponse{
		Code:    http.StatusInternalServerError,
		Message: "INTERNAL_SERVER_ERROR",
		Data:    errString,
	}
	log.Printf("internal server exception: %v\n", err)
	resp.ToJson(w)
}

func notFoundError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
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

func validationError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
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
		resp.ToJson(w)
		return true
	}
}
