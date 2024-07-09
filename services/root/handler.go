package root

import (
	"net/http"

	"github.com/Ayobami6/pickitup/utils"
	"github.com/gorilla/mux"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return new(RootHandler)
}

func (r *RootHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", r.handleRoot).Methods("GET")
}

func (r *RootHandler) handleRoot(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"status":  "success",
		"message": "Welcome to PickItUp API V1!",
		"version": "1.0.0",
		"author": "Ayobami Alaran",
		"github":  "https://github.com/Ayobami6/pickitup",
	}
	if err := utils.WriteJSON(w, http.StatusOK, data); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}