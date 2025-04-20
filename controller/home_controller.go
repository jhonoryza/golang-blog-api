package controller

import (
	"api_blog/response"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HomeIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    nil,
	}
	resp.ToJson(w)
}
