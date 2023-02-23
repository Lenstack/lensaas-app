package infrastructure

import (
	"github.com/Lenstack/lensaas-app/internal/core/applications"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Routes struct {
	Handlers http.Handler
}

func NewRoutes(microservice applications.Microservice) *Routes {
	router := chi.NewRouter()
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(microservice.MiddlewareLogger)
	router.Use(microservice.MiddlewareCORS)

	// Protected Routes
	router.Group(func(router chi.Router) {
		router.Use(microservice.MiddlewareAuth)
		router.Get("/v1/users", microservice.GetUsers)
		router.Get("/v1/users/{id}", microservice.GetUser)
		router.Post("/v1/users", microservice.CreateUser)
		router.Put("/v1/users/{id}", microservice.UpdateUser)
		router.Delete("/v1/users/{id}", microservice.DeleteUser)
	})

	router.Group(func(router chi.Router) {
		router.Post("/v1/authentication/sign_up", microservice.SignUp)   //TODO: implemented ok
		router.Post("/v1/authentication/sign_in", microservice.SignIn)   //TODO: implemented ok
		router.Post("/v1/authentication/sign_out", microservice.SignOut) //TODO: in progress
	})

	return &Routes{Handlers: router}
}
