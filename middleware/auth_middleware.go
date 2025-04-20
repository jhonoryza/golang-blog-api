package middleware

import (
	"api_blog/exception"
	"api_blog/helper"
	"api_blog/repository"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func AuthMiddleware(db *sql.DB, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			exception.UnauthorizedError(w, r, errors.New("Unauthorized"))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := helper.ValidateToken(token)
		if err != nil {
			exception.UnauthorizedError(w, r, errors.New("Unauthorized"))
			return
		}

		// Dapatkan user ID dari claims
		userID := claims.UserID

		// Cek ke database lewat UserRepository
		userRepo := repository.NewUserRepository(db)
		user, err := userRepo.FindById(userID)
		if err != nil {
			exception.UnauthorizedError(w, r, errors.New("Unauthorized"))
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next(w, r.WithContext(ctx), ps)
	}
}
