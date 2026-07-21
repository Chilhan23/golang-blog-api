package middleware

import (
	"fmt"
	"net/http"
	"rest-api/internal/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware (cfg *config.Config) gin.HandlerFunc{
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error" : "Authorization header required"})
			c.Abort() 
			return 
		}

		tokenstring := strings.TrimPrefix(authHeader,"Bearer ")

		if tokenstring == "" || tokenstring == authHeader {
			c.JSON(http.StatusUnauthorized,gin.H{"error" : "invalid authorization"})
			c.Abort()
			return 
		}

		token,err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{},error){
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg(){
				return nil,fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret),nil
		})

		if err != nil || !token.Valid{
			c.JSON(http.StatusUnauthorized,gin.H{"error" : "invalid or expired token"})
			c.Abort()
			return 
		}

		claims,ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized,gin.H{"error" : "invalid  token claim"})
			c.Abort()
			return 
		}

		userID,ok := claims["user_id"].(string)

		if !ok {
			c.JSON(http.StatusUnauthorized,gin.H{"error" : "invalid  token payload"})
			c.Abort()
			return 
		}

		if role, ok := claims["role"].(string); ok {
			c.Set("role", role) 
		}

		if exp,ok := claims["exp"].(float64); ok {
			expirationtime := time.Unix(int64(exp),0)

			if time.Now().After(expirationtime){
				c.JSON(http.StatusUnauthorized,gin.H{"error" : "invalid  token expired"})
				c.Abort()
				return 
			}
		}

		c.Set("user_id",userID)
		c.Next()
	}
}