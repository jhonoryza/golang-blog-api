package response

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (resp *ApiResponse) ToJson(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(resp)
	if err != nil {
		panic(err)
	}
}
