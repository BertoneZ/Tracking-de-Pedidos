package utils

import (
	"tracking/internal/domain"
	"tracking/internal/dto"
)

func ToProductDomain(req dto.UpsertProductRequest) domain.Product {
	return domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}
}


func ToProductResponse(p domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		IsActive:    p.IsActive,
	}
}

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
