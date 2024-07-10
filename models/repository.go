package models

import (
	"net/http"

	"github.com/Ayobami6/pickitup/services/rider/dto"
)

type RiderRepository interface {
	CreateRider(rider *Rider) error
	GetRiders(req *http.Request) ([]dto.RiderListResponse, error)
	GetRider(id int, req *http.Request) (dto.RiderResponse, error)
	CreateRating(riderId uint) (string, error)
}

type UserRepo interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}


type OrderRepo interface {
	CreateOrder(order *Order) (error)
}