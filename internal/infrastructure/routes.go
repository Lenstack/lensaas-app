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
	router.Use(middleware.CleanPath)
	router.Use(microservice.MiddlewareLogger)
	router.Use(microservice.MiddlewareCORS)

	// Protected Routes
	router.Group(func(router chi.Router) {
		router.Use(microservice.MiddlewareAuth)
		router.Get("/v1/users", microservice.GetUsers)           //TODO: not implemented
		router.Get("/v1/users/{id}", microservice.GetUser)       //TODO: not implemented
		router.Post("/v1/users", microservice.CreateUser)        //TODO: not implemented
		router.Put("/v1/users/{id}", microservice.UpdateUser)    //TODO: not implemented
		router.Delete("/v1/users/{id}", microservice.DeleteUser) //TODO: not implemented
	})

	router.Group(func(router chi.Router) {
		router.Post("/v1/authentication/sign_up", microservice.SignUp)             //TODO: implemented ok
		router.Post("/v1/authentication/sign_in", microservice.SignIn)             //TODO: implemented ok
		router.Post("/v1/authentication/sign_out", microservice.SignOut)           //TODO: in progress
		router.Post("/v1/authentication/refresh_token", microservice.RefreshToken) //TODO: implemented ok

		router.Post("/v1/authentication/verification_email", microservice.VerificationEmail) //TODO: implemented ok
		router.Post("/v1/authentication/verification_code", microservice.VerificationCode)   //TODO: in progress

		router.Post("/v1/authentication/password_forgot", microservice.PasswordForgot) //TODO: not implemented
		router.Post("/v1/authentication/password_reset", microservice.PasswordReset)   //TODO: not implemented
	})

	return &Routes{Handlers: router}
}
