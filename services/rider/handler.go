package rider

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ayobami6/pickitup/models"
	"github.com/Ayobami6/pickitup/services/auth"
	"github.com/Ayobami6/pickitup/services/rider/dto"
	"github.com/Ayobami6/pickitup/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type riderHandler struct {
	repo models.RiderRepository
	userRepo models.UserRepo
}

func NewRiderHandler(repo models.RiderRepository, userRepo models.UserRepo) *riderHandler {
	return &riderHandler{repo: repo, userRepo: userRepo}
}

func (h *riderHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register/rider", h.handleRegisterRider).Methods("POST")
	router.HandleFunc("/riders", h.handleGetRiders).Methods("GET")
	router.HandleFunc("/riders/{id}", h.handleGetRider).Methods("GET")
	router.HandleFunc("/riders/charges", auth.RiderAuth(h.handleUpdateCharges, h.repo)).Methods("PATCH")

}

func(h *riderHandler) handleRegisterRider(w http.ResponseWriter, r *http.Request) {
	// pend for now
	var payload dto.RegisterRiderDTO
	if err := utils.ParseJSON(r, &payload); err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Rider Registration Details")
        return
    }

	// validate
	if vErr := utils.Validate.Struct(payload); vErr != nil {
		errs := vErr.(validator.ValidationErrors)
		if strings.Contains(errs.Error(), "Email") {
			utils.WriteError(w, http.StatusBadRequest, "Invalid Email Format")
            return
		} else if strings.Contains(errs.Error(), "Password") {
			utils.WriteError(w, http.StatusBadRequest, "Password Too Weak")
            return
		}
		log.Println(errs.Error())
		utils.WriteError(w, http.StatusBadRequest, "Bad Data!")
		return
	}
	password := payload.Password
	hashedPassword, err := auth.HashPassword(password)
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }

	user := &models.User{
		UserName: payload.UserName,
		Password: hashedPassword,
		Email: payload.Email,
		PhoneNumber: payload.PhoneNumber,
	}
	if h.userRepo == nil {
		log.Println("User Repository not provided")
		utils.WriteError(w, http.StatusInternalServerError, "User Repository not provided")
        return
    }

	newErr := h.userRepo.CreateUser(user)

	if newErr != nil {
		log.Println("Got Here")
		err := newErr.Error()
		if strings.Contains(err, "uni_users_phone_number") {
			utils.WriteError(w, http.StatusConflict, "User with this phone number already exists")
            return
		} else if strings.Contains(err, "uni_users_email") {
			utils.WriteError(w, http.StatusConflict, "User with this email already exists")
            return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }
	// create Rider
	rider := models.Rider{
        UserID: user.ID,
        FirstName: payload.FirstName,
		LastName: payload.LastName,
        Address: payload.Address,
        NextOfKinName: payload.NextOfKinName,
		NextOfKinPhone: payload.NextOfKinPhone,
        DriverLicenseNumber: payload.DriverLicenseNumber,
        NextOfKinAddress: payload.NextOfKinAddress,
		BikeNumber: payload.BikeNumber,
    }

    err = h.repo.CreateRider(&rider)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }
	// send verfication code
	num, err := utils.GenerateAndCacheVerificationCode(payload.Email)
    if err!= nil {
        log.Println("Generate Code Failed: ", err)
    } else {
        // send the email to verify
        msg := fmt.Sprintf("Your verification code is %d\n", num)
        err = utils.SendMail(payload.Email, "Email Verification", payload.UserName, msg)
        if err!= nil {
            utils.WriteError(w, http.StatusInternalServerError, "Failed to send verification email")
            return
        }
    }

    utils.WriteJSON(w, http.StatusCreated, "success", nil, "Rider Successfully Created!")
}

func (h *riderHandler) handleGetRiders(w http.ResponseWriter, r *http.Request) {
	riders, err := h.repo.GetRiders(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
	}
	utils.WriteJSON(w, http.StatusOK, "success", riders)
}

func (h *riderHandler) handleGetRider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    id := vars["id"]
	rider_id, err := strconv.Atoi(id)
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
	
    rider, err := h.repo.GetRider(rider_id, r)
	
    if err!= nil {
        utils.WriteError(w, http.StatusNotFound, "Rider not found")
    }
    utils.WriteJSON(w, http.StatusOK, "success", rider, "Rider Fetch Successfully")
}


func (h *riderHandler) handleUpdateCharges(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	log.Println(userID)
	if userID == -1 {
		auth.Forbidden(w)
		return
	}

	var payload dto.UpdateChargeDTO

	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Bad Data!")
		return 
	}

	var userId uint = uint(userID)

	err = h.repo.UpdateMinAndMaxCharge(payload.MinimumCharge, payload.MaximumCharge, userId)
	if err != nil {
		log.Println(err)
		utils.WriteError(w, http.StatusInternalServerError, "Something Went Wrong")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "status", nil, "Update Successful")

}
