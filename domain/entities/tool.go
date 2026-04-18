package entities

import "database/sql"

type Tool struct {
	ID          int
	Name        string
	Description *string
	Link        string
	IsPublished bool
	Type        string
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}