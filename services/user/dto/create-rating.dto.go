package dto

type CreateRiderRationDTO struct {
	Rating  float64 `json:"rating" validate:"gte=0,lte=5"`
	Comment string  `json:"comment"`
}