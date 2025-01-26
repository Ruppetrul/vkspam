package distributions

import (
	"log"
	"net/http"
)

func Group(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, err := w.Write([]byte("There will be group list."))
		if err != nil {
			log.Println("Error index page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
