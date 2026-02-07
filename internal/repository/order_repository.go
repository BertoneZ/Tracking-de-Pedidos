package repository

import (
	"context"
	"errors"
	"tracking/internal/domain"
	"tracking/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)
type OrderRepositoryInterface interface {
	GetPending(ctx context.Context) ([]domain.Order, error)
	GetOrderById(ctx context.Context, id string) (domain.Order, error)
	GetHistory(ctx context.Context, userID string) ([]domain.Order, error)	
	HasActiveOrder(ctx context.Context, driverID string) (bool, error)
	
	CreateWithItems(ctx context.Context, o *domain.Order) (string, error)
	CompleteOrder(ctx context.Context, orderID string, driverID string) error
	AcceptOrder(ctx context.Context, orderID string, driverID string) error
}
type OrderRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewOrderRepository(db *pgxpool.Pool, rdb *redis.Client) *OrderRepository {
	return &OrderRepository{db: db, rdb: rdb}
}

func (r *OrderRepository) GetPending(ctx context.Context) ([]domain.Order, error) {
	// 1. Traer la cabecera de las órdenes con el nombre del cliente
	query := `
    SELECT o.id, o.customer_id, u.full_name, o.status, 
           o.origin_lat, o.origin_lng, o.dest_lat, o.dest_lng, 
           o.destination_address, o.total_price, o.created_at
    FROM orders o
    JOIN users u ON o.customer_id = u.id -- El JOIN es clave
    WHERE o.status = 'PENDING'
    ORDER BY o.created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		err := rows.Scan(
			&o.ID, &o.CustomerID, &o.CustomerName, &o.Status,
			&o.OriginLat, &o.OriginLng, &o.DestLat, &o.DestLng,
			&o.DestinationAddress, &o.TotalPrice, &o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 2. Por cada orden, buscar sus productos (Ítems)
		itemQuery := `
            SELECT oi.product_id, p.name, oi.quantity, oi.price_at_time
            FROM order_items oi
            JOIN products p ON oi.product_id = p.id
            WHERE oi.order_id = $1`

		itemRows, err := r.db.Query(ctx, itemQuery, o.ID)
		if err != nil {
			return nil, err
		}

		// Inicializamos el slice para que no devuelva null en el JSON
		o.Items = []domain.OrderItem{}

		for itemRows.Next() {
			var item domain.OrderItem
			if err := itemRows.Scan(&item.ProductID, &item.ProductName, &item.Quantity, &item.PriceAtTime); err != nil {
				itemRows.Close()
				return nil, err
			}
			o.Items = append(o.Items, item)
		}
		itemRows.Close()

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
		return utils.ErrOrderNotAvailable
	}
	return nil
}
func (r *OrderRepository) GetOrderById(ctx context.Context, id string) (domain.Order, error) {
	// 1. Buscamos la cabecera de la orden
	queryOrder := `
		SELECT 
			o.id, o.customer_id, u_c.full_name,
			COALESCE(o.driver_id::TEXT, ''), COALESCE(u_d.full_name, ''),
			o.status, o.origin_lat, o.origin_lng, o.dest_lat, o.dest_lng, 
			o.destination_address, o.total_price, o.created_at
		FROM orders o
		JOIN users u_c ON o.customer_id = u_c.id
		LEFT JOIN users u_d ON o.driver_id = u_d.id
		WHERE o.id = $1`

	var o domain.Order
	err := r.db.QueryRow(ctx, queryOrder, id).Scan(
		&o.ID, &o.CustomerID, &o.CustomerName,
		&o.DriverID, &o.DriverName,
		&o.Status, &o.OriginLat, &o.OriginLng, &o.DestLat, &o.DestLng,
		&o.DestinationAddress, &o.TotalPrice, &o.CreatedAt,
	)
	if err != nil {
		return o, err
	}

	// 2. Buscamos los productos de esta orden (JOIN con products para el nombre)
	queryItems := `
		SELECT oi.product_id, p.name, oi.quantity, oi.price_at_time
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1`

	rows, err := r.db.Query(ctx, queryItems, id)
	if err != nil {
		return o, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		// Necesitás agregar el campo ProductName al struct OrderItem en domain si querés verlo
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity, &item.PriceAtTime); err != nil {
			return o, err
		}
		o.Items = append(o.Items, item)
	}

	return o, nil
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
    // 1. Agregamos el JOIN y las columnas que faltaban (nombre, precio, dirección)
    query := `
        SELECT 
            o.id, o.customer_id, u.full_name, o.status, 
			o.origin_lat, o.origin_lng, o.dest_lat, o.dest_lng,
            o.destination_address, o.total_price, o.created_at 
        FROM orders o
        JOIN users u ON o.customer_id = u.id
        WHERE (o.customer_id = $1 OR o.driver_id = $1) AND o.status = 'DELIVERED'
        ORDER BY o.created_at DESC`

    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []domain.Order
    for rows.Next() {
        var o domain.Order
        // 2. Escaneamos los nuevos campos
        err := rows.Scan(
            &o.ID, &o.CustomerID, &o.CustomerName, &o.Status, 
			&o.OriginLat, &o.OriginLng, &o.DestLat, &o.DestLng,
            &o.DestinationAddress, &o.TotalPrice, &o.CreatedAt,
        )
        if err != nil {
            return nil, err
        }

        // 3. Buscamos los ítems de cada orden (Igual que en GetPending)
        itemQuery := `
            SELECT oi.product_id, p.name, oi.quantity, oi.price_at_time
            FROM order_items oi
            JOIN products p ON oi.product_id = p.id
            WHERE oi.order_id = $1`
        
        itemRows, err := r.db.Query(ctx, itemQuery, o.ID)
        if err != nil {
            return nil, err
        }

        o.Items = []domain.OrderItem{} 
        for itemRows.Next() {
            var item domain.OrderItem
            if err := itemRows.Scan(&item.ProductID, &item.ProductName, &item.Quantity, &item.PriceAtTime); err != nil {
                itemRows.Close()
                return nil, err
            }
            o.Items = append(o.Items, item)
        }
        itemRows.Close()

        orders = append(orders, o)
    }
    return orders, nil
}
func (r *OrderRepository) HasActiveOrder(ctx context.Context, driverID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE driver_id = $1 AND status = 'ASSIGNED')`

	var exists bool
	err := r.db.QueryRow(ctx, query, driverID).Scan(&exists)
	return exists, err
}
func (r *OrderRepository) CreateWithItems(ctx context.Context, o *domain.Order) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	// A. Insertar Orden incluyendo las nuevas columnas numéricas
	queryOrder := `
        INSERT INTO orders (
            customer_id, status, destination_address, total_price, 
            origin_lat, origin_lng, dest_lat, dest_lng,
            origin, destination
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8,
            ST_SetSRID(ST_MakePoint($6, $5), 4326)::geography, 
            ST_SetSRID(ST_MakePoint($8, $7), 4326)::geography
        )
        RETURNING id`

	var orderID string
	// El orden de los parámetros es fundamental:
	err = tx.QueryRow(ctx, queryOrder,
		o.CustomerID,         // $1
		o.Status,             // $2
		o.DestinationAddress, // $3
		o.TotalPrice,         // $4
		o.OriginLat,          // $5
		o.OriginLng,          // $6
		o.DestLat,            // $7
		o.DestLng,            // $8
	).Scan(&orderID)

	if err != nil {
		return "", err
	}

	// B. Insertar Items (Sigue igual)
	for _, item := range o.Items {
		queryItem := `INSERT INTO order_items (order_id, product_id, quantity, price_at_time) VALUES ($1, $2, $3, $4)`
		if _, err := tx.Exec(ctx, queryItem, orderID, item.ProductID, item.Quantity, item.PriceAtTime); err != nil {
			return "", err
		}
	}

	return orderID, tx.Commit(ctx)
}
