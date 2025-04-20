package controller

import (
	"api_blog/exception"
	"api_blog/helper"
	"api_blog/repository"
	"api_blog/response"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	DB *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{
		DB: db,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

// Login endpoint
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	userRepo := repository.NewUserRepository(c.DB)
	user, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		exception.UnauthorizedError(w, r, errors.New("Invalid credentials"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		exception.UnauthorizedError(w, r, errors.New("Invalid credentials"))
		return
	}

	token, err := helper.GenerateJWT(1, req.Email)
	if err != nil {
		exception.InternalServerError(w, r, errors.New("Could not generate token"))
		return
	}

	resp := response.ApiResponse{
		Code:    200,
		Message: "OK",
		Data: map[string]any{
			"token": token,
		},
	}

	resp.ToJson(w)
}
