package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"vkspam/database"
	"vkspam/handlers"
	"vkspam/models"
)

const UserContextKey = "user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("jwt_token")

		if len(tokenString) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusUnauthorized,
				false,
				"jwt_token not found.",
			)
			return
		}

		type UserData struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
			jwt.RegisteredClaims
		}

		tokenData, err := jwt.ParseWithClaims(
			tokenString,
			&UserData{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_KEY")), nil
			},
		)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				fmt.Sprintf("Error parse JWT. %s", err.Error()),
			)

			return
		}

		if !tokenData.Valid {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusUnauthorized,
				false,
				fmt.Sprintf("Token is invalid."),
			)

			return
		}

		claims := tokenData.Claims.(*UserData)

		if claims.Id == 0 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				"User check error.",
			)

			return
		}

		user, err := GetUserById(claims.Id)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				fmt.Sprintf(err.Error()),
			)

			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserById(id int) (*models.User, error) {
	db, _ := database.GetDBInstance()

	rows, err := db.Db.Query("SELECT id, email FROM users WHERE id = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var count int
	var user models.User

	for rows.Next() {
		count++
		err = rows.Scan(&user.Id, &user.Email)
		if err != nil {
			return nil, err
		}
	}

	if count > 1 {
		return nil, errors.New("logic error. More than one user. Admin, Wtf")
	}

	if count < 1 {
		return nil, errors.New("logic error. No user")
	}

	return &user, nil
}

func GetUserFromContext(ctx context.Context) models.User {
	return *ctx.Value(UserContextKey).(*models.User)
}
