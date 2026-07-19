package main

import (
	"rest-api/internal/config"
	"rest-api/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main(){
	var cfg *config.Config
	var err error
	cfg, err = config.Load()
	if err != nil {
		panic(err)
	}

	var pool  *pgxpool.Pool
	pool,err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	var route *gin.Engine = gin.Default()
	route.SetTrustedProxies(nil)
	route.GET("/",func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": " Go Gin API is running",
			"status": "success",
			"database": "connected",
		})
	})

	route.Run(":" + cfg.Port)
}	