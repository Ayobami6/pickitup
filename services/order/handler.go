package order

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"unicode"

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
	router.HandleFunc("/orders", auth.Auth(o.handleGetOrders, o.userStore)).Methods("GET")
	router.HandleFunc("/orders/{id}/delivery", auth.UserAuth(o.handleConfirmDeliveryStatus, o.riderStore)).Methods("PATCH")
	router.HandleFunc("/orders/{id}/acknowledge", auth.RiderAuth(o.handleAcknowledge, o.riderStore)).Methods("PATCH")
	
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
	var riderId uint = uint(riderID)
    rider, err := o.riderStore.GetRiderByID(riderId)
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

func (o *orderHandler) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	// get User ID from context
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == -1 {
        auth.Forbidden(w)
        return
    }
	orders, err := o.store.GetOrders(uint(userID))
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get orders")
        return
    }
	utils.WriteJSON(w, http.StatusOK, "success", orders, "Orders retrieved successfully")
}


func (o *orderHandler) handleConfirmDeliveryStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	orderStatus := query.Get("status")
	id, err := strconv.Atoi(params["id"])
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
	orderStatus = string(unicode.ToUpper(rune(orderStatus[0]))) + orderStatus[1:]
	if orderStatus != "Delivered" {
		utils.WriteError(w, http.StatusBadRequest, "Invalid order status")
        return
    }
	var orderID uint = uint(id)
	var convertedOrderStatus models.StatusType = models.StatusType(orderStatus)

	err = o.store.UpdateDeliveryStatus(orderID, convertedOrderStatus)
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to update order status")
        return
    }
	// get orderby the id
	orderResponse, err := o.store.GetOrderByID(orderID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}
	riderID := orderResponse.RiderID
	rider, err := o.riderStore.GetRiderByID(riderID)
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get rider")
        return
    }
	riderUserID := rider.UserID
	chargeAmount := orderResponse.Charge
	riderUser, err := o.userStore.GetUserByID(int(riderUserID))
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get rider user")
        return
    }
	// TODO: add charge amount to rider wallet
	err = riderUser.Credit(o.db, chargeAmount)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to credit rider wallet")
        return
    }
    // TODO: update rider successful ride by 1
	err = rider.UpdateSuccessfulRides(o.db)
	if err!= nil {
        log.Println("Error updating Rider Successful rides")
    }

    // TODO: add email notification
	message := fmt.Sprintf("Your order delivery has been successful confirmed. ₦%.1f has been added to your wallet", chargeAmount)
	subject := "Order Delivery Notification"

	go utils.SendMail(riderUser.Email, subject, riderUser.UserName, message)

    utils.WriteJSON(w, http.StatusOK, "success", nil, "Order Delivery Successfully Confirmed")
}

func (o *orderHandler) handleAcknowledge(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Invalid ID")
		return
	}
	var ID uint = uint(orderID)
	// update the acknowledgement
	err = o.store.UpdateAcknowledgeStatus(ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
	}
	// update delivery status to indelivery
	err = o.store.UpdateDeliveryStatus(ID, models.InDelivery)
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError)
    }
	// get the updated order 
	order, Err := o.store.GetOrderByID(ID)
	if Err != nil {
		log.Println(Err)
	}
	user, userErr := o.userStore.GetUserByID(int(order.UserID))
	if userErr != nil {
		log.Println(userErr)
	}
	rider, riderErr := o.riderStore.GetRiderByID(uint(order.RiderID))
	if riderErr!= nil {
        log.Println(riderErr)
    }
	riderUser, err := o.userStore.GetUserByID(int(rider.ID))
	if err != nil {
		log.Println(err)
	}
	riderMessage := fmt.Sprintf("You have successfully acknowledged pickup order %s\n Once delivered your wallet will be automatically funded with the charge amount .\n Do make sure to ask your client to confirm your delivery before you leave", order.RefID)
	userMessage := fmt.Sprintf("Your Pickup Order %s has been acknowledged.\n Please refer to your previous email for your rider phone number so as to monitor .\n Please make sure to confirm your order delivery on your dashboard once your items has been delivered", order.RefID)
	subject := "PickItUp Order Acknowledgement"
	go utils.SendMail(user.Email, subject, user.UserName, userMessage)
    go utils.SendMail(riderUser.Email, subject, riderUser.UserName, riderMessage)

	utils.WriteJSON(w, http.StatusOK, "success", nil, "Order Successfully Acknowledged!")
}