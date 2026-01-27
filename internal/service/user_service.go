package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"tracking/internal/domain"
	"tracking/internal/repository"
	"tracking/internal/auth"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, role string) (*domain.User, error) {
	// Hashear contraseña
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
		Role:         role,
	}

	err := s.repo.Create(ctx, user)
	return user, err
}
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err // Usuario no encontrado
	}

	// Comparamos el hash de la BD con la password que envía el usuario
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err // Contraseña incorrecta
	}

	// Si todo está bien, generamos el JWT
	return auth.GenerateToken(user.ID, user.Role)
}