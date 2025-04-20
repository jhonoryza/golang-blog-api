package controller

import (
	"api_blog/exception"
	"api_blog/repository"
	"api_blog/requests"
	"api_blog/response"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type ToolController struct {
	DB *sql.DB
}

func NewToolController(db *sql.DB) *ToolController {
	return &ToolController{DB: db}
}

type IndexParam struct {
	Type string `validate:"required"`
}

func (c *ToolController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	typeBy := r.URL.Query().Get("type")

	validate := validator.New()
	indexParam := IndexParam{
		Type: typeBy,
	}

	if err := validate.Struct(indexParam); err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	if sortDir == "" {
		sortDir = "desc"
	}
	if sortBy == "" {
		sortBy = "published_at"
	}
	// add parameter to context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	ctx = context.WithValue(ctx, "typeBy", typeBy)
	r = r.WithContext(ctx)

	// get all tools
	toolRepo := repository.NewToolRepository(c.DB)
	tools := toolRepo.FindAll(r.Context())

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    tools,
	}

	resp.ToJson(w)
}

func (c *ToolController) Store(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req requests.CreateToolRequest

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

	toolRepo := repository.NewToolRepository(c.DB)
	tool, err := toolRepo.Create(req)
	if err != nil {
		exception.InternalServerError(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    201,
		Message: "OK",
		Data: map[string]any{
			"id":         tool.ID,
			"name":       tool.Name,
			"created_at": tool.CreatedAt.Time.In(time.Local).Format(time.RFC822),
		},
	}

	resp.ToJson(w)
}

func (c *ToolController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req requests.UpdateToolRequest

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

	toolRepo := repository.NewToolRepository(c.DB)

	toolId := p.ByName("toolId")
	_, err = toolRepo.FindOneById(toolId)
	if err != nil {
		exception.BadRequestError(w, r, errors.New("record not found"))
		return
	}

	tool, err := toolRepo.Update(req, toolId)
	if err != nil {
		exception.InternalServerError(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    201,
		Message: "OK",
		Data: map[string]any{
			"id":         tool.ID,
			"name":       tool.Name,
			"updated_at": tool.UpdatedAt.Time.In(time.Local).Format(time.RFC822),
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

	toolRepo := repository.NewToolRepository(c.DB)
	rowsAffected, err := toolRepo.Delete(toolId)
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
