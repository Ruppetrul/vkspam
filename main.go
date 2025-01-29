package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"vkspam/database"
	"vkspam/handlers"
	"vkspam/handlers/distributions"
	"vkspam/middleware"
)

func main() {
	checkEnv()
	database.CheckAndMigrate()

	distributionGroupsHandler := distributions.NewDistributionGroupHandler()

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/distributions/group", middleware.AuthMiddleware(distributionGroupsHandler.Group))
	http.HandleFunc("/distributions/group/list", middleware.AuthMiddleware(distributionGroupsHandler.List))

	distributionsHandler := distributions.DistributionHandler{}
	http.HandleFunc("/distribution", middleware.AuthMiddleware(distributionsHandler.Distribution))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}

func checkEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var env string
	if env = os.Getenv("APP_ENV"); env == "" {
		log.Fatal("APP_ENV not found.")
	}
}
