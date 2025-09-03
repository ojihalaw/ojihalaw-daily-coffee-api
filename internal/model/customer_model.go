package model

type CustomerResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type LoginCustomerRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutCustomerRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetCustomerRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type RegisterCustomerRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	UserName    string `json:"user_name" validate:"required,max=50"`
	Email       string `json:"email" validate:"required,max=50"`
	Password    string `json:"password" validate:"required,min=6,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=20"`
}
