package repository

import (
	"api_blog/entity"
	"api_blog/exception"
	"api_blog/helper"
	"api_blog/requests"
	"context"
	"database/sql"
	"time"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) FindOneBySlug(postSlug string) (*entity.Post, error) {
	query := `
		select id, title, slug
		from posts
		where slug = $1
	`
	var post entity.Post
	err := r.DB.QueryRow(query, postSlug).Scan(&post.Id, &post.Title, &post.Slug)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) FindOneById(ctx context.Context, postSlug *string) *entity.Post {
	query := `
	select posts.id, title, summary, content, posts.slug, posts.published_at, author_id, posts.created_at, posts.updated_at,
       is_markdown, is_highlighted, image_url, users.name, string_agg(categories.name, ',') as categories_name
	from posts
			 left join users on posts.author_id = users.id
			 left join post_categories on posts.id = post_categories.post_id
			 left join categories on post_categories.category_id = categories.id
	where posts.slug = $1 and posts.published_at is not null
	group by posts.id, title, summary, content, posts.slug, posts.published_at, author_id, posts.created_at, posts.updated_at, is_markdown, is_highlighted, image_url, users.name
	`

	row := r.DB.QueryRowContext(ctx, query, *postSlug)
	exception.PanicNotFoundIfErr(row.Err())

	var post entity.Post
	err := row.Scan(&post.Id, &post.Title, &post.Summary, &post.Content, &post.Slug, &post.PublishedAt, &post.AuthorId, &post.CreatedAt, &post.UpdatedAt, &post.IsMarkdown, &post.IsHighlighted, &post.ImageUrl, &post.AuthorName, &post.CategoriesName)
	exception.PanicNotFoundIfErr(err)
	return &post
}

func (r *PostRepository) FindAll(ctx context.Context) *[]entity.Post {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)

	allowedSortBy := map[string]bool{
		"title":        true,
		"published_at": true,
		"author_id":    true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values kalau input nggak valid
	if !allowedSortBy[sortBy] {
		sortBy = "published_at"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "desc"
	}

	query := `
	select posts.id, title, summary, posts.slug, posts.published_at, author_id,
							is_highlighted, image_url, users.name, string_agg(categories.name, ',') as categories_name
	from posts
				left join users on posts.author_id = users.id
				left join post_categories on posts.id = post_categories.post_id
				left join categories on post_categories.category_id = categories.id
	where posts.published_at is not null
	`

	args := []any{}

	if search != "" {
		query += ` AND (lower(title) LIKE lower($1) OR lower(summary) LIKE lower($1))`
		args = append(args, "%"+search+"%")
	}

	query += `
	group by posts.id, title, summary, posts.slug, posts.published_at, author_id, is_highlighted, image_url, users.name
	order by ` + sortBy + ` ` + sortDir + ``

	var rows *sql.Rows
	var err error

	if len(args) > 0 {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = r.DB.QueryContext(ctx, query)
	}

	exception.PanicIfErr(err)
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Summary, &post.Slug, &post.PublishedAt, &post.AuthorId, &post.IsHighlighted, &post.ImageUrl, &post.AuthorName, &post.CategoriesName)
		exception.PanicIfErr(err)
		posts = append(posts, post)
	}

	return &posts
}

func (r *PostRepository) Create(req requests.CreatePostRequest) (*entity.Post, error) {
	query := `
		insert into posts (title, summary, content, slug, published_at, created_at, updated_at,
	 	author_id, is_markdown, image_url, image_tw_url, image_thumb_url)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		returning id, title, slug, created_at
	`

	var publishedAt *time.Time
	if req.PublishedAt != nil {
		parsedTime, err := time.Parse(time.DateTime, *req.PublishedAt)
		if err != nil {
			return nil, err
		}
		publishedAt = &parsedTime
	}

	slug := helper.GenerateSlug(req.Title)

	var post entity.Post

	args := []any{}
	args = append(args, req.Title)
	args = append(args, req.Summary)
	args = append(args, req.Content)
	args = append(args, slug)
	args = append(args, publishedAt)
	args = append(args, time.Now())
	args = append(args, time.Now())
	args = append(args, req.AuthorId)
	args = append(args, req.IsMarkdown)
	args = append(args, req.ImageUrl)
	args = append(args, req.ImageTwUrl)
	args = append(args, req.ImageThumbUrl)

	// transaction
	tx, err := r.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

	err = tx.QueryRow(query, args...).Scan(&post.Id, &post.Title, &post.Slug, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Update(req requests.UpdatePostRequest, postSlug string) (*entity.Post, error) {
	query := `
		update posts
		set title = $1, slug = $2, summary = $3, content = $4, published_at = $5, updated_at = $6,
	 	author_id = $7, is_markdown = $8, image_url = $9, image_tw_url = $10, image_thumb_url = $11
		where slug = $12
		returning id, title, slug, updated_at
	`

	var publishedAt *time.Time
	if req.PublishedAt != nil {
		parsedTime, err := time.Parse(time.DateTime, *req.PublishedAt)
		if err != nil {
			return nil, err
		}
		publishedAt = &parsedTime
	}

	var slug *string = req.Slug

	if slug == nil {
		newSlug := helper.GenerateSlug(req.Title)
		slug = &newSlug
	}

	var post entity.Post

	args := []any{}
	args = append(args, req.Title)
	args = append(args, slug)
	args = append(args, req.Slug)
	args = append(args, req.Content)
	args = append(args, publishedAt)
	args = append(args, time.Now())
	args = append(args, req.AuthorId)
	args = append(args, req.IsMarkdown)
	args = append(args, req.ImageUrl)
	args = append(args, req.ImageTwUrl)
	args = append(args, req.ImageThumbUrl)
	args = append(args, postSlug)

	// transaction
	tx, err := r.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

	err = tx.QueryRow(query, args...).Scan(&post.Id, &post.Title, &post.Slug, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Delete(postSlug string) (int64, error) {
	query := `
		delete from posts
		where slug = $1
	`
	// transaction
	tx, err := r.DB.Begin()
	exception.PanicIfErr(err)
	defer exception.CommitOrRollback(tx)

	result, err := tx.Exec(query, postSlug)
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
