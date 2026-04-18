package entities

import "database/sql"

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	ImageUrl    string
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}