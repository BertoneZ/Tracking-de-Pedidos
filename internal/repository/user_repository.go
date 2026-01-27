package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tracking/internal/domain"
)

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
	var u domain.User
	query := `SELECT id, email, password_hash, role FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}