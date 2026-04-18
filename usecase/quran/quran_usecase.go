package quran

import (
	"api_blog/domain/entities"
	"api_blog/domain/repositories"
	"context"
)

type QuranUsecase struct {
	repo repositories.QuranRepository
}

func NewQuranUsecase(repo repositories.QuranRepository) *QuranUsecase {
	return &QuranUsecase{repo: repo}
}

func (u *QuranUsecase) FindAllSurah(ctx context.Context, input FindAllSurahInput) *[]entities.Quran {
	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)
	return u.repo.FindAllSurah(ctx)
}

func (u *QuranUsecase) FindAllAyahFromSurah(ctx context.Context, input FindAllAyahInput) *[]entities.QuranVerse {
	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)
	ctx = context.WithValue(ctx, "surahId", input.SurahId)
	return u.repo.FindAllAyahFromSurah(ctx)
}