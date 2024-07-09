package dto

type RegisterUserDTO struct {
	UserName    string `json:"username"`
	Password    string `json:"password" validate:"required,min=6,alphanumunicode"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}