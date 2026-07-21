package handlers

import (
	"errors"
	"net/http"
	"rest-api/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateComment struct {
	Content string `json:"content" binding:"required"`
}

func CreateCommentHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
			return
		}

		blogIDStr := c.Param("id")
		blogID, err := strconv.Atoi(blogIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
			return
		}

		var input CreateComment
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		comment, err := repository.CreateComment(pool, blogID, userIDStr, input.Content)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create Comment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "comment created succesfully", "comment": comment})
	}
}

func GetCommentByIDBlogHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		blogIDStr := c.Param("id")
		blogID, err := strconv.Atoi(blogIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
			return
		}

		comments, err := repository.GetCommentsByBlogID(pool, blogID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retreived comments at this blog"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Comments retreived successfully", "comments": comments})
	}
}

func DeleteCommentHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userIDStr := userID.(string)

		roleVal, _ := c.Get("role")
		roleStr, _ := roleVal.(string)

		commentIDStr := c.Param("id")
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		deletedComment, err := repository.DeleteComment(pool, commentID, userIDStr, roleStr)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Comment not found or unauthorized to delete"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Comment deleted successfully",
			"comment": deletedComment,
		})
	}
}
