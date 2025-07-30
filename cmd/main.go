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

func main() {
	_ = godotenv.Load()

	r := gin.Default()
	r.Use(corsService.SetupCors())
	routes.InitRoutes(r)

	port := os.Getenv("PORT")
	log.Println("Running on port", port)
	log.Fatal(r.Run(":" + port))
}
