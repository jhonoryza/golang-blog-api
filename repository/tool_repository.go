package repository

import (
	"api_blog/entity"
	"api_blog/exception"
	"context"
	"database/sql"
)

func FindAllTools(ctx context.Context, db *sql.DB) *[]entity.Tool {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)
	typeBy := ctx.Value("typeBy").(string)

	allowedSortBy := map[string]bool{
		"name": true,
		"id":   true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values kalau input nggak valid
	if !allowedSortBy[sortBy] {
		sortBy = "id"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "desc"
	}

	args := []any{}
	query := `SELECT id, name, tools.desc, link, is_published, type, created_at, updated_at
	FROM tools
	where is_published = true
	and type = $1`

	args = append(args, typeBy)

	if search != "" {
		query += ` AND (lower(name) LIKE lower($2) OR lower(link) LIKE lower($2) OR lower(tools.desc) LIKE lower($2))`
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY " + sortBy + " " + sortDir

	rows, err := db.QueryContext(ctx, query, args...)

	exception.PanicIfErr(err)
	defer rows.Close()

	var tools []entity.Tool
	for rows.Next() {
		var tool entity.Tool
		err := rows.Scan(&tool.ID, &tool.Name, &tool.Description, &tool.Link, &tool.IsPublished, &tool.Type, &tool.CreatedAt, &tool.UpdatedAt)
		exception.PanicIfErr(err)
		tools = append(tools, tool)
	}

	return &tools
}
