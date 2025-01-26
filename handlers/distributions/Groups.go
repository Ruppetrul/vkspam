package distributions

import (
	"log"
	"net/http"
	"vkspam/middleware"
)

func Group(w http.ResponseWriter, r *http.Request) {
	middleware.GetUserFromContext(r.Context())

	if r.Method == http.MethodGet {
		_, err := w.Write([]byte("There will be group list."))
		if err != nil {
			log.Println("Error index page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		//TODO create
	}

	if r.Method == http.MethodPut {
		//TODO update
	}

	if r.Method == http.MethodDelete {
		//TODO delete
	}
}
