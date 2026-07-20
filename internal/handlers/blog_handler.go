package handlers

import (
	"errors"
	"log"
	"net/http"
	"rest-api/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateBlogRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID *int   `json:"category_id"`
}

type UpdateBlogRequest struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	CategoryID *int    `json:"category_id"`
}

func CreateBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDInterface.(string)

		var input CreateBlogRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		blog, err := repository.CreateBlog(pool, input.Title, input.Content, userID, input.CategoryID)

		if err != nil {
			log.Printf("CreateBlog Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully", "blog": blog})
	}
}

func GetALLBlogsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		blogs, err := repository.GetAllBlogs(pool)

		if err != nil {
			log.Printf("GetAllBlogs Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blogs retrieved successfully", "blogs": blogs})
	}
}

func GetBlogsByUserIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDInterface.(string)
		blogs, err := repository.GetBlogsByUserID(pool, userID)

		if err != nil {
			log.Printf("GetBlogsByUserID Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user blogs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blogs retrieved successfully", "blogs": blogs})
	}
}

func GetBlogByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
			return
		}

		blog, err := repository.GetBlogByID(pool, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
				return
			}

			log.Printf("GetBlogByID Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blog retrieved successfully", "blog": blog})

	}
}

func UpdateBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDInterface.(string)
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
			return
		}

		var input UpdateBlogRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Content == nil && input.CategoryID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "At least one field (title, content, or category_id) must be filled",
			})
			return
		}

		exist, err := repository.GetBlogByID(pool, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
				return
			}
			log.Printf("UpdateBlog GetByID Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if exist.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to edit this blog"})
			return
		}

		title := exist.Title
		if input.Title != nil {
			title = *input.Title
		}

		content := exist.Content
		if input.Content != nil {
			content = *input.Content
		}

		categoryID := exist.CategoryID
		if input.CategoryID != nil {
			categoryID = input.CategoryID
		}

		blog, err := repository.UpdateBlog(pool, id, title, content, userID, categoryID)

		if err != nil {
			log.Printf("UpdateBlog Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully", "blog": blog})

	}
}

func DeleteBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDInterface.(string)
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
			return
		}

		blog, err := repository.DeleteBlog(pool, id, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
				return
			}

			log.Printf("DeleteBlog Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Blog deleted successfully",
			"blog":    blog,
		})
	}
}