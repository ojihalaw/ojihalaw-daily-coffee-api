package model

type UserLoginResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	Status      string `gorm:"size:20;default:active"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	AccessExpiresIn  int64  `json:"access_expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

type AuthResponse struct {
	UserLoginResponse UserLoginResponse `json:"user"`
	TokenResponse     TokenResponse     `json:"token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken     string `json:"access_token"`
	AccessExpiresIn int64  `json:"access_expires_in"`
}
