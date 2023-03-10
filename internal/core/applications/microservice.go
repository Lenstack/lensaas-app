package applications

import (
	"github.com/Lenstack/lensaas-app/internal/core/services"
)

type Microservice struct {
	EmailService  services.EmailService
	TokenService  services.TokenService
	UserService   services.UserService
	StripeService services.StripeService
}

func NewMicroservice(emailService services.EmailService, tokenService services.TokenService, userService services.UserService, stripeService services.StripeService) *Microservice {
	return &Microservice{EmailService: emailService, TokenService: tokenService, UserService: userService, StripeService: stripeService}
}
