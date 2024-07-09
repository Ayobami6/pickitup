package dto

type RegisterRiderDTO struct {
	Email               string `json:"email" validate:"required,email"`
	FirstName           string `json:"first_name" validate:"required"`
	LastName            string `json:"last_name" validate:"required"`
	NextOfKinAddress    string `json:"next_of_kin_address" validate:"required"`
	NextOfKinName       string `json:"next_of_kin_name" validate:"required"`
	NextOfKinPhone      string `json:"next_of_kin_phone" validate:"required"`
	DriverLicenseNumber string `json:"driver_license_number" validate:"required"`
	BikeNumber          string `json:"bike_number" validate:"required"`
	Address             string `json:"address" validate:"required"`
	UserName            string `json:"username" validate:"required"`
	Password            string `json:"password" validate:"required,min=6,alphanumunicode"`
	PhoneNumber         string `json:"phone_number" validate:"required"`
}