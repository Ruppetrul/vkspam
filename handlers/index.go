package handlers

import (
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		log.Println("Error index page:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
