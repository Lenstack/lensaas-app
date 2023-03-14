package applications

import "net/http"

func (m *Microservice) RevokeUserToken(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	return
}
