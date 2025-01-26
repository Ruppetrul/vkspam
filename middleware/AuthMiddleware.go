package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"vkspam/database"
	"vkspam/models"
)

const UserContextKey = "user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		fmt.Print(token)
		if len(token) == 0 {
			_, _ = w.Write([]byte("Token not found."))
			return
		}

		user, err := userCheck(token)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func userCheck(jwtToken string) (*models.User, error) {
	db, _ := database.GetDBInstance()

	rows, err := db.Db.Query("SELECT * FROM users WHERE jwt_token = $1;", jwtToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var count int
	var user models.User

	for rows.Next() {
		count++
		err = rows.Scan(&user.Id, &user.Name, &user.Token, &user.InviteCode)
		if err != nil {
			return nil, err
		}
	}

	if count > 1 {
		return nil, errors.New("Logic error. More than one user. Admin, Wtf?")
	}

	if count < 1 {
		return nil, errors.New("Logic error. No user.")
	}

	return &user, nil
}

func GetUserFromContext(ctx context.Context) *models.User {
	return ctx.Value(UserContextKey).(*models.User)
}
