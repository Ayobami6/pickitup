package dto

type CreateOrderDTO struct {
	Item           string `json:"item"`
	Quantity       int    `json:"quantity"`
	PickUpAddress  string `json:"pickup_address"`
	DropOffAddress string `json:"dropoff_address"`
}