package handlers

import (
	"net/http"
	"rest-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateCategory struct{
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

func CreateCategoryHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		var input CreateCategory
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category,err := repository.CreateCategory(pool,input.Name,input.Slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "category created successfully", "category": category})
	}
}