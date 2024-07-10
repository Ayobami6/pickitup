package order

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Ayobami6/pickitup/models"
	"github.com/Ayobami6/pickitup/services/auth"
	"github.com/Ayobami6/pickitup/services/order/dto"
	"github.com/Ayobami6/pickitup/utils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type orderHandler struct {
	store models.OrderRepo
	userStore models.UserRepo
	riderStore models.RiderRepository
	db *gorm.DB
}

func NewOrderHandler(store models.OrderRepo, us models.UserRepo, rs models.RiderRepository, db *gorm.DB) *orderHandler {
    return &orderHandler{store: store, userStore:us, riderStore: rs, db: db}
}


func (o *orderHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders/{rider_id}", auth.Auth(o.handleCreateOrder, o.userStore)).Methods("POST")
	
}

func (o *orderHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// deserialize order 
	var payload dto.CreateOrderDTO 
	if err := utils.ParseJSON(r, &payload); err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Order Details")
        return
    }
	// get rider id from param
	params := mux.Vars(r)
    riderID, err := strconv.Atoi(params["rider_id"])
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid rider id")
        return
    }
    // get the rider
    rider, err := o.riderStore.GetRiderByID(riderID)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get rider")
        return
    }
    // check if the rider is available
	riderAvailableStatus := rider.AvailabilityStatus
	switch {
		case riderAvailableStatus == "Unavailable":
        case riderAvailableStatus == "OnBreak":
		case riderAvailableStatus == "Busy":
            utils.WriteError(w, http.StatusNotFound, "Rider is currently unavailable")
            return
        default:
			break  
	}
	minCharge := rider.MinimumCharge
	maxCharge := rider.MaximumCharge
	charge := minCharge + ((maxCharge - minCharge)/2)

	// get the request context
	ctx := r.Context()
	// get user Id from context
	userID := auth.GetUserIDFromContext(ctx)
	if userID == -1 {
        auth.Forbidden(w)
        return
    }
	// get the user
	user, err := o.userStore.GetUserByID(userID)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get user")
        return
    }
	// check if the user has enough balance
	if !user.Verified {
		utils.WriteError(w, http.StatusNotFound, "User is not verified")
        return
	}

	if bal := user.WalletBalance; bal <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Insufficient balance")
        return
	}
	// charge the user
	err = user.Debit(o.db, charge)
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to deduct balance")
        return
    }
	// lets create the order
	order := &models.Order{
        RiderID: rider.ID,
        UserID: user.ID,
        Charge: charge,
        Item: payload.Item,
		Quantity: payload.Quantity,
		PickUpAddress: payload.PickUpAddress,
        DropOffAddress: payload.DropOffAddress,
    }
    err = o.store.CreateOrder(order)
	log.Println(err)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to create order")
        return
    }
	data := map[string]string{
		"ref_id": order.RefID,
	}
	riderUser, err := o.userStore.GetUserByID(int(rider.UserID))
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get rider user")
        return
    }
	riderMessage := fmt.Sprintf("You have New Pick Up Order with ID %s\n Containing item %s which is to be picked up at %s \n and delivered at %s Please go to your dashboard to accept the order and transit immediately or reject \n", order.RefID, order.Item, order.PickUpAddress, order.DropOffAddress)
	userMessage := fmt.Sprintf("Your Order %s has been placed successfully \n Here is your rider phone number %s\n\n", order.RefID, riderUser.PhoneNumber)
	subject := "PickItUp Order Notification"
	// send mail to user and rider user
	go utils.SendMail(user.Email, subject, user.UserName, userMessage)
    go utils.SendMail(riderUser.Email, subject, riderUser.UserName, riderMessage)

    // write the response
    utils.WriteJSON(w, http.StatusCreated, "success", data, "Order created successfully")


}