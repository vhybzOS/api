package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vhybZApp/api.git/azure"
	"github.com/vhybZApp/api.git/config"
	"github.com/vhybZApp/api.git/database"
)

// @title           RESTful API with JWT Authentication
// @version         1.0
// @description     A simple RESTful server with JWT authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/refresh", refresh)
		auth.GET("/profile", authMiddleware(), getProfile)
	}

	// Azure OpenAI routes
	azureGroup := r.Group("/azure")
	{
		azureGroup.POST("/chat/completions", authMiddleware(), azure.ChatCompletion)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
