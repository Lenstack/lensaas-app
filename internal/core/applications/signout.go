package applications

import (
	"encoding/json"
	"github.com/Lenstack/lensaas-app/internal/core/models"
	"net/http"
	"time"
)

// SignOut TODO: 1. Get user from request, 2. Validate request, 3. Call SignOut method from UserService, 4. Return success message
func (m *Microservice) SignOut(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	cookieRefresh, err := req.Cookie("refresh_token")
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	cookieRefreshToken := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	}

	cookieAccessToken := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(wr, &cookieRefreshToken)
	http.SetCookie(wr, &cookieAccessToken)

	message, err := m.UserService.RevokeToken(cookieRefresh.Value)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(wr).Encode(&models.Error{Message: err.Error(), Code: http.StatusBadRequest})
		if err != nil {
			return
		}
		return
	}

	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode(&models.SignOutResponse{Message: message})
	if err != nil {
		return
	}
}
