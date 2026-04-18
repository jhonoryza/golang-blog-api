package entities

import "database/sql"

type Service struct {
	ID          int
	Name        string
	Description string
	Icon        string
	Link        string
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}