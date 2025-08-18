package routes

import (
	"log"
	"os"
	"strconv"

	"github.com/buildyow/byow-user-service/delivery/http"
	"github.com/buildyow/byow-user-service/delivery/queue"
	"github.com/buildyow/byow-user-service/docs"
	"github.com/buildyow/byow-user-service/infrastructure/db"
	"github.com/buildyow/byow-user-service/infrastructure/jwt"
	loggerZap "github.com/buildyow/byow-user-service/infrastructure/logger"
	"github.com/buildyow/byow-user-service/infrastructure/validation"
	"github.com/buildyow/byow-user-service/repository"
	"github.com/buildyow/byow-user-service/usecase"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.Engine, conn *amqp091.Connection, ch *amqp091.Channel) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	defer logger.Sync()
	r.Use(ginzap.Ginzap(logger, "", true))      // Logging request
	r.Use(ginzap.RecoveryWithZap(logger, true)) // Logging panic recovery
	r.Use(loggerZap.LogRequestBody(logger))     // Logging request body
	// Connect DB
	client, err := db.Connect(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	database := client.Database(os.Getenv("DB_NAME"))
	userRepo := repository.NewUserMongoRepo(database, ch)
	companyRepo := repository.NewCompanyMongoRepo(database)

	// Initialize database indexes
	if err := db.CreateIndexes(database, logger); err != nil {
		logger.Warn("Failed to create database indexes", zap.Error(err))
	}

	// Initialize JWT blacklist service
	blacklistService := jwt.NewBlacklistService(database, logger)
	blacklistService.StartCleanupWorker()

	// Usecase
	userUC := &usecase.UserUsecase{
		Repo:        userRepo,
		CompanyRepo: companyRepo,
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	userUC.JWTExpire, _ = strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	userUC.EmailConfig.Host = os.Getenv("EMAIL_HOST")
	userUC.EmailConfig.Port, _ = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	userUC.EmailConfig.User = os.Getenv("EMAIL_USER")
	userUC.EmailConfig.Pass = os.Getenv("EMAIL_PASS")

	//Consume RabbitMq
	consumer := queue.NewOTPConsumer(*userUC)
	msgs, err := ch.Consume(
		"otp_email_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}
	go consumer.ConsumeOTP(msgs)

	companyUC := &usecase.CompanyUsecase{
		Repo: companyRepo,
		UserID: func(c *gin.Context) string {
			userID, exists := c.Get("user_id")
			if !exists {
				return ""
			}
			if userIDStr, ok := userID.(string); ok {
				return userIDStr
			}
			return ""
		},
	}

	// Handler
	userHandler := http.NewUserHandler(userUC, companyUC)
	companyHandler := http.NewCompanyHandler(companyUC)

	// Public Routes
	auth := r.Group("/auth/users")
	{
		auth.POST("/register",
			validation.ValidateRegistrationRequest(),
			validation.ValidateFileUpload(10<<20, []string{"image/jpeg", "image/png", "image/gif"}), // 10MB limit
			userHandler.Register)
		auth.POST("/login",
			validation.ValidateLoginRequest(),
			userHandler.Login)
		auth.POST("/change-password-otp", userHandler.ChangePasswordWithOTP)
		auth.GET("/forgot-password/send-otp", userHandler.SendOTPForgotPassword)
	}
	verification := r.Group("/verification/users")
	{
		verification.GET("/send-otp", userHandler.SendOTPVerification)
		verification.POST("/verify-otp", userHandler.VerifyOTP)
	}

	// Protected Routes
	protected := r.Group("/api")
	protected.Use(jwt.JWTMiddleware(blacklistService))
	usersPath := protected.Group("/users")
	{
		usersPath.GET("/me", userHandler.UserMe)
		usersPath.GET("/onboard", userHandler.OnBoard)
		usersPath.POST("/update", userHandler.UpdateUser)
		usersPath.POST("/logout", userHandler.Logout)
		usersPath.POST("/change-email", userHandler.ChangeEmail)
		usersPath.GET("/change-email/send-otp", userHandler.SendOTPEmailChange)
		usersPath.POST("/change-phone", userHandler.ChangePhone)
		usersPath.GET("/change-phone/send-otp", userHandler.SendOTPPhoneChange)
		usersPath.POST("/change-password-old", userHandler.ChangePasswordWithOldPassword)
		usersPath.GET("/set-company/:company_id", userHandler.SetCompanyID)
	}
	companyPath := protected.Group("/companies")
	{
		companyPath.GET("/all", companyHandler.FindAll)
		companyPath.POST("/create", companyHandler.Create)
		companyPath.GET("/:id", companyHandler.FindByID)
		companyPath.PUT("/:id", companyHandler.Update)
		companyPath.DELETE("/:id", companyHandler.DeleteById)
	}

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "BYOW User Service is healthy",
			"version": "1.0.0",
		})
	})

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
