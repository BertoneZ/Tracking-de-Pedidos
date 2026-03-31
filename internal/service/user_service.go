package service

import (
	"context"
	"errors"
	"os"
	"strings"
	"tracking/internal/auth"
	"tracking/internal/domain"
	"tracking/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register(ctx context.Context, email, password, fullName, role string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	BootstrapAdmin(ctx context.Context, email, password, fullName, secret string) (*domain.User, error)
	ListUsers(ctx context.Context, role string, active *bool) ([]domain.User, error)
	DeactivateUser(ctx context.Context, actorUserID, targetUserID string) error
}
type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

const DefaultRole = "customer"

func (s *UserService) Register(ctx context.Context, email, password, fullName, role string) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	finalRole := strings.ToLower(strings.TrimSpace(role))
	if finalRole == "" || finalRole == "admin" {
		finalRole = DefaultRole
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
		FullName:     fullName,
		Role:         finalRole,
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

func (s *UserService) BootstrapAdmin(ctx context.Context, email, password, fullName, secret string) (*domain.User, error) {
	expectedSecret := os.Getenv("ADMIN_BOOTSTRAP_SECRET")
	if expectedSecret == "" || secret != expectedSecret {
		return nil, errors.New("secret de bootstrap inválido")
	}

	exists, err := s.repo.AdminExists(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ya existe un admin en el sistema")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
		FullName:     fullName,
		Role:         "admin",
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, role string, active *bool) ([]domain.User, error) {
	return s.repo.ListUsers(ctx, role, active)
}

func (s *UserService) DeactivateUser(ctx context.Context, actorUserID, targetUserID string) error {
	if actorUserID == targetUserID {
		return errors.New("no puedes desactivarte a vos mismo")
	}

	return s.repo.DeactivateUser(ctx, targetUserID)
}
