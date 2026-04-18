package user

import (
	"api_blog/domain/entities"
	"api_blog/domain/repositories"
)

type UserUsecase struct {
	repo repositories.UserRepository
}

func NewUserUsecase(repo repositories.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) FindByEmail(email string) *entities.User {
	return u.repo.FindByEmail(email)
}

func (u *UserUsecase) FindByID(id int) *entities.User {
	return u.repo.FindByID(id)
}