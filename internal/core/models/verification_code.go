package models

type VerificationCodeRequest struct {
	Code string `json:"code"  validate:"required,numeric,len=7"`
}

type VerificationCodeResponse struct {
	Message string `json:"message"`
}
