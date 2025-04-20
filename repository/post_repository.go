package repository

import (
	"api_blog/entity"
	"api_blog/exception"
	"context"
	"database/sql"
)

func FindOnePostById(ctx context.Context, tx *sql.Tx, postSlug *string) *entity.Post {
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

	row := tx.QueryRowContext(ctx, query, *postSlug)
	exception.PanicNotFoundIfErr(row.Err())

	var post entity.Post
	err := row.Scan(&post.Id, &post.Title, &post.Summary, &post.Content, &post.Slug, &post.PublishedAt, &post.AuthorId, &post.CreatedAt, &post.UpdatedAt, &post.IsMarkdown, &post.IsHighlighted, &post.ImageUrl, &post.AuthorName, &post.CategoriesName)
	exception.PanicNotFoundIfErr(err)
	return &post
}

func FindAllPosts(ctx context.Context, tx *sql.Tx) *[]entity.Post {
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
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = tx.QueryContext(ctx, query)
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
