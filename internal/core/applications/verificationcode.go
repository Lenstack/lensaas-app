package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"net/http"
)

// VerificationCode TODO 1. Get email, code from request, 2. Validate email, 3. Call EmailVerification method from UserService, 4. Return success message
func (m *Microservice) VerificationCode(wr http.ResponseWriter, req *http.Request) {
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

	body := &models.VerificationCodeRequest{}

	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	validateErrors := utils.Validate(body)
	if len(validateErrors) > 0 {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(validateErrors)
		if err != nil {
			return
		}
		return
	}

	message, err := m.UserService.VerifyCode(token, body.Code)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(&models.VerificationCodeResponse{Message: message})
	if err != nil {
		return
	}
}
