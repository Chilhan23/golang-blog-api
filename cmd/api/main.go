package main

import (
	"rest-api/internal/config"
	"rest-api/internal/database"
	"rest-api/internal/handlers"
	"rest-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()
	if err != nil {
		panic(err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	var route *gin.Engine = gin.Default()
	route.Use(CORSMiddleware())
	route.SetTrustedProxies(nil)

	route.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  " Go Gin API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	// User Routes
	route.POST("/auth/register", handlers.CreateUserHandler(pool))
	route.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	// Public Blog Routes
	route.GET("/blogs", handlers.GetALLBlogsHandler(pool))
	route.GET("/blogs/:id", handlers.GetBlogByIDHandler(pool))

	// Protected Blog Routes
	protectedBlogs := route.Group("/blogs")
	protectedBlogs.Use(middleware.AuthMiddleware(cfg))
	{
		protectedBlogs.POST("", handlers.CreateBlogHandler(pool))
		protectedBlogs.GET("/user", handlers.GetBlogsByUserIDHandler(pool))
		protectedBlogs.POST("/:id/like", handlers.ToggleLikeHandler(pool))
		protectedBlogs.PUT("/:id", handlers.UpdateBlogHandler(pool))
		protectedBlogs.DELETE("/:id", handlers.DeleteBlogHandler(pool))
	}

	// Public Category Routes
	route.GET("/categories", handlers.GetAllCategoryHandler(pool))
	route.GET("/categories/:id", handlers.GetCategoryByIDHandler(pool))

	// Admin Category Routes
	adminCategories := route.Group("/categories")
	adminCategories.Use(middleware.AuthMiddleware(cfg), middleware.AdminMiddleware())
	{
		adminCategories.POST("", handlers.CreateCategoryHandler(pool))
		adminCategories.PUT("/:id", handlers.UpdateCategoryHandler(pool))
		adminCategories.DELETE("/:id", handlers.DeleteCategoryHandler(pool))
	}

	route.Run(":" + cfg.Port)
}