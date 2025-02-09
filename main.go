package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"vkspam/database"
	"vkspam/handlers"
	"vkspam/handlers/auth"
	"vkspam/handlers/distributions"
	"vkspam/middleware"
)

type AppBaseResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

func main() {
	checkEnv()
	database.CheckAndMigrate()

	distributionGroupsHandler := distributions.NewDistributionGroupHandler()
	authHandler := auth.NewLoginHandler()

	http.HandleFunc("/", handlers.Index)

	http.HandleFunc("/auth/login", authHandler.Login)
	http.HandleFunc("/auth/register", authHandler.Register)

	http.HandleFunc("/distributions/group", middleware.AuthMiddleware(distributionGroupsHandler.Group))
	http.HandleFunc("/distributions/group/list", middleware.AuthMiddleware(distributionGroupsHandler.List))

	http.HandleFunc("/distributions/run", distributionGroupsHandler.Run)
	distributionsHandler := distributions.NewDistributionHandler()
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
