package main

import (
	_"tracking/docs"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"tracking/internal/db"
	"tracking/internal/routes"
)
// @title API de Logística Rafaela
// @version 1.0
// @description Servidor de tracking de pedidos en tiempo real.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool, err := db.ConnectPostgres()
	if err != nil {
		log.Fatal("No se pudo conectar a la BD:", err)
	}
	rdb := db.ConnectRedis()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.RegisterUserRoutes(r, pool)
	routes.RegisterOrderRoutes(r, pool, rdb)

	r.Run(":8080")
}
