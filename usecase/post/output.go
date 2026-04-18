package post

import "time"

type PostOutput struct {
	Id             int        `json:"id"`
	Title          string     `json:"title"`
	ImageUrl       string     `json:"image_url"`
	PublishedAt    *time.Time `json:"published_at"`
	Summary        *string    `json:"summary"`
	Slug           *string    `json:"slug"`
	IsHighlighted  bool       `json:"is_highlighted"`
	AuthorName     *string    `json:"author_name"`
	CategoriesName *string    `json:"categories_name"`
}

type PostDetailOutput struct {
	Id             int        `json:"id"`
	AuthorId       int        `json:"author_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	ImageUrl       string     `json:"image_url"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	PublishedAt    *time.Time `json:"published_at"`
	Summary        *string    `json:"summary"`
	Slug           *string    `json:"slug"`
	IsMarkdown     bool       `json:"is_markdown"`
	IsHighlighted  bool       `json:"is_highlighted"`
	AuthorName     *string    `json:"author_name"`
	CategoriesName *string    `json:"categories_name"`
}

type CreateOutput struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Slug      *string    `json:"slug"`
	CreatedAt *time.Time `json:"created_at"`
}

type UpdateOutput struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Slug      *string    `json:"slug"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DeleteOutput struct {
	RowsAffected int64 `json:"rowsAffected"`
}