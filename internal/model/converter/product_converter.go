package converter

import (
	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/model"
)

func ProductToResponse(product *entity.Product) *model.ProductResponse {
	return &model.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Slug:        product.Slug,
		SKU:         product.SKU,
		Variant:     product.Variant,
		Price:       product.Price,
		Stock:       product.Stock,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}
}
