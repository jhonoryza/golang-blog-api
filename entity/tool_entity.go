package entity

import (
	"database/sql"
)

type Tool struct {
	ID          int
	Name        string
	Description *string
	Link        string
	IsPublished bool
	Type        string // dev tools,php packages,go packages,tutorial
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}
