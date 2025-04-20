package controller

import (
	"api_blog/exception"
	"api_blog/repository"
	"api_blog/response"
	"context"
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type ToolController struct {
	db *sql.DB
}

func NewToolController(db *sql.DB) *ToolController {
	return &ToolController{db: db}
}

type IndexParam struct {
	Type string `validate:"required"`
}

func (tc *ToolController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	tools := repository.FindAllTools(r.Context(), tc.db)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    tools,
	}

	resp.ToJson(w)
}
