package applications

import "github.com/Lenstack/lensaas-app/internal/core/services"

type Microservice struct {
	UserService services.UserService
}

func NewMicroservice(userService services.UserService) *Microservice {
	return &Microservice{UserService: userService}
}
