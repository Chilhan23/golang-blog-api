package handlers

import (
	"net/http"
	"rest-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreateBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context){
		var input CreateBlogRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		blog, err := repository.CreateBlog(pool, input.Title, input.Content)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"blog": blog})
	}
}

func GetALLBlogsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		blogs, err := repository.GetAllBlogs(pool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) 
			return 
		}

		c.JSON(http.StatusOK, blogs)
	}
}