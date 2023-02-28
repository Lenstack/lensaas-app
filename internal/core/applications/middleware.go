package applications

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"net/http"
	"strings"
)

// MiddlewareAuth TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(&models.Error{Message: "Unauthorized", Code: http.StatusUnauthorized})
			if err != nil {
				return
			}
			return
		}

		parts := strings.Fields(token)
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(&models.Error{Message: "Unauthorized", Code: http.StatusUnauthorized})
			if err != nil {
				return
			}
			return
		}

		clearedToken := parts[1]

		userId, err := m.TokenService.ValidateToken(clearedToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(&models.Error{Message: "Unauthorized", Code: http.StatusUnauthorized})
			if err != nil {
				return
			}
			return
		}

		// Set the userId in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// MiddlewarePermission TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewarePermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		// If the user is not authorized, return a 401
		next.ServeHTTP(w, r)
	})
}

// MiddlewareLogger TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		fmt.Println("Request: ", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// MiddlewareCORS TODO 1. Add a middleware to the microservice, 2. Add a middleware to the routes
func (m *Microservice) MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
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
