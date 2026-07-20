package handlers

import (
	"net/http"
	"rest-api/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateBlogRequest struct {
	Title   *string `json:"title"`
    Content *string `json:"content"`
}

func CreateBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){

		userIDInterface,exists := c.Get("user_id")

		if !exists{
			c.JSON(http.StatusInternalServerError,gin.H{"error" : "user_id not found bruh :D"})
			return 
		}
		userID := userIDInterface.(string)

		var input CreateBlogRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		blog, err := repository.CreateBlog(pool, input.Title, input.Content,userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully", "blog": blog})
	}
}

func GetALLBlogsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		blogs, err := repository.GetAllBlogs(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) 
			return 
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blogs retrieved successfully", "blogs": blogs})
	}
}

func GetBlogsByUserIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface,exists := c.Get("user_id")

		if !exists{
			c.JSON(http.StatusInternalServerError,gin.H{"error" : "user_id not found bruh :D"})
			return 
		}
		userID := userIDInterface.(string)
		blogs, err := repository.GetBlogsByUserID(pool,userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) 
			return 
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blogs retrieved successfully", "blogs": blogs})
	}
}

func GetBlogByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		idStr :=  c.Param("id")

		id,err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid blog id"})
			return 
		}

		blog,err := repository.GetBlogByID(pool,id)
		if err != nil {
			if err == pgx.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error" : "Blog Not Found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Blog retrieved successfully", "blog": blog})

	}
}

func UpdateBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {

		userIDInterface,exists := c.Get("user_id")

		if !exists{
			c.JSON(http.StatusInternalServerError,gin.H{"error" : "user_id not found bruh :D"})
			return 
		}
		userID := userIDInterface.(string)
		idSTr := c.Param("id")

		id,err := strconv.Atoi(idSTr)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid Blog ID"})
			return 
		}

		var input UpdateBlogRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
			return
		}

		if input.Title == nil && input.Content == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "at leats one field ( Title / content) must be filled",
			})
			return
		}

		exist, err := repository.GetBlogByID(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if exist.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to edit this blog"}) 
			return
		}

		title := exist.Title
		if input.Title != nil{
			title = *input.Title
		}

		
		content := exist.Content
		if input.Content != nil {
			content = *input.Content
		}

		blog,err := repository.UpdateBlog(pool,id,title,content,userID)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error" : err.Error()})
			return
		}
		

		c.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully", "blog": blog})

		
	}
}

func DeleteBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		userIDInterface,exists := c.Get("user_id")

		if !exists{
			c.JSON(http.StatusInternalServerError,gin.H{"error" : "user_id not found bruh :D"})
			return 
		}
		userID := userIDInterface.(string)
		idSTr := c.Param("id")

		id,err := strconv.Atoi(idSTr)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid Blog ID"})
			return 
		}

		blog,err := repository.DeleteBlog(pool,id,userID)
		if err != nil {
			if err == pgx.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error" : "Blog Not Found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Blog deleted successfully",
			"blog":    blog,
		})
	}
}