package models

import (
	"net/http"

	orderDto "github.com/Ayobami6/pickitup/services/order/dto"
	"github.com/Ayobami6/pickitup/services/rider/dto"
)

type RiderRepository interface {
	CreateRider(rider *Rider) error
	GetRiders(req *http.Request) ([]dto.RiderListResponse, error)
	GetRider(id int, req *http.Request) (dto.RiderResponse, error)
	GetRiderByUserID(userID uint) (*Rider, error)
	GetRiderByID(id uint)(*Rider, error)
	UpdateRating(riderId uint) (error)
	UpdateMinAndMaxCharge(minCharge float64, maxCharge float64, userID uint) (error)
	GetRiderReviews(rideId uint)([]Review, error)
}

type UserRepo interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateRating(rating *Review) error
}


type OrderRepo interface {
	CreateOrder(order *Order) (error)
	GetOrders(userID uint)([]orderDto.OrderResponseDTO, error)
	UpdateDeliveryStatus(orderId uint, status StatusType) error
	UpdateAcknowledgeStatus(orderId uint) error
	GetOrderByID(orderId uint) (orderDto.OrderResponseDTO, error)
}