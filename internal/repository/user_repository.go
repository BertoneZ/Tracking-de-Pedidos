package repository

import (
	"context"
	"fmt"
	"tracking/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, u *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	AdminExists(ctx context.Context) (bool, error)
	ListUsers(ctx context.Context, role string, active *bool) ([]domain.User, error)
	DeactivateUser(ctx context.Context, userID string) error
}
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (email, password_hash, full_name, role, is_active) 
              VALUES ($1, $2, $3, $4, true) RETURNING id, created_at, is_active`

	return r.db.QueryRow(ctx, query,
		u.Email,
		u.PasswordHash,
		u.FullName,
		u.Role,
	).Scan(&u.ID, &u.CreatedAt, &u.IsActive)
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	query := `SELECT id, email, password_hash, full_name, role, is_active FROM users WHERE email = $1 AND is_active = true`

	var u domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.Role,
		&u.IsActive,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, fmt.Errorf("el usuario ya existe")
		}

		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, full_name, role, is_active, created_at FROM users WHERE id = $1`

	var u domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) AdminExists(ctx context.Context) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin' AND is_active = true)`
	var exists bool
	err := r.db.QueryRow(ctx, query).Scan(&exists)
	return exists, err
}

func (r *UserRepository) ListUsers(ctx context.Context, role string, active *bool) ([]domain.User, error) {
	query := "SELECT id, email, COALESCE(full_name, ''), role, COALESCE(is_active, true), COALESCE(created_at, NOW()) FROM users"
	var args []interface{}

	if role != "" && active != nil {
		query += " WHERE role = $1 AND is_active = $2"
		args = []interface{}{role, *active}
	} else if role != "" {
		query += " WHERE role = $1"
		args = []interface{}{role}
	} else if active != nil {
		query += " WHERE is_active = $1"
		args = []interface{}{*active}
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *UserRepository) DeactivateUser(ctx context.Context, userID string) error {
	query := `UPDATE users SET is_active = false WHERE id = $1 AND is_active = true`
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("usuario no encontrado o ya inactivo")
	}

	return nil
}
