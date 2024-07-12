package dto

type OrderResponseDTO struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	Charge    float64 `json:"charge"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	RiderID   uint    `json:"rider_id"`
	RefID     string  `json:"ref_id"`
	Item      string  `json:"item"`
}
