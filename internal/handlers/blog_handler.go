package handlers

import (
	"net/http"
	"rest-api/internal/repository"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateBlogRequest struct {
	Title   string `json:"title"`
    Content string `json:"content"`
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

func GetBlogByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		idStr :=  c.Param("id")

		id,err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid todo id"})
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

		c.JSON(http.StatusOK,blog)

	}
}

func UpdateBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
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

		if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Title and content must be filled",
			})
			return
		}

		blog,err := repository.UpdateBlog(pool,id,input.Title,input.Content)
		if err != nil {
			if err == pgx.ErrNoRows{
				c.JSON(http.StatusNotFound, gin.H{"error" : "Blog not found!"})
				return
			}

			c.JSON(http.StatusInternalServerError,gin.H{"error" : err.Error()})
		}

		c.JSON(http.StatusOK,blog)

		
	}
}