package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"net/http"
)

// SignIn TODO: 1. Get email and password from request, 2. Validate request, 3. Call SignIn method from UserService, 4. Return token
func (m *Microservice) SignIn(wr http.ResponseWriter, req *http.Request) {
	body := &models.SignInRequest{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		_, err := wr.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	validateErrors := utils.Validate(body)
	if len(validateErrors) > 0 {
		wr.Header().Set("Content-Type", "application/json")
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(validateErrors)
		if err != nil {
			return
		}
		return
	}

	message, err := m.UserService.SignIn(body.Email, body.Password)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(err.Error())
		if err != nil {
			return
		}
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(message)
	if err != nil {
		return
	}
}
