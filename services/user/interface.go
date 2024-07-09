package user

type userRepo interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}