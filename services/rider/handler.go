package rider

import (
	"net/http"
	"strings"

	"github.com/Ayobami6/pickitup/models"
	"github.com/Ayobami6/pickitup/services/rider/dto"
	"github.com/Ayobami6/pickitup/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type riderHandler struct {
	repo models.RiderRepository
	userRepo models.UserRepo
}

func NewRiderHandler(repo models.RiderRepository) *riderHandler {
	return &riderHandler{repo: repo}
}

func (h *riderHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register/rider", h.handleRegisterRider).Methods("POST")

}

func(h *riderHandler) handleRegisterRider(w http.ResponseWriter, r *http.Request) {
	// pend for now
	var payload dto.RegisterRiderDTO
	if err := utils.ParseJSON(r, &payload); err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Rider Registration Details")
        return
    }

	// validate
	if vErr := utils.Validate.Struct(payload); vErr!= nil {
        errs := vErr.(validator.ValidationErrors)
        if strings.Contains(errs.Error(), "Email") {
            utils.WriteError(w, http.StatusBadRequest, "Invalid Email Address")
            return
        }
    }

    // create rider
    
    // TODO: add authentication and authorization for rider registration.
    // TODO: add password hashing and salting for security.
    // TODO: add email verification.
    // TODO: add phone number verification.
    // TODO: add address validation.
}