package service

import (
	"context"
	"errors"
	"strings"
	"tracking/internal/dto"
	"tracking/internal/repository"
	"tracking/internal/utils"

	"github.com/jackc/pgx/v5"
)

type ProductServiceInterface interface {
	GetProductsService(ctx context.Context) ([]dto.ProductResponse, error)
	CreateProductService(ctx context.Context, req dto.UpsertProductRequest) (dto.ProductResponse, error)
	UpdateProductService(ctx context.Context, id string, req dto.UpdateProductRequest) (dto.ProductResponse, error)
	DeleteProductService(ctx context.Context, id string) error
}
type ProductService struct {
	repo repository.ProductRepositoryInterface
}

func NewProductService(repo repository.ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProductsService(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.repo.GetProductsRepo(ctx)
	if err != nil {
		return nil, err
	}

	return utils.SliceProductDomainToProductResponseListDto(products), nil
}

func (s *ProductService) CreateProductService(ctx context.Context, req dto.UpsertProductRequest) (dto.ProductResponse, error) {
	productDomain := utils.ToProductDomain(req)
	created, err := s.repo.Create(ctx, productDomain)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return utils.ToProductResponse(created), nil
}

func (s *ProductService) UpdateProductService(ctx context.Context, id string, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ProductResponse{}, pgx.ErrNoRows
		}
		return dto.ProductResponse{}, err
	}

	updatedInput := current

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name != "" {
			updatedInput.Name = name
		}
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return dto.ProductResponse{}, errors.New("precio debe ser mayor a 0")
		}
		updatedInput.Price = *req.Price
	}

	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)
		if description != "" {
			updatedInput.Description = description
		}
	}

	updated, err := s.repo.Update(ctx, id, updatedInput)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return utils.ToProductResponse(updated), nil
}

func (s *ProductService) DeleteProductService(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
