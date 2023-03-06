package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"net/http"
	"time"
)

func (m *Microservice) RefreshToken(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")

	cookieValue, err := req.Cookie("refresh_token")
	if err != nil {
		wr.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: "Unauthorized", Code: http.StatusUnauthorized})
		if err != nil {
			return
		}
		return
	}

	accessToken, expiresIn, err := m.UserService.RefreshToken(cookieValue.Value)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	cookieAccessToken := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(m.TokenService.ExpirationTimeAccess),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(wr, &cookieAccessToken)

	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(&models.RefreshTokenResponse{AccessToken: accessToken, ExpiresIn: time.Now().Add(expiresIn)})
	if err != nil {
		return
	}
}
