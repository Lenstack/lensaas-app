package infrastructure

import "github.com/stripe/stripe-go/v74"

type Stripe struct {
}

func NewStripe(environment string, stripeKey string) *Stripe {
	if environment == "development" {
		stripeAppInfo := &stripe.AppInfo{
			Name:    "Lenstack",
			Version: "0.0.1",
			URL:     "https://lenstack.com",
		}
		stripe.SetAppInfo(stripeAppInfo)
	}
	return &Stripe{}
}
