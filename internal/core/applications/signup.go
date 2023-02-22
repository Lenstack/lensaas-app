package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/entities"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"net/http"
)

// SignUp TODO: 1. Get user from request, 2. Validate request, 3. Call SignUp method from UserService, 4. Return success message
func (m *Microservice) SignUp(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")

	body := &models.SignUpRequest{}
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
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(validateErrors)
		if err != nil {
			return
		}
		return
	}

	message, err := m.UserService.SignUp(entities.User{Name: body.Name, Email: body.Email, Password: body.Password})
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		_, err := wr.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}
	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(message)
	if err != nil {
		return
	}
}
