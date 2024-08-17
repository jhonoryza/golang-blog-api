package controller

import (
	"api_blog/exception"
	"api_blog/repository"
	"api_blog/response"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"net/http"
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

	// get all posts
	posts := repository.FindAll(r.Context(), tx)
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
	post := repository.FindById(r.Context(), tx, &postSlug)
	postResponse := response.NewPostResponse(post)

	// return response
	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponse,
	}
	resp.ToJson(w)
}
