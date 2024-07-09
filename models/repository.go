package models

type RiderRepository interface {
	CreateRider(rider *Rider) error
}

type UserRepo interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}
