package response

import (
	"api_blog/entity"
	"api_blog/helper"
	"os"
	"time"
)

type PostResponses struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	ImageUrl       string `json:"image_url"`
	PublishedAt    string `json:"published_at"`
	Summary        string `json:"summary"`
	Slug           string `json:"slug"`
	IsHighlighted  bool   `json:"is_highlighted"`
	AuthorName     string `json:"author_name"`
	CategoriesName string `json:"categories_name"`
}

func NewPostResponses(posts *[]entity.Post) *[]PostResponses {
	var postResponses []PostResponses
	for _, post := range *posts {
		summary := helper.CastNilString(post.Summary)
		slug := helper.CastNilString(post.Slug)
		authorName := helper.CastNilString(post.AuthorName)
		categoriesName := helper.CastNilString(post.CategoriesName)
		imageUrl := os.Getenv("IMAGE_BASE_URL") + post.ImageUrl
		postResponses = append(postResponses, PostResponses{
			Id:             post.Id,
			Title:          post.Title,
			ImageUrl:       imageUrl,
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
	Id             int    `json:"id"`
	AuthorId       int    `json:"author_id"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	ImageUrl       string `json:"image_url"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	PublishedAt    string `json:"published_at"`
	Summary        string `json:"summary"`
	Slug           string `json:"slug"`
	IsMarkdown     bool   `json:"is_markdown"`
	IsHighlighted  bool   `json:"is_highlighted"`
	AuthorName     string `json:"author_name"`
	CategoriesName string `json:"categories_name"`
}

func NewPostResponse(post *entity.Post) *PostResponse {
	summary := helper.CastNilString(post.Summary)
	slug := helper.CastNilString(post.Slug)
	authorName := helper.CastNilString(post.AuthorName)
	categoriesName := helper.CastNilString(post.CategoriesName)
	imageUrl := os.Getenv("IMAGE_BASE_URL") + post.ImageUrl
	return &PostResponse{
		Id:             post.Id,
		AuthorId:       post.AuthorId,
		Title:          post.Title,
		Content:        post.Content,
		ImageUrl:       imageUrl,
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
