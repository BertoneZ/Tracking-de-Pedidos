	package db
	import (
		"context"

		"github.com/jackc/pgx/v5/pgxpool"
	)

	func ConnectPostgres() (*pgxpool.Pool, error) {
		// Reemplazamos os.Getenv por el string directo (Hardcoded)
		dsn := "postgres://user_logistics:pass_logistics@localhost:5433/logistics_db?sslmode=disable"
		
		pool, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			return nil, err
		}

		// El Ping confirma que la contraseña 'pass_logistics' fue aceptada
		if err := pool.Ping(context.Background()); err != nil {
			return nil, err
		}
		
		return pool, nil
	}