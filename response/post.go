package response

import (
	"api_blog/entity"
	"api_blog/helper"
	"time"
)

type PostResponses struct {
	Id                                          int
	Title, ImageUrl, PublishedAt, Summary, Slug string
	IsHighlighted                               bool
	AuthorName, CategoriesName                  string
}

func NewPostResponses(posts *[]entity.Post) *[]PostResponses {
	var postResponses []PostResponses
	for _, post := range *posts {
		summary := helper.CastNilString(post.Summary)
		slug := helper.CastNilString(post.Slug)
		authorName := helper.CastNilString(post.AuthorName)
		categoriesName := helper.CastNilString(post.CategoriesName)
		postResponses = append(postResponses, PostResponses{
			Id:             post.Id,
			Title:          post.Title,
			ImageUrl:       post.ImageUrl,
			PublishedAt:    post.PublishedAt.Time.In(time.Local).Format(time.RFC822),
			Summary:        summary,
			Slug:           slug,
			IsHighlighted:  post.IsHighlighted,
			AuthorName:     authorName,
			CategoriesName: categoriesName,
		})
	}
	return &postResponses
}

type PostResponse struct {
	Id, AuthorId                                                               int
	Title, Content, ImageUrl, CreatedAt, UpdatedAt, PublishedAt, Summary, Slug string
	IsMarkdown, IsHighlighted                                                  bool
	AuthorName, CategoriesName                                                 string
}

func NewPostResponse(post *entity.Post) *PostResponse {
	summary := helper.CastNilString(post.Summary)
	slug := helper.CastNilString(post.Slug)
	authorName := helper.CastNilString(post.AuthorName)
	categoriesName := helper.CastNilString(post.CategoriesName)
	return &PostResponse{
		Id:             post.Id,
		AuthorId:       post.AuthorId,
		Title:          post.Title,
		Content:        post.Content,
		ImageUrl:       post.ImageUrl,
		CreatedAt:      post.CreatedAt.Time.In(time.Local).Format(time.RFC822),
		UpdatedAt:      post.UpdatedAt.Time.In(time.Local).Format(time.RFC822),
		PublishedAt:    post.PublishedAt.Time.In(time.Local).Format(time.RFC822),
		Summary:        summary,
		Slug:           slug,
		IsMarkdown:     post.IsMarkdown,
		IsHighlighted:  post.IsHighlighted,
		AuthorName:     authorName,
		CategoriesName: categoriesName,
	}
}
