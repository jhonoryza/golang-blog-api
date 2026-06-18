package controller

import (
	request "api_blog/delivery/http/request"
	"api_blog/delivery/http/response"
	"api_blog/infrastructure/exception"
	"api_blog/usecase/post"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type PostController struct {
	Usecase *post.PostUsecase
}

func NewPostController(uc *post.PostUsecase) *PostController {
	return &PostController{Usecase: uc}
}

func (c *PostController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	input := post.FindAllInput{
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
	}

	posts := c.Usecase.FindAll(r.Context(), input)
	postResponses := response.NewPostResponsesFromOutputs(posts)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponses,
	}
	resp.ToJson(w)
}

func (c *PostController) IndexAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	input := post.FindAllInput{
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
	}

	posts := c.Usecase.FindAllIncludingUnpublished(r.Context(), input)
	postResponses := response.NewPostResponsesFromOutputs(posts)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponses,
	}
	resp.ToJson(w)
}

func (c *PostController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postSlug := p.ByName("postSlug")

	postDetail := c.Usecase.FindOneById(r.Context(), postSlug)
	postResponse := response.NewPostResponseFromDetail(postDetail)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    postResponse,
	}
	resp.ToJson(w)
}

func (c *PostController) Store(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req request.CreatePostRequest

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

	result, err := c.Usecase.Create(req)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	var createdAt string
	if result.CreatedAt != nil {
		createdAt = result.CreatedAt.In(time.Local).Format(time.RFC822)
	}

	resp := response.ApiResponse{
		Code:    201,
		Message: "OK",
		Data: map[string]any{
			"id":         result.Id,
			"title":      result.Title,
			"slug":       result.Slug,
			"created_at": createdAt,
		},
	}

	resp.ToJson(w)
}

func (c *PostController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req request.UpdatePostRequest
	postSlug := p.ByName("postSlug")
	if postSlug == "" {
		exception.ErrorHandler(w, r, errors.New("postSlug is required"))
		return
	}

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

	_, err = c.Usecase.FindOneBySlug(postSlug)
	if err != nil {
		exception.BadRequestError(w, r, errors.New("record not found"))
		return
	}

	result, err := c.Usecase.Update(req, postSlug)
	if err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	var updatedAt string
	if result.UpdatedAt != nil {
		updatedAt = result.UpdatedAt.In(time.Local).Format(time.RFC822)
	}

	resp := response.ApiResponse{
		Code:    200,
		Message: "OK",
		Data: map[string]any{
			"id":         result.Id,
			"title":      result.Title,
			"slug":       result.Slug,
			"updated_at": updatedAt,
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

	rowsAffected, err := c.Usecase.Delete(postSlug)
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