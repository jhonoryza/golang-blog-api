package controller

import (
	request "api_blog/delivery/http/request"
	"api_blog/delivery/http/response"
	"api_blog/infrastructure/exception"
	"api_blog/usecase/upload"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type UploadController struct {
	Usecase *upload.UploadUsecase
}

func NewUploadController(uc *upload.UploadUsecase) *UploadController {
	return &UploadController{Usecase: uc}
}

func (c *UploadController) Presign(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req request.PresignUploadRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	result, err := c.Usecase.PresignUpload(r.Context(), upload.PresignInput{
		Filename: req.Filename,
		Size:     req.Size,
	})
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    result,
	}
	resp.ToJson(w)
}
