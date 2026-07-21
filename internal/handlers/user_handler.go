package handlers

import (
	"errors"
	"net/http"
	"rest-api/internal/config"
	"rest-api/internal/models"
	"rest-api/internal/repository"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


type UserRegister struct{
	Username string `json:"username" binding:"required"`
	Email string    `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct{
	Token string `json:"token"`
}


func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var RegisterRequest UserRegister

		if err := c.BindJSON(&RegisterRequest); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error" : err.Error()})
			return
		}

		if strings.TrimSpace(RegisterRequest.Username) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error" : "Username must be filled"})
			return
		}

		if strings.TrimSpace(RegisterRequest.Email) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error" : "Email must be filled"})
			return
		}

		if len(RegisterRequest.Password ) < 6 {
			c.JSON(http.StatusBadRequest,gin.H{"error" : "Password must be more than 6 words"})
			return
		}

		

		HashPass,err := bcrypt.GenerateFromPassword([]byte(RegisterRequest.Password),bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error" : " failed to hash password" +  err.Error()})
			return
		}
		user := &models.User{
			Username: RegisterRequest.Username,
			Email: RegisterRequest.Email,
			Password: string(HashPass),
		}

		CreateUser,err := repository.CreateUser(pool,user)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 = Unique Violation in Postgres
				if strings.Contains(pgErr.ConstraintName, "username") || strings.Contains(pgErr.Detail, "username") {
					c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
					return
				}
				if strings.Contains(pgErr.ConstraintName, "email") || strings.Contains(pgErr.Detail, "email") {
					c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
					return
				}
				
				c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,CreateUser)

	}
}

func LoginHandler(pool *pgxpool.Pool,cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var LoginRequest UserLogin

		if err := c.BindJSON(&LoginRequest); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error" : err.Error()})
			return
		}

		user,err := repository.GetUserByUsername(pool,LoginRequest.Username)
		if err != nil{
			 c.JSON(http.StatusUnauthorized, gin.H{"error" : "invalid credential"})
			 return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(LoginRequest.Password))
		if err != nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error" : "Invalid credential"})
			return 
		}
		claims := jwt.MapClaims{
			"user_id" : user.ID,
			"username" : user.Username,
			"role"		: user.Role,
			"exp" : time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
		tokenstring,err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error" : "failed generate token jwt" + err.Error()})
			return 
		}

		c.JSON(http.StatusOK,LoginResponse{Token : tokenstring})
	}
}