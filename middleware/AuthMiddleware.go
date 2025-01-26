package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"vkspam/database"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		fmt.Print(token)
		if len(token) == 0 {
			_, _ = w.Write([]byte("Token not found."))
			return
		}

		_, err := userCheck(token)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		next.ServeHTTP(w, r)
	}
}

func userCheck(jwtToken string) (success bool, err error) {
	db, _ := database.GetDBInstance()

	rows, err := db.Db.Query("SELECT * FROM users WHERE jwt_token = $1;", jwtToken)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}

	if count > 1 {
		return false, errors.New("Logic error. More than one user. Admin, Wtf?")
	}

	if count < 1 {
		return false, errors.New("Logic error. No user.")
	}

	return true, nil
}
