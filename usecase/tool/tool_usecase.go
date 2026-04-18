package tool

import (
	"api_blog/domain/entities"
	"api_blog/domain/repositories"
	request "api_blog/delivery/http/request"
	"context"
)

type ToolUsecase struct {
	repo repositories.ToolRepository
}

func NewToolUsecase(repo repositories.ToolRepository) *ToolUsecase {
	return &ToolUsecase{repo: repo}
}

func (u *ToolUsecase) FindAll(ctx context.Context, input FindAllInput) *[]entities.Tool {
	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)
	ctx = context.WithValue(ctx, "typeBy", input.Type)
	return u.repo.FindAll(ctx)
}

func (u *ToolUsecase) FindOneById(toolId string) *entities.Tool {
	return u.repo.FindOneById(toolId)
}

func (u *ToolUsecase) Create(req request.CreateToolRequest) (*entities.Tool, error) {
	return u.repo.Create(req)
}

func (u *ToolUsecase) Update(req request.UpdateToolRequest, toolId string) (*entities.Tool, error) {
	return u.repo.Update(req, toolId)
}

func (u *ToolUsecase) Delete(toolId string) (int64, error) {
	return u.repo.Delete(toolId)
}