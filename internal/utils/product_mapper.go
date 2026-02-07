package utils

import (
	"tracking/internal/domain"
	"tracking/internal/dto"
)

// 1. DTO Request -> Domain (Para crear o actualizar)
func ToProductDomain(req dto.UpsertProductRequest) domain.Product {
	return domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}
}

// 2. Domain -> DTO Response (Para devolver un producto individual)
func ToProductResponse(p domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
	}
}

// 3. Slice Domain -> Slice DTO Response (Para el listado de productos)
func SliceProductDomainToProductResponseListDto(products []domain.Product) []dto.ProductResponse {
	if products == nil {
		return []dto.ProductResponse{}
	}

	res := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		res[i] = ToProductResponse(p)
	}
	return res
}