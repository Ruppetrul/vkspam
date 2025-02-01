package auth

import (
	"encoding/json"
	"net/http"
	"vkspam/database"
	"vkspam/handlers"
	"vkspam/repositories"
	"vkspam/services"
)

type LoginHandler struct {
	service services.UserService
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func NewLoginHandler() *LoginHandler {
	db, _ := database.GetDBInstance()

	repository := repositories.NewUserRepository(db.Db)
	service := services.NewUserService(repository)

	return &LoginHandler{
		service: service,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusMethodNotAllowed,
			false,
			"Only POST method allowed",
		)

		return
	}

	email := r.FormValue("email")
	if len(email) == 0 {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusBadRequest,
			false,
			"Missing required parameter 'id'",
		)
		return
	}

	password := r.FormValue("password")
	if len(password) == 0 {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusBadRequest,
			false,
			"Missing required parameter 'password'",
		)
		return
	}

	token, err := h.service.TryLogin(email, password)
	if err != nil {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusBadRequest,
			false,
			err.Error(),
		)
		return
	}

	loginResponse := LoginResponse{
		Success: true,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	err = json.NewEncoder(w).Encode(loginResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
