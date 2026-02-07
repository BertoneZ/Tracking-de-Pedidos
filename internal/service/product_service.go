package service

import (
	"context"
	"tracking/internal/utils"
	"tracking/internal/dto"
	"tracking/internal/repository"
)
type ProductServiceInterface interface {
	GetProductsService(ctx context.Context) ([]dto.ProductResponse, error)
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
	// Usamos el mapper que creamos antes
	return utils.SliceProductDomainToProductResponseListDto(products), nil
}