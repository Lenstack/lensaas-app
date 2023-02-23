package models

type VerificationCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code"  validate:"required,numeric,len=7"`
}

type VerificationCodeResponse struct {
	Message string `json:"message"`
}
