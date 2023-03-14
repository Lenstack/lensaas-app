package models

type SignOutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SignOutResponse struct {
	Message string `json:"message"`
}
