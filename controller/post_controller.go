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

type PostController struct {
	DB *sql.DB
}

func NewPostController(db *sql.DB) *PostController {
	return &PostController{
		DB: db,
	}
}

func (c *PostController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
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
	r = r.WithContext(ctx)

	// get all posts
	postRepo := repository.NewPostRepository(c.DB)
	posts := postRepo.FindAll(r.Context())
	postResponses := response.NewPostResponses(posts)

	// return response
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponses,
	}
	resp.ToJson(w)
}

func (c *PostController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postSlug := p.ByName("postSlug")

	// get specific posts
	postRepo := repository.NewPostRepository(c.DB)
	post := postRepo.FindOneById(r.Context(), &postSlug)
	postResponse := response.NewPostResponse(post)

	// return response
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponse,
	}
	resp.ToJson(w)
}

func (c *PostController) Store(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req requests.CreatePostRequest

	// decode json body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	// validate
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	postRepo := repository.NewPostRepository(c.DB)
	post, err := postRepo.Create(req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    201,
		Message: "OK",
		Data: map[string]any{
			"id":         post.Id,
			"title":      post.Title,
			"created_at": post.CreatedAt.Time.In(time.Local).Format(time.RFC822),
		},
	}

	resp.ToJson(w)
}

func (c *PostController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req requests.UpdatePostRequest
	postSlug := p.ByName("postSlug")
	if postSlug == "" {
		exception.ErrorHandler(w, r, errors.New("postSlug is required"))
		return
	}

	// decode json body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	// validate
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	postRepo := repository.NewPostRepository(c.DB)
	post, err := postRepo.Update(req, postSlug)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	resp := response.ApiResponse{
		Code:    200,
		Message: "OK",
		Data: map[string]any{
			"id":         post.Id,
			"title":      post.Title,
			"slug":       post.Slug,
			"updated_at": post.UpdatedAt.Time.In(time.Local).Format(time.RFC822),
		},
	}

	resp.ToJson(w)
}

func (c *PostController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postSlug := p.ByName("postSlug")
	if postSlug == "" {
		exception.ErrorHandler(w, r, errors.New("postSlug is required"))
		return
	}

	postRepo := repository.NewPostRepository(c.DB)
	rowsAffected, err := postRepo.Delete(postSlug)
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
