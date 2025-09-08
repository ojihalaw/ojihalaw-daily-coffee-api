package model

import "github.com/google/uuid"

type ProductResponse struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Slug        string    `json:"slug"`
	SKU         string    `json:"sku"`
	Variant     string    `json:"variant"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CategoryID  uuid.UUID `json:"category_id"`
	CreatedAt   string    `json:"created_at,omitempty"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
}

type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required,max=100"`
	Slug        string    `json:"slug"`
	SKU         string    `json:"sku"`
	Variant     string    `json:"variant" validate:"required,max=50"`
	Price       int       `json:"price" validate:"required"`
	Stock       int       `json:"stock" validate:"required"`
	Description string    `json:"description" validate:"required"`
	ImageURL    string    `json:"image_url"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
}

type UpdateProductRequest struct {
	Name        string    `json:"name" validate:"required,max=100"`
	Slug        string    `json:"slug"`
	SKU         string    `json:"sku"`
	Variant     string    `json:"variant" validate:"required,max=50"`
	Price       int       `json:"price" validate:"required"`
	Stock       int       `json:"stock" validate:"required"`
	Description string    `json:"description" validate:"required"`
	ImageURL    string    `json:"image_url"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
}

type UpdateSpecialProductRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	IsSpecial bool      `json:"is_special"`
}
