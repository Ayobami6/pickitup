package dto

type CreateOrderDTO struct {
	UserID         string `json:"user_id"`
	RiderID        string `json:"rider_id"`
	Item           string `json:"item"`
	Quantity       int    `json:"quantity"`
	PickUpAddress  string `json:"pickup_address"`
	DropOffAddress string `json:"dropoff_address"`
}