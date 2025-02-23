package handlers

import (
	"encoding/json"
	"net/http"
)

type AppBaseResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type AppErrorResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

func ReturnAppBaseResponse(w http.ResponseWriter, code int, success bool, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var response interface{}
	if code >= 200 && code <= 299 {
		response = AppBaseResponse{
			Success: success,
			Data:    data,
		}
	} else {
		response = AppErrorResponse{
			Success: false,
			Message: data,
		}
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
