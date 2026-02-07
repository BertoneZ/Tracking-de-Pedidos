package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tracking/internal/domain"
	"github.com/jackc/pgx/v5/pgconn" 
	"fmt"
)
type UserRepositoryInterface interface {
	Create(ctx context.Context, u *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
    query := `INSERT INTO users (email, password_hash, full_name, role) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at`
    
    return r.db.QueryRow(ctx, query, 
        u.Email, 
        u.PasswordHash, 
        u.FullName, 
        u.Role,
    ).Scan(&u.ID, &u.CreatedAt)
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
   
    query := `SELECT id, email, password_hash, full_name, role FROM users WHERE email = $1`
    
    var u domain.User
    err := r.db.QueryRow(ctx, query, email).Scan(
        &u.ID, 
        &u.Email, 
        &u.PasswordHash, 
        &u.FullName, 
        &u.Role,
    )
    if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
        	return nil, fmt.Errorf("el usuario ya existe")
    	}

        return nil, err
    }
    return &u, nil
}