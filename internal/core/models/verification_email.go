package models

type VerificationEmailRequest struct {
	Token string `json:"token"`
}

type VerificationEmailResponse struct {
	Message string `json:"message"`
}
