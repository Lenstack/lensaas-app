package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"github.com/Lenstack/lensaas-app/internal/utils"
	"net/http"
	"time"
)

// SignIn TODO: 1. Get email and password from request, 2. Validate request, 3. Call SignIn method from UserService, 4. Return token
func (m *Microservice) SignIn(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	body := &models.SignInRequest{}

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

	accessToken, refreshToken, expiresIn, err := m.UserService.SignIn(body.Email, body.Password)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	// Set cookie expires in 1 day
	cookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(expiresIn),
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(wr, &cookie)

	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(&models.SignInResponse{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: time.Now().Add(expiresIn)})
	if err != nil {
		return
	}
}
