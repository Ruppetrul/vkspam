package main

import (
	"log"
	"net/http"
	"vkspam/handlers"
)

func main() {
	http.HandleFunc("/", handlers.Index)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}
