package repository

import (
	"api_blog/entity"
	"api_blog/exception"
	"context"
	"database/sql"
)

func FindById(ctx context.Context, tx *sql.Tx, postSlug *string) *entity.Post {
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

func FindAll(ctx context.Context, tx *sql.Tx) *[]entity.Post {
	query := `
	select posts.id, title, summary, posts.slug, posts.published_at, author_id,
       is_highlighted, image_url, users.name, string_agg(categories.name, ',') as categories_name
	from posts
			 left join users on posts.author_id = users.id
			 left join post_categories on posts.id = post_categories.post_id
			 left join categories on post_categories.category_id = categories.id
	where posts.published_at is not null 
	group by posts.id, title, summary, posts.slug, posts.published_at, author_id, is_highlighted, image_url, users.name
	order by published_at desc
			  `
	rows, err := tx.QueryContext(ctx, query)
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
