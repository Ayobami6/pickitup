package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ayobami6/pickitup/config"
	"github.com/Ayobami6/pickitup/models"
	"github.com/Ayobami6/pickitup/services/auth"
	"github.com/Ayobami6/pickitup/services/user/dto"
	"github.com/Ayobami6/pickitup/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type userHandler struct {
	repo models.UserRepo
	riderRepo models.RiderRepository
}

func NewUserHandler(repo models.UserRepo, riderRepo models.RiderRepository) *userHandler {
	return &userHandler{repo: repo, riderRepo: riderRepo}
}

func (h *userHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/{rider_id}/ratings", auth.UserAuth(h.handleGiveRatings, h.riderRepo)).Methods("POST")
}

func (h *userHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserDTO
	err := utils.ParseJSON(r, &payload)
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Payload")
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
		utils.WriteError(w, http.StatusBadRequest, "Bad Data!")
		return
	}
	email := payload.Email

	_, errr := h.repo.GetUserByEmail(email)
	if errr == nil {
        utils.WriteError(w, http.StatusConflict, "User with this email already exists")
        return
    }
	// hash the user password before save
	password := payload.Password
	hashedPassword, err := auth.HashPassword(password)
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }
	
	// create user with the new hashed password
	newErr := h.repo.CreateUser(&models.User{
		UserName: payload.UserName,
		Password: hashedPassword,
		Email: payload.Email,
		PhoneNumber: payload.PhoneNumber,
	})
	if newErr != nil {
		err := newErr.Error()
		if strings.Contains(err, "uni_users_phone_number") {
			utils.WriteError(w, http.StatusConflict, "User with this phone number already exists")
            return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }
	// genarate and cache random number
	num, err := utils.GenerateAndCacheVerificationCode(payload.Email)
	if err!= nil {
        log.Println("Generate Code Failed: ", err)
    } else {
		// send the email to verify
		msg := fmt.Sprintf("Your verification code is %d\n", num)
		err = utils.SendMail(payload.Email, "Email Verification", payload.UserName, msg)
        if err!= nil {
            log.Printf("Email sending failed due to %v\n", err)
        }
	}

	utils.WriteJSON(w, http.StatusCreated, "success", nil, "User Successfully Created!")

}

func (h *userHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginDTO
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid Login Credentials")
		return
	}
	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		err := err.(validator.ValidationErrors)
		if strings.Contains(err.Error(), "Email") {
			utils.WriteError(w, http.StatusBadRequest, "Invalid email address")
			return
		}
	}
	u, err := h.repo.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid email address or password")
		return
	}
	// check password

	if !auth.CheckPassword(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid password")
		return
	}
	secret := []byte(config.GetEnv("JWT_SECRET", "secret"))
	token, err := auth.CreateJWT(secret, int(u.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "It from us!")
		return
	}
	utils.WriteJSON(w, http.StatusOK, "login Successful", map[string]string{"token": token})
}


func (h *userHandler) handleGiveRatings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    riderID := vars["rider_id"]

    var payload dto.CreateRiderRationDTO
    if err := utils.ParseJSON(r, &payload); err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Payload")
        return
    }
    // validate
    if vErr := utils.Validate.Struct(payload); vErr!= nil {
        err := vErr.(validator.ValidationErrors)
		if strings.Contains(err.Error(), "Rating"){
			utils.WriteError(w, http.StatusBadRequest, "Invalid Rating")
            return
		}
        utils.WriteError(w, http.StatusBadRequest, "Bad Data!")
        return
    }
	riderId, err := strconv.Atoi(riderID)
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Rider ID")
        return
    }

	var rID uint = uint(riderId)

	// create rating
	cErr := h.repo.CreateRating(&models.Review{
		RiderID: rID,
		Rating: payload.Rating,
        Comment: payload.Comment,
	})
	if cErr!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }
    // update user rating
    err = h.riderRepo.UpdateRating(rID)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }

    utils.WriteJSON(w, http.StatusOK, "success", nil, "Ratings Successfully Submitted")
}