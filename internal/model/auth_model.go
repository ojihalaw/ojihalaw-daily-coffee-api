package model

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

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
