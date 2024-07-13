package dto

type RiderListResponse struct {
	ID                 uint    `json:"id"`
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	RiderID            string  `json:"rider_id"`
	BikeNumber         string  `json:"bike_number"`
	Address            string  `json:"address"`
	Rating             float64 `json:"rating"`
	SuccessfulRides    int64   `json:"successful_rides"`
	Level              string  `json:"level"`
	CurrentLocation    string  `json:"current_location"`
	SelfUrl            string  `json:"self_url"`
	MaximumCharge      float64 `json:"maximum_charge"`
	MinimumCharge      float64 `json:"minimum_charge"`
	AvailabilityStatus string  `json:"availability_status"`
}

type RiderResponse struct {
	ID                 uint             `json:"id"`
	FirstName          string           `json:"first_name"`
	LastName           string           `json:"last_name"`
	RiderID            string           `json:"rider_id"`
	BikeNumber         string           `json:"bike_number"`
	Address            string           `json:"address"`
	Rating             float64          `json:"rating"`
	SuccessfulRides    int64            `json:"successful_rides"`
	Level              string           `json:"level"`
	CurrentLocation    string           `json:"current_location"`
	CreatedAt          string           `json:"created_at"`
	SelfUrl            string           `json:"self_url"`
	MaximumCharge      float64          `json:"maximum_charge"`
	MinimumCharge      float64          `json:"minimum_charge"`
	Reviews            []ReviewResponse `json:"reviews"`
	AvailabilityStatus string           `json:"availability_status"`
}

type ReviewResponse struct {
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}