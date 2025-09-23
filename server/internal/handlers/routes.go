package handlers

import (
	"github.com/gin-gonic/gin"
	"todo-backend/internal/services"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, todoService services.TodoService, categoryService services.CategoryService) {
	// Create handlers
	todoHandler := NewTodoHandler(todoService)
	categoryHandler := NewCategoryHandler(categoryService)

	// API version group
	api := r.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "OK",
				"message": "Todo API is running",
			})
		})

		// Todo routes
		todos := api.Group("/todos")
		{
			todos.POST("", todoHandler.CreateTodo)           // POST /api/todos
			todos.GET("", todoHandler.ListTodos)             // GET /api/todos
			todos.GET("/:id", todoHandler.GetTodo)           // GET /api/todos/:id
			todos.PUT("/:id", todoHandler.UpdateTodo)        // PUT /api/todos/:id
			todos.DELETE("/:id", todoHandler.DeleteTodo)     // DELETE /api/todos/:id
			todos.PATCH("/:id/complete", todoHandler.ToggleTodoComplete) // PATCH /api/todos/:id/complete
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)      // POST /api/categories
			categories.GET("", categoryHandler.ListCategories)       // GET /api/categories
			categories.GET("/all", categoryHandler.GetAllCategories) // GET /api/categories/all
			categories.GET("/:id", categoryHandler.GetCategory)      // GET /api/categories/:id
			categories.PUT("/:id", categoryHandler.UpdateCategory)   // PUT /api/categories/:id
			categories.DELETE("/:id", categoryHandler.DeleteCategory) // DELETE /api/categories/:id
		}
	}
}