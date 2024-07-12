package dto

type CreateRiderRationDTO struct {
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}