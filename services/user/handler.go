package user

import (
	"net/http"
	"strings"

	"github.com/Ayobami6/pickitup/services/auth"
	"github.com/Ayobami6/pickitup/services/user/dto"
	"github.com/Ayobami6/pickitup/utils"
	"github.com/gorilla/mux"
)

type userHandler struct {
	repo userRepoImpl
}

func NewUserHandler(repo userRepoImpl) *userHandler {
	return &userHandler{repo: repo}
}

func (h *userHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *userHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterUserDTO
	err := utils.ParseJSON(r, &payload)
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Payload")
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
	newErr := h.repo.CreateUser(&User{
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

	utils.WriteJSON(w, http.StatusCreated, "success", nil, "User Successfully Created!")



}