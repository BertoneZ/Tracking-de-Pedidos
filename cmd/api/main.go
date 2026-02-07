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
	"tracking/internal/repository"
	"tracking/internal/service"
	"os"
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
		log.Println("Aviso: No se encontró el archivo .env, usando variables de entorno del sistema")
	}
	log.Printf("DEBUG: La URL de la base es: %s", os.Getenv("DB_URL"))
	pool, err := db.ConnectPostgres()
	if err != nil {
		log.Fatal("No se pudo conectar a la BD:", err)
	}
	rdb := db.ConnectRedis()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.RegisterUserRoutes(r, pool)
	routes.RegisterOrderRoutes(r, pool, rdb)
	productRepo := repository.NewProductRepository(pool)
	productService := service.NewProductService(productRepo)
	routes.RegisterProductRoutes(r, productService)

	r.Run(":8080")
}
