package handlers

import (
	"errors"
	"log"
	"net/http"
	"rest-api/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateCategory struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UpdateCategory struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}

func CreateCategoryHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateCategory
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category, err := repository.CreateCategory(pool, input.Name, input.Slug)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "Category name or slug already exists"})
				return
			}
			log.Printf("CreateCategory Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category": category})
	}
}

func GetAllCategoryHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := repository.GetAllCategory(pool)

		if err != nil {
			log.Printf("GetAllCategory Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Categories retrieved successfully", "categories": categories})
	}
}

func GetCategoryByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")

		id, err := strconv.Atoi(idstr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		category, err := repository.GetCategoryByID(pool, id)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
				return
			}

			log.Printf("GetCategoryByID Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category retrieved successfully", "category": category})
	}
}

func UpdateCategoryHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		var input UpdateCategory

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name == nil && input.Slug == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "At least one field (name or slug) must be filled",
			})
			return
		}

		exists, err := repository.GetCategoryByID(pool, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
				return
			}
			log.Printf("UpdateCategory GetByID Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		name := exists.Name
		if input.Name != nil {
			name = *input.Name
		}

		slug := exists.Slug
		if input.Slug != nil {
			slug = *input.Slug
		}

		category, err := repository.UpdateCategory(pool, id, name, slug)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "Category name or slug already exists"})
				return
			}
			log.Printf("UpdateCategory Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "category": category})
	}
}

func DeleteCategoryHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstr := c.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		category, err := repository.DeleteCategory(pool, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
				return
			}

			log.Printf("DeleteCategory Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Category deleted successfully",
			"category": category,
		})
	}
}