package utils

import (
	"tracking/internal/domain"
	"tracking/internal/dto"
)

// DTO Request -> Domain 
func ToProductDomain(req dto.UpsertProductRequest) domain.Product {
	return domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}
}

// Domain -> DTO Response 
func ToProductResponse(p domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
	}
}

// Slice Domain -> Slice DTO Response 
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