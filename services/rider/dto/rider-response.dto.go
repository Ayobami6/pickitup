package dto

type RiderListResponse struct {
	ID              uint    `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	RiderID         string  `json:"rider_id"`
	BikeNumber      string  `json:"bike_number"`
	Address         string  `json:"address"`
	Rating          float64 `json:"rating"`
	SuccessfulRides int64   `json:"successful_rides"`
	Level           string  `json:"level"`
	CurrentLocation string  `json:"current_location"`
}

type RiderResponse struct {
	ID              uint    `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	RiderID         string  `json:"rider_id"`
	BikeNumber      string  `json:"bike_number"`
	Address         string  `json:"address"`
	Rating          float64 `json:"rating"`
	SuccessfulRides int64   `json:"successful_rides"`
	Level           string  `json:"level"`
	CurrentLocation string  `json:"current_location"`
	CreatedAt       string  `json:"created_at"`
	// Reviews         []ReviewResponse `json:"reviews"`
}