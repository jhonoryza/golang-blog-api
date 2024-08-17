package entity

import (
	"database/sql"
)

type Post struct {
	Id, AuthorId               int
	Title, Content, ImageUrl   string
	Summary, Slug              *string
	PublishedAt                sql.NullTime
	CreatedAt, UpdatedAt       sql.NullTime
	IsMarkdown, IsHighlighted  bool
	AuthorName, CategoriesName *string
}
