package model

type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
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
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name     string `json:"name" validate:"required,max=100"`
}
