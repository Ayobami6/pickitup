package dto

type LoginDTO struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}