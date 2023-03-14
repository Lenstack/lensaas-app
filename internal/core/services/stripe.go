package services

type IStripeService interface {
	CreateCustomer()
	CreateSubscription()
}

type StripeService struct {
}

func NewStripeService() *StripeService {
	return &StripeService{}
}
