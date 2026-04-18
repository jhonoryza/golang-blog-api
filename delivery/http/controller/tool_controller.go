package controller

import (
	"api_blog/infrastructure/exception"
	request "api_blog/delivery/http/request"
	"api_blog/delivery/http/response"
	"api_blog/usecase/tool"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type ToolController struct {
	Usecase *tool.ToolUsecase
	DB      *sql.DB
}

func NewToolController(uc *tool.ToolUsecase, db *sql.DB) *ToolController {
	return &ToolController{Usecase: uc, DB: db}
}

type IndexParam struct {
	Type string `validate:"required"`
}

func (c *ToolController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	typeBy := r.URL.Query().Get("type")

	validate := validator.New()
	indexParam := IndexParam{Type: typeBy}

	if err := validate.Struct(indexParam); err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	input := tool.FindAllInput{
		Type:    typeBy,
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
	}

	if input.SortDir == "" {
		input.SortDir = "desc"
	}
	if input.SortBy == "" {
		input.SortBy = "published_at"
	}

	tools := c.Usecase.FindAll(r.Context(), input)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    tools,
	}

	resp.ToJson(w)
}

func (c *ToolController) Store(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req request.CreateToolRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	result, err := c.Usecase.Create(req)
	if err != nil {
		exception.InternalServerError(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    201,
		Message: "OK",
		Data: map[string]any{
			"id":         result.ID,
			"name":       result.Name,
			"created_at": result.CreatedAt.Time.In(time.Local).Format(time.RFC822),
		},
	}

	resp.ToJson(w)
}

func (c *ToolController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req request.UpdateToolRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	toolId := p.ByName("toolId")
	existing := c.Usecase.FindOneById(toolId)
	if existing == nil {
		exception.BadRequestError(w, r, errors.New("record not found"))
		return
	}

	result, err := c.Usecase.Update(req, toolId)
	if err != nil {
		exception.InternalServerError(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    201,
		Message: "OK",
		Data: map[string]any{
			"id":         result.ID,
			"name":       result.Name,
			"updated_at": result.UpdatedAt.Time.In(time.Local).Format(time.RFC822),
		},
	}

	resp.ToJson(w)
}

func (c *ToolController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	toolId := p.ByName("toolId")
	if toolId == "" {
		exception.ErrorHandler(w, r, errors.New("toolId is required"))
		return
	}

	rowsAffected, err := c.Usecase.Delete(toolId)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    200,
		Message: "OK",
		Data: map[string]any{
			"rowsAffected": rowsAffected,
		},
	}

	resp.ToJson(w)
}