package routes

import (
	"os"
	"strconv"

	"github.com/buildyow/byow-user-service/delivery/http"
	"github.com/buildyow/byow-user-service/docs"
	"github.com/buildyow/byow-user-service/infrastructure/db"
	"github.com/buildyow/byow-user-service/infrastructure/jwt"
	loggerZap "github.com/buildyow/byow-user-service/infrastructure/logger"
	"github.com/buildyow/byow-user-service/repository"
	"github.com/buildyow/byow-user-service/usecase"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.Engine) {
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
	userRepo := repository.NewUserMongoRepo(client.Database(os.Getenv("DB_NAME")))

	// Usecase
	userUC := &usecase.UserUsecase{
		Repo:      userRepo,
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
	userUC.JWTExpire, _ = strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	userUC.EmailConfig.Host = os.Getenv("EMAIL_HOST")
	userUC.EmailConfig.Port, _ = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	userUC.EmailConfig.User = os.Getenv("EMAIL_USER")
	userUC.EmailConfig.Pass = os.Getenv("EMAIL_PASS")

	// Handler
	userHandler := http.NewUserHandler(userUC)

	// Public Routes
	auth := r.Group("/auth/users")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
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
	protected.Use(jwt.JWTMiddleware())
	{
		protected.GET("/users/me", userHandler.UserMe)
		protected.GET("/users/onboard", userHandler.OnBoard)
		protected.POST("/users/update", userHandler.UpdateUser)
		protected.POST("/users/logout", userHandler.Logout)
		protected.POST("/users/change-email", userHandler.ChangeEmail)
		protected.GET("/users/change-email/send-otp", userHandler.SendOTPEmailChange)
		protected.POST("/users/change-phone", userHandler.ChangePhone)
		protected.GET("/users/change-phone/send-otp", userHandler.SendOTPPhoneChange)
		protected.POST("/users/change-password-old", userHandler.ChangePasswordWithOldPassword)
	}

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
