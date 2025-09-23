package main

import (
	"log"
	"fmt"

	"github.com/gin-gonic/gin"
	"todo-backend/internal/config"
	"todo-backend/internal/handlers"
	"todo-backend/internal/middleware"
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
	"todo-backend/internal/services"
	"todo-backend/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		TimeZone: cfg.Database.TimeZone,
	}

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Run auto-migration for development
	if cfg.IsDevelopment() {
		log.Println("Running auto-migration...")
		if err := db.AutoMigrate(models.AllModels()...); err != nil {
			log.Fatalf("Failed to run auto-migration: %v", err)
		}
		log.Println("Auto-migration completed")
	}

	// Initialize repositories
	todoRepo := repository.NewTodoRepository(db.GetDB())
	categoryRepo := repository.NewCategoryRepository(db.GetDB())

	// Initialize services
	todoService := services.NewTodoService(todoRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Initialize Gin router
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.StructuredLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.Security())
	router.Use(middleware.RateLimitHeaders())

	// Setup routes
	handlers.SetupRoutes(router, todoService, categoryService)

	// Handle 404
	router.NoRoute(middleware.NotFoundHandler())

	// Handle 405
	router.NoMethod(middleware.MethodNotAllowedHandler())

	// Start server
	port := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on port %s", cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Env)
	log.Printf("Health check available at: http://localhost%s/api/health", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}