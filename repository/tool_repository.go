package repository

import (
	"api_blog/entity"
	"api_blog/exception"
	"api_blog/requests"
	"context"
	"database/sql"
	"time"
)

type ToolRepository struct {
	DB *sql.DB
}

func NewToolRepository(db *sql.DB) *ToolRepository {
	return &ToolRepository{DB: db}
}

func (r *ToolRepository) FindAll(ctx context.Context) *[]entity.Tool {
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

	rows, err := r.DB.QueryContext(ctx, query, args...)

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

func (r *ToolRepository) Create(req requests.CreateToolRequest) (*entity.Tool, error) {
	query := `
		insert into tools(name, "desc", link, is_published, type, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7)
		returning id, name, created_at
	`
	// transaction
	tx, err := r.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

	args := []any{}
	args = append(args, req.Name)
	args = append(args, req.Description)
	args = append(args, req.Link)
	args = append(args, req.IsPublished)
	args = append(args, req.Type)
	args = append(args, time.Now())
	args = append(args, time.Now())

	var tool entity.Tool
	err = tx.QueryRow(query, args...).Scan(&tool.ID, &tool.Name, &tool.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

func (r *ToolRepository) Update(req requests.UpdateToolRequest, toolId string) (*entity.Tool, error) {
	query := `
		update tools
		set name = $1, "desc" = $2, link = $3, is_published = $4, type = $5, updated_at = $6
		where id = $7
		returning id, name, updated_at
	`
	// transaction
	tx, err := r.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

	args := []any{}
	args = append(args, req.Name)
	args = append(args, req.Description)
	args = append(args, req.Link)
	args = append(args, req.IsPublished)
	args = append(args, req.Type)
	args = append(args, time.Now())
	args = append(args, toolId)

	var tool entity.Tool
	err = tx.QueryRow(query, args...).Scan(&tool.ID, &tool.Name, &tool.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

func (r *ToolRepository) FindOneById(toolId string) (*entity.Tool, error) {
	query := `
		select id
		from tools
		where id = $1
	`

	var tool entity.Tool
	err := r.DB.QueryRow(query, toolId).Scan(&tool.ID)
	if err != nil {
		return nil, err
	}

	return &tool, nil
}

func (r *ToolRepository) Delete(toolId string) (int64, error) {
	query := `
		delete from tools
		where id = $1
	`

	// transaction
	tx, err := r.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

	result, err := tx.Exec(query, toolId)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}
	return rowsAffected, nil
}
