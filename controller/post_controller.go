package controller

import (
	"api_blog/exception"
	"api_blog/repository"
	"api_blog/response"
	"context"
	"database/sql"
	"net/http"

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

func (postController *PostController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// transaction
	tx, err := postController.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

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
	posts := repository.FindAllPosts(r.Context(), tx)
	postResponses := response.NewPostResponses(posts)

	// return response
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponses,
	}
	resp.ToJson(w)
}

func (postController *PostController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postSlug := p.ByName("postSlug")

	// transaction
	tx, err := postController.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

	// get specific posts
	post := repository.FindOnePostById(r.Context(), tx, &postSlug)
	postResponse := response.NewPostResponse(post)

	// return response
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponse,
	}
	resp.ToJson(w)
}
