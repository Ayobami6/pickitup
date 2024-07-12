package dto

type UpdateChargeDTO struct {
	MinimumCharge float64 `json:"min_charge"`
	MaximumCharge float64 `json:"max_charge"`
}