package repository
import (
	"context"
	"tracking/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"errors"
	"github.com/redis/go-redis/v9"
)
type OrderRepository struct {
	db *pgxpool.Pool
	rdb *redis.Client
}

func NewOrderRepository(db *pgxpool.Pool, rdb *redis.Client) *OrderRepository {
	return &OrderRepository{db: db, rdb: rdb}
}

func (r *OrderRepository) Create(ctx context.Context, o *domain.Order) error {
	query := `
		INSERT INTO orders (customer_id, status, origin, destination)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326)::geography, ST_SetSRID(ST_MakePoint($5, $6), 4326)::geography)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		o.CustomerID,
		o.Status,
		o.OriginLng, o.OriginLat, // Importante: PostGIS usa Longitud, Latitud
		o.DestLng, o.DestLat,
	).Scan(&o.ID, &o.CreatedAt)
}

func (r *OrderRepository) GetPending(ctx context.Context) ([]domain.Order, error) {
	query := `SELECT id, customer_id, status, 
	          ST_Y(origin::geometry) as lat, ST_X(origin::geometry) as lng,
	          ST_Y(destination::geometry) as dlat, ST_X(destination::geometry) as dlng,
	          created_at 
	          FROM orders WHERE status = 'PENDING' 
	          ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		err := rows.Scan(&o.ID, &o.CustomerID, &o.Status, &o.OriginLat, &o.OriginLng, &o.DestLat, &o.DestLng, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
func (r *OrderRepository) AcceptOrder(ctx context.Context, orderID string, driverID string) error {
	query := `UPDATE orders 
	          SET driver_id = $1, status = 'ASSIGNED' 
	          WHERE id = $2 AND status = 'PENDING'`
	
	result, err := r.db.Exec(ctx, query, driverID, orderID)
	if err != nil {
		return err
	}

	// Si ninguna fila fue afectada, es porque el pedido ya no estaba PENDING
	if result.RowsAffected() == 0 {
		return errors.New("el pedido ya fue tomado por otro repartidor o no existe")
	}
	return nil
}
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT id, customer_id, driver_id, status FROM orders WHERE id = $1`
	
	var o domain.Order
	var driverID *string // Puntero para manejar el NULL de la base de datos

	// Escaneamos los datos
	err := r.db.QueryRow(ctx, query, id).Scan(&o.ID, &o.CustomerID, &driverID, &o.Status)
	if err != nil {
		return nil, err
	}

	// Si el valor no es NULL en la BD, se lo asignamos al struct
	if driverID != nil {
		o.DriverID = *driverID 
	} else {
		o.DriverID = "" // Opcional: aseguramos que esté vacío si es NULL
	}

	return &o, nil
}
func (r *OrderRepository) CompleteOrder(ctx context.Context, orderID string, driverID string) error {
	// 1. Actualizamos Postgres
	query := `UPDATE orders SET status = 'DELIVERED' 
	          WHERE id = $1 AND driver_id = $2 AND status = 'ASSIGNED'`
	
	res, err := r.db.Exec(ctx, query, orderID, driverID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("no se pudo completar el pedido (revisar ID o estado)")
	}

	// 2. Limpiamos Redis para este driver (opcional, pero recomendado)
	// Se usa ZREM porque GEOADD crea un Sorted Set internamente
	r.rdb.ZRem(ctx, "drivers_locations", driverID)

	return nil
}
func (r *OrderRepository) GetHistory(ctx context.Context, userID string) ([]domain.Order, error) {
	query := `SELECT id, customer_id, driver_id, status, created_at 
	          FROM orders 
	          WHERE (customer_id = $1 OR driver_id = $1) AND status = 'DELIVERED'
	          ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		// Usamos puntero para driver_id por si hay nulos, aunque en DELIVERED no debería
		var driverID *string 
		err := rows.Scan(&o.ID, &o.CustomerID, &driverID, &o.Status, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		if driverID != nil { o.DriverID = *driverID }
		orders = append(orders, o)
	}
	return orders, nil
}