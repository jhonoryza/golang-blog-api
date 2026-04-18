package response

import (
	"api_blog/domain/entities"
	"api_blog/domain/services"
	"api_blog/usecase/post"
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

func NewPostResponses(posts *[]entities.Post) *[]PostResponses {
	var postResponses []PostResponses
	for _, post := range *posts {
		summary := services.CastNilString(post.Summary)
		slug := services.CastNilString(post.Slug)
		authorName := services.CastNilString(post.AuthorName)
		categoriesName := services.CastNilString(post.CategoriesName)
		imageUrl := os.Getenv("IMAGE_BASE_URL") + "/blog/laravelblog/storage/" + post.ImageUrl
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

func NewPostResponsesFromOutputs(posts *[]post.PostOutput) *[]PostResponses {
	if posts == nil {
		return nil
	}
	var postResponses []PostResponses
	for _, p := range *posts {
		summary := services.CastNilString(p.Summary)
		slug := services.CastNilString(p.Slug)
		authorName := services.CastNilString(p.AuthorName)
		categoriesName := services.CastNilString(p.CategoriesName)
		imageUrl := os.Getenv("IMAGE_BASE_URL") + "/blog/laravelblog/storage/" + p.ImageUrl
		publishedAt := ""
		if p.PublishedAt != nil {
			publishedAt = p.PublishedAt.In(time.Local).Format(time.RFC822)
		}
		postResponses = append(postResponses, PostResponses{
			Id:             p.Id,
			Title:          p.Title,
			ImageUrl:       imageUrl,
			PublishedAt:    publishedAt,
			Summary:        summary,
			Slug:           slug,
			IsHighlighted:  p.IsHighlighted,
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

func NewPostResponse(post *entities.Post) *PostResponse {
	summary := services.CastNilString(post.Summary)
	slug := services.CastNilString(post.Slug)
	authorName := services.CastNilString(post.AuthorName)
	categoriesName := services.CastNilString(post.CategoriesName)
	imageUrl := os.Getenv("IMAGE_BASE_URL") + "/blog/laravelblog/storage/" + post.ImageUrl
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

func NewPostResponseFromDetail(postDetail *post.PostDetailOutput) *PostResponse {
	if postDetail == nil {
		return nil
	}
	summary := services.CastNilString(postDetail.Summary)
	slug := services.CastNilString(postDetail.Slug)
	authorName := services.CastNilString(postDetail.AuthorName)
	categoriesName := services.CastNilString(postDetail.CategoriesName)
	imageUrl := os.Getenv("IMAGE_BASE_URL") + "/blog/laravelblog/storage/" + postDetail.ImageUrl

	createdAt, updatedAt, publishedAt := "", "", ""
	if postDetail.CreatedAt != nil {
		createdAt = postDetail.CreatedAt.In(time.Local).Format(time.RFC822)
	}
	if postDetail.UpdatedAt != nil {
		updatedAt = postDetail.UpdatedAt.In(time.Local).Format(time.RFC822)
	}
	if postDetail.PublishedAt != nil {
		publishedAt = postDetail.PublishedAt.In(time.Local).Format(time.RFC822)
	}

	return &PostResponse{
		Id:             postDetail.Id,
		AuthorId:       postDetail.AuthorId,
		Title:          postDetail.Title,
		Content:        postDetail.Content,
		ImageUrl:       imageUrl,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		PublishedAt:    publishedAt,
		Summary:        summary,
		Slug:           slug,
		IsMarkdown:     postDetail.IsMarkdown,
		IsHighlighted:  postDetail.IsHighlighted,
		AuthorName:     authorName,
		CategoriesName: categoriesName,
	}
}