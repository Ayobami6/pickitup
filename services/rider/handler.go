package rider

import (
	"net/http"

	"github.com/gorilla/mux"
)

type riderHandler struct {
	repo riderRepositoryImpl
}

func NewRiderHandler(repo riderRepositoryImpl) *riderHandler {
	return &riderHandler{repo: repo}
}

func (h *riderHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register/rider", h.handleRegisterRider).Methods("POST")

}

func(h *riderHandler) handleRegisterRider(w http.ResponseWriter, r *http.Request) {
	// pend for now
}