// @title Build Your Own Website User Service API
// @version 1.0
// @description This is the user service for the Byow app
// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"log"
	"os"

	corsService "github.com/buildyow/byow-user-service/infrastructure/cors"
	"github.com/buildyow/byow-user-service/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// setupServer creates and configures the Gin router
func setupServer() *gin.Engine {
	r := gin.Default()
	r.Use(corsService.SetupCors())
	routes.InitRoutes(r)
	return r
}

// getPort returns the port from environment variable, with fallback to "8080"
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// loadEnv loads the .env file, ignoring errors
func loadEnv() {
	_ = godotenv.Load()
}

func main() {
	loadEnv()

	r := setupServer()
	port := getPort()
	
	log.Println("Running on port", port)
	log.Fatal(r.Run(":" + port))
}
