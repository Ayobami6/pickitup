package dto

type CreateRiderRationDTO struct {
	RiderID int     `json:"rider_id"`
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}