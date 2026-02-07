package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func ConnectPostgres() (*pgxpool.Pool, error) {

	connStr := os.Getenv("DB_URL")

	// Agregá este log para ver qué está leyendo realmente tu app
	fmt.Println("Intentando conectar a:", connStr)

	return pgxpool.New(context.Background(), connStr)
}

// Agregá este log para ver qué está leyendo realmente tu app
