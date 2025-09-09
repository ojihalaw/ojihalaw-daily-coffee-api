package model

type CustomerResponse struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
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

type UpdateCustomerRequest struct {
	Name        string `json:"name"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role" `
}
