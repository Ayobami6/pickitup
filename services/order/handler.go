package order

import (
	"net/http"

	"github.com/Ayobami6/pickitup/models"
	"github.com/gorilla/mux"
)

type orderHandler struct {
	store models.OrderRepo
	userStore models.UserRepo
	riderStore models.RiderRepository
}

func NewOrderHandler(store models.OrderRepo, us models.UserRepo, rs models.RiderRepository) *orderHandler {
    return &orderHandler{store: store, userStore:us, riderStore: rs}
}


func (o *orderHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders", o.handleCreateOrder).Methods("POST")
	
}

func (o *orderHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// 
}