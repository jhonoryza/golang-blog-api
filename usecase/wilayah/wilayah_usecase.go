package wilayah

import (
	"api_blog/domain/entities"
	"api_blog/domain/repositories"
	"context"
)

type WilayahUsecase struct {
	repo repositories.WilayahRepository
}

func NewWilayahUsecase(repo repositories.WilayahRepository) *WilayahUsecase {
	return &WilayahUsecase{repo: repo}
}

func (u *WilayahUsecase) FindAllProvinces(ctx context.Context, input FindAllProvincesInput) *[]entities.Province {
	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)
	return u.repo.FindAllProvinces(ctx)
}

func (u *WilayahUsecase) FindCitiesByProvince(ctx context.Context, input FindCitiesByProvinceInput) *[]entities.City {
	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)
	return u.repo.FindCitiesByProvince(ctx, input.ProvinceId)
}

func (u *WilayahUsecase) FindDistrictsByCity(ctx context.Context, input FindDistrictsByCityInput) *[]entities.District {
	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)
	return u.repo.FindDistrictsByCity(ctx, input.CityId)
}

func (u *WilayahUsecase) FindSubdistrictsByDistrict(ctx context.Context, input FindSubdistrictsByDistrictInput) *[]entities.Subdistrict {
	ctx = context.WithValue(ctx, "search", input.Search)
	ctx = context.WithValue(ctx, "sortBy", input.SortBy)
	ctx = context.WithValue(ctx, "sortDir", input.SortDir)
	return u.repo.FindSubdistrictsByDistrict(ctx, input.DistrictId)
}