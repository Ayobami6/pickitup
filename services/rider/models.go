package rider

import (
	"crypto/rand"
	"encoding/hex"

	"gorm.io/gorm"
)

type Rider struct {
	gorm.Model
	RiderID    string `json:"rider_id" gorm:"uniqueIndex;size:8"`
	BikeNumber string `json:"bike_number"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserID uint 
	NextOfKinName string `json:"next_of_kin_name"`
	NextOfKinPhone string `json:"next_of_kin_phone"`
	DriverLicenseNumber string `json:"driver_license_number"`
	NextOfKinAddress string `json:"next_of_kin_address"`
	Address string `json:"address"`
	SuccessfulRides int64 `json:"successful_rides"`
	Rating float64  `json:"rating"`
	Level string  `json:"level"`
	CurrentLocation string `json:"current_location"`
	Review []Review `json:"reviews" gorm:"references:reviews"`
}



type Review struct {
	gorm.Model
	RiderID uint `json:"rider_id"`
	Rating float64 `json:"rating"`
	Comment string `json:"comment"`
}

func generateRandomID(n int) (string, error) {
    bytes := make([]byte, n)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes)[:n], nil
}

func (u *Rider) BeforeCreate(tx *gorm.DB) (err error) {
    u.RiderID, err = generateRandomID(8)
    if err != nil {
        return err
    }
    return nil
}