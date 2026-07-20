package main

import (
	"rest-api/internal/config"
	"rest-api/internal/database"
	"rest-api/internal/handlers"
	"rest-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main(){
	var cfg *config.Config
	var err error
	cfg, err = config.Load()
	if err != nil {
		panic(err)
	}

	var pool  *pgxpool.Pool
	pool,err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	var route *gin.Engine = gin.Default()
	route.SetTrustedProxies(nil)
	route.GET("/",func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": " Go Gin API is running",
			"status": "success",
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
		protectedBlogs.PUT("/:id", handlers.UpdateBlogHandler(pool))
		protectedBlogs.DELETE("/:id", handlers.DeleteBlogHandler(pool))
	}

	protectedCategories := route.Group("/categories")
	protectedCategories.Use(middleware.AuthMiddleware(cfg))
	{
		protectedCategories.POST("", handlers.CreateCategoryHandler(pool))
	}

	route.Run(":" + cfg.Port)
}