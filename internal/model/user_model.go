package model

type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type LoginUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type RegisterUserRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	UserName    string `json:"user_name" validate:"required,max=50"`
	Password    string `json:"password" validate:"required,min=6,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=20"`
}
