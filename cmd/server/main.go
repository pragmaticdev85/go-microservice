package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/pragmaticdev85/go-microservice/docs"
	"github.com/pragmaticdev85/go-microservice/internal/config"
	"github.com/pragmaticdev85/go-microservice/internal/controllers"
	"github.com/pragmaticdev85/go-microservice/internal/repositories"
	"github.com/pragmaticdev85/go-microservice/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go Microservice API
// @version 1.0
// @description This is a sample microservice with Go, Gin, MongoDB
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize MongoDB client
	dbClient, err := repositories.NewMongoDBClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer dbClient.Disconnect()

	// Initialize repositories
	exampleRepo := repositories.NewExampleRepository(dbClient, cfg.MongoDBName)

	// Initialize services
	exampleService := services.NewExampleService(exampleRepo)

	// Initialize controllers
	exampleController := controllers.NewExampleController(exampleService)

	// Setup Gin router
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		example := api.Group("/examples")
		{
			example.POST("", exampleController.CreateExample)
			// example.GET("", exampleController.GetExamples)
			example.GET("/:id", exampleController.GetExampleByID)
			// example.PUT("/:id", exampleController.UpdateExample)
			// example.DELETE("/:id", exampleController.DeleteExample)
		}
	}

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Server running on port %s", cfg.Port)
	router.Run(":" + cfg.Port)
}
