package applications

import (
	"net/http"
)

// AuthMiddleware TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		next.ServeHTTP(w, r)
	})
}

// MiddlewareLogger TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		next.ServeHTTP(w, r)
	})
}

// MiddlewareCORS TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Stop here for a Preflighted OPTIONS request.
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareRecovery TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewareRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		next.ServeHTTP(w, r)
	})
}

// MiddlewareRateLimit TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewareRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		next.ServeHTTP(w, r)
	})
}
