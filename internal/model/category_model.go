package model

type CategoryResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Slug      string `json:"slug,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}
