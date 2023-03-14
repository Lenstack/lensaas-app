package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"net/http"
)

// VerificationEmail TODO 1. Get token from request, 2. Validate token, 3. Call EmailVerification method from UserService, 4. Return success message
func (m *Microservice) VerificationEmail(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")

	token := req.URL.Query().Get("token")
	if token == "" {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: "token is required", Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	message, err := m.UserService.VerifyEmail(token)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(&models.VerificationEmailResponse{Message: message})
	if err != nil {
		return
	}
}
