package main

import (
	"log"
	"net/http"
	"os"
	"vkspam/handlers"
)

func main() {
	checkEnv()

	http.HandleFunc("/", handlers.Index)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}

func checkEnv() {
	var env string
	if env = os.Getenv("APP_ENV"); env == "" {
		log.Fatal("APP_ENV not found.")
	}
}
