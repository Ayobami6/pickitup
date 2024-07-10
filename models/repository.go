package models

import "github.com/Ayobami6/pickitup/services/rider/dto"

type RiderRepository interface {
	CreateRider(rider *Rider) error
	GetRiders() ([]dto.RiderListResponse, error)
	GetRider(id int) (Rider, error)
	CreateRating(riderId uint) (string, error)
}

type UserRepo interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}
