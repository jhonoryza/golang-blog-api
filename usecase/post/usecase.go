package post

import (
	"api_blog/domain/entities"
	"api_blog/domain/repositories"
	request "api_blog/delivery/http/request"
	"context"
)

type PostUsecase struct {
	repo repositories.PostRepository
}

func NewPostUsecase(repo repositories.PostRepository) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (u *PostUsecase) FindAll(ctx context.Context, input FindAllInput) *[]PostOutput {
	if input.SortDir == "" {
		input.SortDir = "desc"
	}
	if input.SortBy == "" {
		input.SortBy = "published_at"
	}

	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)

	posts := u.repo.FindAll(ctx)
	return ToPostOutputs(posts)
}

func (u *PostUsecase) FindOneById(ctx context.Context, slug string) *PostDetailOutput {
	post := u.repo.FindOneById(ctx, &slug)
	return ToPostDetailOutput(post)
}

func (u *PostUsecase) FindOneBySlug(slug string) (*PostOutput, error) {
	post, err := u.repo.FindOneBySlug(slug)
	if err != nil {
		return nil, err
	}
	result := ToPostOutput(post)
	return &result, nil
}

func (u *PostUsecase) Create(req request.CreatePostRequest) (*CreateOutput, error) {
	post, err := u.repo.Create(req)
	if err != nil {
		return nil, err
	}
	return &CreateOutput{
		Id:        post.Id,
		Title:     post.Title,
		Slug:      post.Slug,
		CreatedAt: &post.CreatedAt.Time,
	}, nil
}

func (u *PostUsecase) Update(req request.UpdatePostRequest, slug string) (*UpdateOutput, error) {
	post, err := u.repo.Update(req, slug)
	if err != nil {
		return nil, err
	}
	return &UpdateOutput{
		Id:        post.Id,
		Title:     post.Title,
		Slug:      post.Slug,
		UpdatedAt: &post.UpdatedAt.Time,
	}, nil
}

func (u *PostUsecase) Delete(slug string) (int64, error) {
	return u.repo.Delete(slug)
}

func ToPostOutputs(posts *[]entities.Post) *[]PostOutput {
	if posts == nil {
		return nil
	}
	var result []PostOutput
	for _, p := range *posts {
		result = append(result, PostOutput{
			Id:             p.Id,
			Title:          p.Title,
			ImageUrl:       p.ImageUrl,
			PublishedAt:    nil,
			Summary:        p.Summary,
			Slug:           p.Slug,
			IsHighlighted:  p.IsHighlighted,
			AuthorName:     p.AuthorName,
			CategoriesName: p.CategoriesName,
		})
	}
	return &result
}

func ToPostOutput(post *entities.Post) PostOutput {
	if post == nil {
		return PostOutput{}
	}
	return PostOutput{
		Id:             post.Id,
		Title:          post.Title,
		ImageUrl:       post.ImageUrl,
		PublishedAt:    nil,
		Summary:        post.Summary,
		Slug:           post.Slug,
		IsHighlighted:  post.IsHighlighted,
		AuthorName:     post.AuthorName,
		CategoriesName: post.CategoriesName,
	}
}

func ToPostDetailOutput(post *entities.Post) *PostDetailOutput {
	if post == nil {
		return nil
	}
	return &PostDetailOutput{
		Id:             post.Id,
		AuthorId:       post.AuthorId,
		Title:          post.Title,
		Content:        post.Content,
		ImageUrl:       post.ImageUrl,
		CreatedAt:      &post.CreatedAt.Time,
		UpdatedAt:      &post.UpdatedAt.Time,
		PublishedAt:    &post.PublishedAt.Time,
		Summary:        post.Summary,
		Slug:           post.Slug,
		IsMarkdown:     post.IsMarkdown,
		IsHighlighted:  post.IsHighlighted,
		AuthorName:     post.AuthorName,
		CategoriesName: post.CategoriesName,
	}
}