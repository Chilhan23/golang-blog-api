package handlers

import (
	"net/http"
	"strconv"

	"rest-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ToggleLikeHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		blogIDStr := c.Param("id")
		blogID, err := strconv.Atoi(blogIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
			return
		}

		isLiked, totalLikes, err := repository.TogleLike(pool, userID.(string), blogID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like status"})
			return
		}

		message := "Blog liked successfully"
		if !isLiked {
			message = "Blog unliked successfully"
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     message,
			"is_liked":    isLiked,
			"total_likes": totalLikes,
		})
	}
}