package repositories

import (
	"api_blog/domain/entities"
	request "api_blog/delivery/http/request"
	"context"
)

type PostRepository interface {
	FindAll(ctx context.Context) *[]entities.Post
	FindOneById(ctx context.Context, postSlug *string) *entities.Post
	FindOneBySlug(postSlug string) (*entities.Post, error)
	Create(req request.CreatePostRequest) (*entities.Post, error)
	Update(req request.UpdatePostRequest, postSlug string) (*entities.Post, error)
	Delete(postSlug string) (int64, error)
}

type UserRepository interface {
	FindByEmail(email string) *entities.User
	FindByID(id int) *entities.User
	Create(user entities.User) (*entities.User, error)
}

type ToolRepository interface {
	FindAll(ctx context.Context) *[]entities.Tool
	FindOneById(toolId string) *entities.Tool
	Create(req request.CreateToolRequest) (*entities.Tool, error)
	Update(req request.UpdateToolRequest, toolId string) (*entities.Tool, error)
	Delete(toolId string) (int64, error)
}

type QuranRepository interface {
	FindAllSurah(ctx context.Context) *[]entities.Quran
	FindAllAyahFromSurah(ctx context.Context) *[]entities.QuranVerse
}

type WilayahRepository interface {
	FindAllProvinces(ctx context.Context) *[]entities.Province
	FindCitiesByProvince(ctx context.Context, provinceId string) *[]entities.City
	FindDistrictsByCity(ctx context.Context, cityId string) *[]entities.District
	FindSubdistrictsByDistrict(ctx context.Context, districtId string) *[]entities.Subdistrict
}

type ProductRepository interface {
	FindAll() *[]entities.Product
	FindOneById(productId string) *entities.Product
}

type ServiceRepository interface {
	FindAll() *[]entities.Service
}