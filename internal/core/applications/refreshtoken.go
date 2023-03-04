package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"net/http"
	"time"
)

func (m *Microservice) RefreshToken(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	body := &models.RefreshTokenRequest{}

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

	accessToken, expiresIn, err := m.UserService.RefreshToken(body.RefreshToken)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(&models.RefreshTokenResponse{AccessToken: accessToken, ExpiresIn: time.Now().Add(expiresIn)})
	if err != nil {
		return
	}
}
