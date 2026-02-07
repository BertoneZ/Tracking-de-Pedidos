package service

import (
	"context"
	"tracking/internal/auth"
	"tracking/internal/domain"
	"tracking/internal/repository"
	"golang.org/x/crypto/bcrypt"
)
type UserServiceInterface interface {
	Register(ctx context.Context, email, password, fullName, role string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
}
type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password, fullName, role string) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
		FullName:   fullName,
		Role:         role,
	}
	err = s.repo.Create(ctx, user)
	return user, err
}
func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", err 
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	return user, token, err
}