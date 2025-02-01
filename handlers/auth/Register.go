package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vkspam/handlers"
)

func (h *LoginHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusMethodNotAllowed,
			false,
			fmt.Sprintf("Only POST method allowed"),
		)
		return
	}

	email := r.FormValue("email")
	if len(email) < 1 {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusBadRequest,
			false,
			fmt.Sprintf("Missing required parameter 'email'"),
		)

		return
	}
	password := r.FormValue("password")
	if len(password) < 1 {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusBadRequest,
			false,
			fmt.Sprintf("Missing required parameter 'password'"),
		)

		return
	}

	isUserExists, err := h.service.CheckEmailExist(email)

	if err != nil {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusInternalServerError,
			false,
			fmt.Sprintf("Error when check email exists. %s", err.Error()),
		)

		return
	}

	if true == isUserExists {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusBadRequest,
			false,
			fmt.Sprintf("User with this email already exists"),
		)

		return
	}

	user, err := h.service.Register(email, password)
	if err != nil {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusInternalServerError,
			false,
			fmt.Sprintf("Error when register user. %s", err.Error()),
		)

		return
	}

	response := LoginResponse{
		true,
		user.Token,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusInternalServerError,
			false,
			fmt.Sprintf("Error when response after register user. %s", err.Error()),
		)

		return
	}
}
