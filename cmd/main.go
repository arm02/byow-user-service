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
	"os/signal"
	"syscall"

	corsService "github.com/buildyow/byow-user-service/infrastructure/cors"
	"github.com/buildyow/byow-user-service/infrastructure/rabbitmq"
	"github.com/buildyow/byow-user-service/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// setupServer creates and configures the Gin router
func setupServer() *gin.Engine {
	r := gin.Default()
	conn, ch := rabbitmq.Connect("amqp://guest:guest@localhost:5672/")
	r.Use(corsService.SetupCors())
	routes.InitRoutes(r, conn, ch)
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Closing RabbitMQ connection...")
		ch.Close()
		conn.Close()
		os.Exit(0)
	}()

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
