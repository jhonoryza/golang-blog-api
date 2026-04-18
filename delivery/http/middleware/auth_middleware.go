package middleware

import (
	"api_blog/infrastructure/exception"
	"api_blog/domain/services"
	"api_blog/usecase/user"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func AuthMiddleware(userUsecase *user.UserUsecase, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			exception.UnauthorizedError(w, r, errors.New("Unauthorized"))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := services.ValidateToken(token)
		if err != nil {
			exception.UnauthorizedError(w, r, errors.New("Unauthorized"))
			return
		}

		userID := claims.UserID

		user := userUsecase.FindByID(userID)
		if user == nil {
			exception.UnauthorizedError(w, r, errors.New("Unauthorized"))
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next(w, r.WithContext(ctx), ps)
	}
}