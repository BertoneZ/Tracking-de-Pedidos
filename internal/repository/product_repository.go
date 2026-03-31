package repository

import (
	"context"
	"tracking/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepositoryInterface interface {
	GetProductsRepo(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id string) (domain.Product, error)
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	Update(ctx context.Context, id string, p domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id string) error
}
type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductsRepo(ctx context.Context) ([]domain.Product, error) {
	query := `SELECT id, name, price, description, is_active FROM products WHERE is_active = true`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.IsActive); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (domain.Product, error) {
	query := `SELECT id, name, price, description, is_active FROM products WHERE id = $1 AND is_active = true`
	var p domain.Product
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.IsActive)
	return p, err
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	query := `INSERT INTO products (name, price, description, is_active) VALUES ($1, $2, $3, true) RETURNING id, name, price, description, is_active`
	var created domain.Product
	err := r.db.QueryRow(ctx, query, p.Name, p.Price, p.Description).Scan(
		&created.ID,
		&created.Name,
		&created.Price,
		&created.Description,
		&created.IsActive,
	)
	return created, err
}

func (r *ProductRepository) Update(ctx context.Context, id string, p domain.Product) (domain.Product, error) {
	query := `UPDATE products SET name = $1, price = $2, description = $3 WHERE id = $4 AND is_active = true RETURNING id, name, price, description, is_active`
	var updated domain.Product
	err := r.db.QueryRow(ctx, query, p.Name, p.Price, p.Description, id).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Price,
		&updated.Description,
		&updated.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Product{}, pgx.ErrNoRows
		}
		return domain.Product{}, err
	}

	return updated, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE products SET is_active = false WHERE id = $1 AND is_active = true`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
