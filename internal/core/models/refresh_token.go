package models

import "time"

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   time.Time `json:"expires_in"`
}
