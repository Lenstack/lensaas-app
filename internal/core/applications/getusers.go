package applications

import (
	"net/http"
)

// GetUsers TODO 1. Get the tenantId from the request context, 2. Get the users from the database, 3. Return the users
func (m *Microservice) GetUsers(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	return
}
