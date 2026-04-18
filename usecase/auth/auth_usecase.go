package auth

import (
	"api_blog/domain/entities"
	"api_blog/domain/repositories"
	"errors"
)

type AuthUsecase struct {
	userRepo repositories.UserRepository
}

func NewAuthUsecase(userRepo repositories.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (u *AuthUsecase) Login(email, password string) (*entities.User, error) {
	user := u.userRepo.FindByEmail(email)
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (u *AuthUsecase) FindByID(id int) *entities.User {
	return u.userRepo.FindByID(id)
}