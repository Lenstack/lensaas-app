package infrastructure

import (
	"github.com/Lenstack/lensaas-app/internal/core/applications"
	"net/http"
)

type Routes struct {
	Handlers http.Handler
}

func NewRoutes(microservice applications.Microservice) *Routes {
	mux := http.NewServeMux()
	microservice.MiddlewareCORS(mux)
	microservice.MiddlewareLogger(mux)
	microservice.MiddlewareRecovery(mux)
	microservice.MiddlewareRateLimit(mux)

	mux.HandleFunc("/sign-in", microservice.SignIn)
	mux.HandleFunc("/sign-up", microservice.SignUp)
	mux.HandleFunc("/sign-out", microservice.SignOut)
	return &Routes{Handlers: mux}
}
